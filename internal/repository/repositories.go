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

func NewRepositoriesFromConfig(c config.Configuration, logger *slog.Logger) Repositories {
	metric := metric.NewMetricClient(c, logger.With("client", "metric"))

	cache, err := db.NewCacheClient(c.Cache(), metric, logger.With("client", "cache"))

	if err != nil {
		panic("error building cache client: " + err.Error())
	}

	database, err := db.NewDBClient(c.PrimaryDatabase(), metric, logger.With("client", "primary-db"))

	if err != nil {
		panic("error building primary database client: " + err.Error())
	}

	replica, err := db.NewDBClient(c.ReplicaDatabase(), metric, logger.With("client", "replica-db"))

	if err != nil {
		replica = database
	}

	memory := db.NewMemoryDatabase(replica, cache, metric, logger.With("client", "memory"))
	return NewRepositories(database, memory, logger)
}

func NewRepositories(database db.SQLClient, memory db.SQLReader, logger *slog.Logger) Repositories {
	return Repositories{
		Url: NewURLRepository(database, memory, logger.With("repository", "url")),
	}
}
