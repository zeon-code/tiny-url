package test

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/http/handler"
	"github.com/zeon-code/tiny-url/internal/repository"
	"github.com/zeon-code/tiny-url/internal/service"
)

type FakeDependencies struct {
	// Database
	DBMock    sqlmock.Sqlmock
	DBBackend *sqlx.DB
	DBMetric  *FakeMetric

	// Cache
	CacheBackend *FakeRedis
	CacheMetric  *FakeMetric

	// Memory
	MemoryMetric *FakeMetric

	// HTTP
	HTTPMetric *FakeMetric
}

func NewFakeDependencies() FakeDependencies {
	var sqldb *sql.DB

	fake := FakeDependencies{
		DBMetric:     NewFakeMetric(),
		CacheMetric:  NewFakeMetric(),
		MemoryMetric: NewFakeMetric(),
		HTTPMetric:   NewFakeMetric(),
		CacheBackend: NewFakeRedisBackend(),
	}

	sqldb, fake.DBMock, _ = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	fake.DBBackend = sqlx.NewDb(sqldb, "postgres")

	return fake
}

func (d FakeDependencies) DB() *db.PostgresClient {
	return db.NewPostgresClient(d.DBBackend, d.DBMetric, d.Logger())
}

func (d FakeDependencies) Cache() *db.RedisClient {
	return db.NewRedisClient(d.CacheBackend, d.CacheMetric, d.Logger())
}

func (d FakeDependencies) Memory() db.SQLReader {
	return db.NewMemoryDatabase(d.DB(), d.Cache(), d.MemoryMetric, d.Logger())
}

func (d FakeDependencies) Logger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func (d FakeDependencies) Repositories() repository.Repositories {
	return repository.NewRepositories(d.DB(), d.Memory(), d.Logger())
}

func (d FakeDependencies) Services() service.Services {
	return service.NewServices(d.Repositories(), d.Logger())
}

func (d FakeDependencies) Router() http.Handler {
	return handler.NewRouter(d.Services(), d.HTTPMetric, d.Logger())
}
