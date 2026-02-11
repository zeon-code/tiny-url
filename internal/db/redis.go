package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeon-code/tiny-url/internal/pkg/config"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
)

type RedisBackend interface {
	Ping(context.Context) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Del(context.Context, ...string) *redis.IntCmd
	Incr(context.Context, string) *redis.IntCmd
	SetNX(context.Context, string, interface{}, time.Duration) *redis.BoolCmd
	Close() error
}

// RedisClient provides a thin abstraction over redis.Client,
// centralizing cache operations and normalizing cache-related
// error handling. It delegates commands to the underlying Redis
// backend while mapping low-level errors to domain-level errors.
type RedisClient struct {
	backend RedisBackend
	metric  observability.MetricClient
	logger  observability.Logger
}

func NewRedisClientFromConfig(conf config.DatabaseConfiguration, observer observability.Observer) (*RedisClient, error) {
	dsn, err := conf.DSN()

	if err != nil {
		return nil, err
	}

	opt, err := redis.ParseURL(dsn)

	if err != nil {
		return nil, err
	}

	opt.DialerRetries = 3
	opt.DialTimeout = 50 * time.Millisecond
	opt.ReadTimeout = 100 * time.Millisecond
	opt.WriteTimeout = 100 * time.Millisecond

	rdb, err := observability.NewInstrumentedRedis(opt, observer)

	if err != nil {
		return nil, err
	}

	return NewRedisClient(rdb, observer), err
}

func NewRedisClient(backend RedisBackend, observer observability.Observer) *RedisClient {
	return &RedisClient{backend: backend, logger: observer.Logger().With("client", "redis")}
}

// Get retrieves the cached value associated with the given key.
// If the key exists, the raw cached bytes are returned. If the key
// does not exist returns error.
//
// Returns a mapped cache error for consistent error handling.
func (p RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := p.backend.Get(ctx, key).Bytes()

	if err != nil {
		return []byte{}, mapCacheError(err)
	}

	return data, nil
}

// Set stores the given value in the cache under the provided key with
// the specified TTL. The operation is performed using a set-if-not-exists
// strategy to avoid overwriting existing entries.
//
// Returns a mapped cache error for consistent error handling.
func (p RedisClient) Set(ctx context.Context, value any, key string, ttl time.Duration) error {
	err := p.backend.SetNX(ctx, key, value, ttl).Err()

	if err != nil {
		return mapCacheError(err)
	}

	return nil
}

// Set stores the given value in the cache under the provided key with
// the specified TTL. The operation is performed using a set-if-not-exists
// strategy to avoid overwriting existing entries.
//
// Returns a mapped cache error for consistent error handling.
func (p RedisClient) Del(ctx context.Context, key string) error {
	err := p.backend.Del(ctx, key).Err()

	if err != nil {
		return mapCacheError(err)
	}

	return nil
}

// Incr atomically increments the integer value stored at the given key
// and returns the updated value. If the key does not exist, it is
// initialized before being incremented.
//
// Returns a mapped cache error for consistent error handling.
func (p RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	current, err := p.backend.Incr(ctx, key).Result()

	if err != nil {
		return 0, mapCacheError(err)
	}

	return current, nil
}

// Incr atomically increments the integer value stored at the given key
// and returns the updated value. If the key does not exist, it is
// initialized before being incremented.
//
// Returns a mapped cache error for consistent error handling.
func (p RedisClient) Close() error {
	return mapCacheError(p.backend.Close())
}

// Ping checks the connectivity and responsiveness of the cache backend.
//
// It sends a lightweight ping command to Redis to verify that the
// connection is alive and ready to accept requests.
//
// Returns a mapped cache error for consistent error handling.
func (p RedisClient) Ping(ctx context.Context) error {
	err := p.backend.Ping(ctx).Err()

	if err != nil {
		return mapCacheError(err)
	}

	return nil
}
