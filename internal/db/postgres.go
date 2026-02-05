package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"github.com/zeon-code/tiny-url/internal/pkg/metric"

	"github.com/jmoiron/sqlx"
)

type PostgresBackend interface {
	SelectContext(ctx context.Context, value any, query string, args ...any) error
	GetContext(ctx context.Context, value any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type PostgresTxBackend interface {
	PostgresBackend

	Commit() error
	Rollback() error
}

type PostgresClientBackend interface {
	PostgresBackend

	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// PostgresProxy provides a thin abstraction over sqlx.DB,
// centralizing database access and normalizing error handling.
// It delegates query execution to the underlying sqlx backend
// while mapping low-level database errors to domain-level errors.
type PostgresProxy struct {
	backend PostgresBackend
	metric  metric.MetricClient
	logger  *slog.Logger
}

// Select executes a query against the database and populates
// the provided value with the result. `value` must be a pointer
// to the destination type, `query` is the SQL query string, and
// `args` are any query parameters.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresProxy) Select(ctx context.Context, value any, query string, args ...any) error {
	startAt := time.Now()
	err := p.backend.SelectContext(ctx, value, query, args...)

	p.track(query, startAt, err)
	return mapDBError(err)
}

// Get executes a query against the database and populates the
// provided value with the first row returned. `value` must be a
// pointer to the destination type. `query` is the SQL query string,
// and `args` are any query parameters.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresProxy) Get(ctx context.Context, value any, query string, args ...any) error {
	startAt := time.Now()
	err := p.backend.GetContext(ctx, value, query, args...)

	p.track(query, startAt, err)
	return mapDBError(err)
}

// Exec executes a query against the database that does not return rows,
// such as INSERT, UPDATE, or DELETE. `query` is the SQL query string,
// and `args` are any query parameters.
//
// Returns a mapped error using mapDBError for consistent error handling.
func (p PostgresProxy) Exec(ctx context.Context, query string, args ...any) error {
	startAt := time.Now()
	_, err := p.backend.ExecContext(ctx, query, args...)

	p.track(query, startAt, err)
	return mapDBError(err)
}

func (p PostgresProxy) track(query string, startAt time.Time, err error) {
	p.metric.DBQuery(query, time.Since(startAt))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		p.metric.DBError(query, err.Error())
	}
}

type PostgresClient struct {
	PostgresProxy
}

func newPostgresClient(c config.DatabaseConfiguration, m metric.MetricClient, l *slog.Logger) (SQLClient, error) {
	backend, err := sqlx.Connect(c.GetDriver(), c.GetDNS())

	if err != nil {
		return nil, mapDBError(err)
	}

	return NewPostgresClient(backend, m, l), nil
}

func NewPostgresClient(b PostgresClientBackend, m metric.MetricClient, l *slog.Logger) *PostgresClient {
	return &PostgresClient{
		PostgresProxy: PostgresProxy{
			backend: b,
			metric:  m,
			logger:  l,
		},
	}
}

// BeginTx begins a new SQL transaction with the given context and options.
// It returns a transactional SQLTX that guarantees atomic execution.
// The caller is responsible for committing or rolling back the transaction.
func (p PostgresClient) BeginTx(ctx context.Context, opt *sql.TxOptions) (SQLTX, error) {
	startAt := time.Now()
	backend, ok := p.backend.(*sqlx.DB)

	if !ok {
		p.track("START TRANSACTION;", startAt, ErrDBInvalidBackend)
		return nil, ErrDBInvalidBackend
	}

	tx, err := backend.BeginTxx(ctx, opt)
	p.track("START TRANSACTION;", startAt, err)

	if err != nil {
		return nil, mapDBError(err)
	}

	return newPostgresTx(tx, p.metric), nil
}

type PostgresTX struct {
	PostgresProxy
}

func newPostgresTx(tx PostgresTxBackend, m metric.MetricClient) SQLTX {
	return &PostgresTX{
		PostgresProxy: PostgresProxy{
			backend: tx,
			metric:  m,
		},
	}
}

// Commit commits the current transaction and releases all associated resources.
// Once committed, the transaction is closed and further operations will fail.
func (p PostgresTX) Commit() error {
	startAt := time.Now()
	backend, ok := p.backend.(PostgresTxBackend)

	if !ok {
		p.track("COMMIT TRANSACTION;", startAt, ErrDBInvalidBackend)
		return ErrDBInvalidBackend
	}

	err := backend.Commit()
	p.track("COMMIT TRANSACTION;", startAt, err)
	return mapDBError(err)
}

// Rollback roll back the transaction and releases all associated resources.
// Calling Rollback after Commit has no effect.
func (p PostgresTX) Rollback() error {
	startAt := time.Now()
	backend, ok := p.backend.(PostgresTxBackend)

	if !ok {
		p.track("ROLLBACK;", startAt, ErrDBInvalidBackend)
		return ErrDBInvalidBackend
	}

	err := backend.Rollback()
	p.track("ROLLBACK;", startAt, err)
	return mapDBError(err)
}
