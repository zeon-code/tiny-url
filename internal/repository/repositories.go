package repository

import (
	"log/slog"

	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"github.com/zeon-code/tiny-url/internal/pkg/metric"
)

type Repositories struct {
	Url URLRepository
}

// NewRepositoriesFromConfig builds and wires all repository dependencies using the
// provided application configuration and logger.
//
// It initializes metric, cache, primary database, and replica database clients.
// If the replica database configuration is missing or fails to initialize, the
// primary database client is used as a fallback to ensure read availability.
//
// The function panics if critical dependencies (cache or primary database)
// cannot be created, as the application cannot operate without them.
//
// Returns a fully initialized Repositories instance
func NewRepositoriesFromConfig(conf config.Configuration, metrics metric.MetricClient, logger *slog.Logger) Repositories {
	cache, err := db.NewCacheClient(conf.Cache(), metrics, logger.With("client", "cache"))

	if err != nil {
		panic("error building cache client: " + err.Error())
	}

	database, err := db.NewDBClient(conf.PrimaryDatabase(), metrics, logger.With("client", "primary-db"))

	if err != nil {
		panic("error building primary database client: " + err.Error())
	}

	replica, err := db.NewDBClient(conf.ReplicaDatabase(), metrics, logger.With("client", "replica-db"))

	if err != nil {
		replica = database
	}

	memory := db.NewMemoryDatabase(replica, cache, metrics, logger.With("client", "memory"))
	return NewRepositories(database, memory, logger)
}

// NewRepositories constructs a Repositories container using the provided
// database and memory-backed readers.
//
// The database client is used for write operations, while the memory client
// (typically backed by cache and/or replicas) is used for read operations.
func NewRepositories(database db.SQLClient, memory db.SQLReader, logger *slog.Logger) Repositories {
	return Repositories{
		Url: NewURLRepository(database, memory, logger.With("repository", "url")),
	}
}
