package repository

import (
	"context"
	"errors"

	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
)

type HealthRepository interface {
	Ping(context.Context) (string, error)
}

type HealthStore struct {
	primary db.SQLClient
	memory  db.SQLReader
	logger  observability.Logger
}

func NewHealthRepository(primary db.SQLClient, memory db.SQLReader, observer observability.Observer) HealthRepository {
	return HealthStore{
		primary: primary,
		memory:  memory,
		logger:  observer.Logger().With("repository", "health"),
	}
}

func (r HealthStore) Ping(ctx context.Context) (string, error) {
	if err := r.primary.Ping(ctx); err != nil {
		return "db_primary_unavailable", errors.New("error dependency not ready")
	}

	if err := r.memory.Ping(ctx); err != nil {
		return "memory_unavailable", errors.New("error dependency not ready")
	}

	return "", nil
}
