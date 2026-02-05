package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrDBInvalidBackend = errors.New("error db invalid backend instance")
	ErrDBNotFound       = errors.New("error db resource not found")
)

func mapDBError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrDBNotFound
	case errors.Is(err, context.Canceled):
		return err
	case errors.Is(err, context.DeadlineExceeded):
		return err
	}

	return err
}

var (
	ErrCacheNotFound    = errors.New("error cache not found")
	ErrCacheUnavailable = errors.New("error cache unavailable")
)

func mapCacheError(err error) error {
	switch {
	case errors.Is(err, redis.Nil):
		return ErrCacheNotFound
	case errors.Is(err, context.Canceled):
		return err
	case errors.Is(err, context.DeadlineExceeded):
		return err
	}

	return ErrCacheUnavailable
}
