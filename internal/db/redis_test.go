package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/db"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestRedisCacheClient(t *testing.T) {
	key := "key"
	ctx := context.Background()

	t.Run("proxy get command", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Value = "hit"

		buffer, err := fake.Cache().Get(ctx, key)

		assert.NoError(t, err)
		assert.Equal(t, "hit", string(buffer))
		assert.True(t, fake.CacheMetric.LastCacheHit)
	})

	t.Run("proxy get command with miss", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Err = redis.Nil

		buffer, err := fake.Cache().Get(ctx, key)

		assert.Equal(t, []byte{}, buffer)
		assert.Equal(t, db.ErrCacheNotFound, err)
		assert.True(t, fake.CacheMetric.LastCacheMiss)
	})

	t.Run("proxy get command with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Err = redis.ErrClosed

		buffer, err := fake.Cache().Get(ctx, key)

		assert.Equal(t, []byte{}, buffer)
		assert.Equal(t, db.ErrCacheUnavailable, err)
		assert.Equal(t, "failed to read redis key: redis: client is closed", fake.CacheMetric.LastCacheErr)
		assert.Equal(t, key, fake.CacheMetric.LastCacheKeyErr)
	})

	t.Run("proxy set command", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		err := fake.Cache().Set(ctx, "value", key, 1*time.Minute)

		assert.NoError(t, err)
		assert.Empty(t, fake.CacheMetric.LastCacheErr)
		assert.Empty(t, fake.CacheMetric.LastCacheKeyErr)
	})

	t.Run("proxy set command with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Err = redis.ErrClosed

		err := fake.Cache().Set(ctx, "value", key, 1*time.Minute)

		assert.Equal(t, db.ErrCacheUnavailable, err)
		assert.Equal(t, "failed to write value into key: redis: client is closed", fake.CacheMetric.LastCacheErr)
		assert.Equal(t, key, fake.CacheMetric.LastCacheKeyErr)
	})

	t.Run("proxy del command", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		err := fake.Cache().Del(ctx, key)

		assert.NoError(t, err)
		assert.Empty(t, fake.CacheMetric.LastCacheErr)
		assert.Empty(t, fake.CacheMetric.LastCacheKeyErr)
	})

	t.Run("proxy del command with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Err = redis.ErrClosed

		err := fake.Cache().Del(ctx, key)

		assert.Equal(t, db.ErrCacheUnavailable, err)
		assert.Equal(t, "failed to delete redis key: redis: client is closed", fake.CacheMetric.LastCacheErr)
		assert.Equal(t, key, fake.CacheMetric.LastCacheKeyErr)
	})

	t.Run("proxy incr command", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Value = int64(1)

		value, err := fake.Cache().Incr(ctx, key)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), value)
		assert.Empty(t, fake.CacheMetric.LastCacheErr)
		assert.Empty(t, fake.CacheMetric.LastCacheKeyErr)
	})

	t.Run("proxy incr command with error", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		fake.CacheBackend.Err = redis.ErrClosed

		_, err := fake.Cache().Incr(ctx, key)

		assert.Equal(t, db.ErrCacheUnavailable, err)
		assert.Equal(t, "failed to increment redis key: redis: client is closed", fake.CacheMetric.LastCacheErr)
		assert.Equal(t, key, fake.CacheMetric.LastCacheKeyErr)
	})
}
