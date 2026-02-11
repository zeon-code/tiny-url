package db

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
)

type dbFetch func(context.Context, any, string, ...any) error

// MemoryDatabaseClient decorates a SQLReader with a transparent,
// best-effort cache layer.
//
// It attempts to serve read operations from cache when caching is enabled
// via context. Cache failures or corrupt entries never block database access;
// in those cases the request transparently falls back to the underlying
// SQLReader.
//
// Values passed to this client must be pointers to JSON-marshalable types.
type MemoryDatabaseClient struct {
	db     SQLClient
	cache  CacheClient
	metric observability.MetricClient
	logger observability.Logger
}

func NewMemoryDatabase(db SQLClient, cache CacheClient, observer observability.Observer) (SQLReader, error) {
	metrics, err := observer.Metric()

	if err != nil {
		return nil, err
	}

	return MemoryDatabaseClient{
		db:     db,
		cache:  cache,
		metric: metrics,
		logger: observer.Logger().With("client", "memory"),
	}, nil
}

// Select executes a database select operation with optional caching.
//
// When caching is enabled in the context, Select attempts to retrieve the
// result from cache before querying the database. On a cache hit, the cached
// JSON payload is decoded into value and the database is not queried.
//
// If the cache entry is missing, invalid, or cannot be decoded, the cache
// entry is discarded and the query is executed against the database. On a
// successful database response, the result is stored in cache according to
// the cache policy defined in the context.
//
// Cache errors do not affect the database execution path.
func (c MemoryDatabaseClient) Select(ctx context.Context, value any, query string, args ...any) error {
	return c.load(ctx, c.db.Select, value, query, args...)
}

// Get executes a database get operation with optional caching.
//
// Behavior mirrors Select, but is intended for single-row or primary-key
// lookups. When caching is enabled, Get attempts to resolve the value from
// cache before querying the database, falling back to the database on cache
// misses or decode failures.
//
// Cache errors do not affect the database execution path.
func (c MemoryDatabaseClient) Get(ctx context.Context, value any, query string, args ...any) error {
	return c.load(ctx, c.db.Get, value, query, args...)
}

// Get executes a database get operation with optional caching.
//
// Behavior mirrors Select, but is intended for single-row or primary-key
// lookups. When caching is enabled, Get attempts to resolve the value from
// cache before querying the database, falling back to the database on cache
// misses or decode failures.
//
// Cache errors do not affect the database execution path.
func (c MemoryDatabaseClient) Ping(ctx context.Context) error {
	if err := c.cache.Ping(ctx); err != nil {
		c.logger.Warn(ctx, "error cache is not available", slog.Any("error", err))
		observability.TraceError(ctx, "cache is not available", err)
	}

	return c.db.Ping(ctx)
}

// Close closes the underlying database connection and cache client,
// releasing all associated resources.
//
// Both close operations are attempted, and any resulting errors are
// combined using errors.Join and returned to the caller.
func (c MemoryDatabaseClient) Close() error {
	return errors.Join(c.cache.Close(), c.db.Close())
}

func (c MemoryDatabaseClient) load(ctx context.Context, fetch dbFetch, value any, query string, args ...any) error {
	startAt := time.Now()
	memory := cache.CacheFromContext(ctx)

	if !memory.IsEnabled {
		c.metric.MemoryBypassed(ctx)
		return fetch(ctx, value, query, args...)
	}

	if data, err := c.cache.Get(ctx, memory.Policy.Key); err == nil {
		if err := json.Unmarshal(data, value); err == nil {
			c.metric.MemoryHit(ctx, memory.Policy.Key, time.Since(startAt))
			return nil
		}

		c.metric.MemoryInvalid(ctx, memory.Policy.Key)
		c.cache.Del(ctx, memory.Policy.Key)
	}

	if err := fetch(ctx, value, query, args...); err != nil {
		return err
	}

	if data, err := json.Marshal(value); err == nil {
		c.cache.Set(ctx, data, memory.Policy.Key, memory.Policy.TTL)
	}

	c.metric.MemoryMiss(ctx, memory.Policy.Key, time.Since(startAt))
	return nil
}
