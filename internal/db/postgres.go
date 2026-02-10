package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"

	"github.com/jmoiron/sqlx"
)

type PostgresBackend interface {
	SelectContext(context.Context, any, string, ...any) error
	GetContext(context.Context, any, string, ...any) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type PostgresTxBackend interface {
	PostgresBackend

	Commit() error
	Rollback() error
}

type PostgresClientBackend interface {
	PostgresBackend

	Close() error
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

// PostgresProxy provides a thin abstraction over sqlx.DB,
// centralizing database access and normalizing error handling.
// It delegates query execution to the underlying sqlx backend
// while mapping low-level database errors to domain-level errors.
type PostgresProxy struct {
	backend PostgresBackend
	logger  observability.Logger
}

// Select executes a query against the database and populates
// the provided value with the result. `value` must be a pointer
// to the destination type, `query` is the SQL query string, and
// `args` are any query parameters.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresProxy) Select(ctx context.Context, value any, query string, args ...any) error {
	err := p.backend.SelectContext(ctx, value, query, args...)
	return mapDBError(err)
}

// Get executes a query against the database and populates the
// provided value with the first row returned. `value` must be a
// pointer to the destination type. `query` is the SQL query string,
// and `args` are any query parameters.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresProxy) Get(ctx context.Context, value any, query string, args ...any) error {
	err := p.backend.GetContext(ctx, value, query, args...)
	return mapDBError(err)
}

// Exec executes a query against the database that does not return rows,
// such as INSERT, UPDATE, or DELETE. `query` is the SQL query string,
// and `args` are any query parameters.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresProxy) Exec(ctx context.Context, query string, args ...any) error {
	_, err := p.backend.ExecContext(ctx, query, args...)
	return mapDBError(err)
}

type PostgresClient struct {
	PostgresProxy

	isConnectionClosed bool
}

func NewPostgresClientFromConfig(conf config.DatabaseConfiguration, observer observability.Observer) (SQLClient, error) {
	dns, err := conf.DSN()

	if err != nil {
		return nil, err
	}

	db, err := observability.NewInstrumentedDB(observer, conf.Driver(), dns)

	if err != nil {
		return nil, mapDBError(err)
	}

	return NewPostgresClient(sqlx.NewDb(db, "postgres"), observer), nil
}

func NewPostgresClient(backend PostgresClientBackend, observer observability.Observer) *PostgresClient {
	return &PostgresClient{
		PostgresProxy: PostgresProxy{
			backend: backend,
			logger:  observer.Logger().WithGroup("postgres-client"),
		},
	}
}

// BeginTx begins a new SQL transaction with the given context and options.
// It returns a transactional SQLTX that guarantees atomic execution.
// The caller is responsible for committing or rolling back the transaction.
func (p PostgresClient) BeginTx(ctx context.Context, opt *sql.TxOptions) (SQLTX, error) {
	backend, ok := p.backend.(*sqlx.DB)

	if !ok {
		return nil, ErrDBInvalidBackend
	}

	tx, err := backend.BeginTxx(ctx, opt)

	if err != nil {
		return nil, mapDBError(err)
	}

	return newPostgresTx(tx, p.logger), nil
}

// Close closes the underlying PostgreSQL backend connection.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresClient) Close() error {
	if !p.isConnectionClosed {
		p.isConnectionClosed = true
		return mapDBError(p.backend.(PostgresClientBackend).Close())
	}

	return nil
}

type PostgresTX struct {
	PostgresProxy
}

func newPostgresTx(tx PostgresTxBackend, logger observability.Logger) SQLTX {
	return &PostgresTX{
		PostgresProxy: PostgresProxy{
			backend: tx,
			logger:  logger.WithGroup("postgres-tx-client"),
		},
	}
}

// Commit commits the current transaction and releases all associated resources.
// Once committed, the transaction is closed and further operations will fail.
func (p PostgresTX) Commit() error {
	backend, ok := p.backend.(PostgresTxBackend)

	if !ok {
		return ErrDBInvalidBackend
	}

	err := backend.Commit()
	return mapDBError(err)
}

// Rollback roll back the transaction and releases all associated resources.
// Calling Rollback after Commit has no effect.
func (p PostgresTX) Rollback() error {
	backend, ok := p.backend.(PostgresTxBackend)

	if !ok {
		return ErrDBInvalidBackend
	}

	err := backend.Rollback()
	return mapDBError(err)
}
