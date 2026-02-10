package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
)

type SQLReader interface {
	Close() error

	Select(context.Context, any, string, ...any) error
	Get(context.Context, any, string, ...any) error
}

type SQLTX interface {
	Commit() error
	Rollback() error

	Select(context.Context, any, string, ...any) error
	Get(context.Context, any, string, ...any) error
	Exec(context.Context, string, ...any) error
}

type SQLClient interface {
	SQLReader

	Exec(context.Context, string, ...any) error
	BeginTx(context.Context, *sql.TxOptions) (SQLTX, error)
}

func NewDBClient(conf config.DatabaseConfiguration, observer observability.Observer) (SQLClient, error) {
	return NewPostgresClientFromConfig(conf, observer)
}

type CacheClient interface {
	Del(context.Context, string) error
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, any, string, time.Duration) error
	Incr(context.Context, string) (int64, error)
	Close() error
}

func NewCacheClient(conf config.DatabaseConfiguration, observer observability.Observer) (CacheClient, error) {
	return NewRedisClientFromConfig(conf, observer)
}
