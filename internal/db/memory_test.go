package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
)

func TestMemoryClient(t *testing.T) {
	type Row struct {
		Name string `db:"name"`
	}

	t.Run("get by default should bypass cache", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id = $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego")

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.Memory().Get(context.Background(), &Row{}, query, 1)

		assert.NoError(t, err)
		assert.True(t, fake.MemoryMetric.LastMemoryBypass)
	})

	t.Run("get by default should bypass cache with policy", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id = $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego")

		ctx := cache.WithCachePolicy(
			context.Background(),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "get-policy-key",
			},
		)

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.Memory().Get(ctx, &Row{}, query, 1)

		assert.NoError(t, err)
		assert.True(t, fake.MemoryMetric.LastMemoryBypass)
	})

	t.Run("get should cache", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		ctx := cache.WithCachePolicy(
			cache.WithCache(context.Background()),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "get-policy-key",
			},
		)

		fake.CacheBackend.Value = `{"name": "diego"}`
		err := fake.Memory().Get(ctx, &Row{}, "SELECT * FROM anything WHERE id = $1", 1)

		assert.NoError(t, err)
		assert.NotNil(t, fake.MemoryMetric.LastMemoryHitLatency)
		assert.Equal(t, "get-policy-key", fake.MemoryMetric.LastMemoryHitKey)
	})

	t.Run("get should invalidate the cache when the value is invalid.", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id = $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego")

		ctx := cache.WithCachePolicy(
			cache.WithCache(context.Background()),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "select-policy-key",
			},
		)

		fake.CacheBackend.Value = "[}"
		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.Memory().Get(ctx, &Row{}, query, 1)

		assert.NoError(t, err)
		assert.True(t, fake.MemoryMetric.LastMemoryInvalid)
	})

	t.Run("get should fetch from DB when cache does not exists", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id = $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego")

		ctx := cache.WithCachePolicy(
			cache.WithCache(context.Background()),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "select-policy-key",
			},
		)

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		fake.CacheBackend.Err = redis.Nil

		err := fake.Memory().Get(ctx, &Row{}, query, 1)

		assert.NoError(t, err)
		assert.NotNil(t, fake.MemoryMetric.LastMemoryMissLatency)
		assert.Equal(t, "select-policy-key", fake.MemoryMetric.LastMemoryMissKey)
	})

	t.Run("select by default should bypass cache", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id > $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego").AddRow("maria")

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.Memory().Select(context.Background(), &[]Row{}, query, 1)

		assert.NoError(t, err)
		assert.True(t, fake.MemoryMetric.LastMemoryBypass)
	})

	t.Run("select by default should bypass cache with policy", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id > $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego").AddRow("maria")

		ctx := cache.WithCachePolicy(
			context.Background(),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "select-policy-key",
			},
		)

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.Memory().Select(ctx, &[]Row{}, "SELECT * FROM anything WHERE id > $1", 1)

		assert.NoError(t, err)
		assert.True(t, fake.MemoryMetric.LastMemoryBypass)
	})

	t.Run("select should cache", func(t *testing.T) {
		fake := test.NewFakeDependencies()

		ctx := cache.WithCachePolicy(
			cache.WithCache(context.Background()),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "select-policy-key",
			},
		)

		fake.CacheBackend.Value = "[]"
		err := fake.Memory().Select(ctx, &[]Row{}, "SELECT * FROM anything WHERE id > $1", 1)

		assert.NoError(t, err)
		assert.NotNil(t, fake.MemoryMetric.LastMemoryHitLatency)
		assert.Equal(t, "select-policy-key", fake.MemoryMetric.LastMemoryHitKey)
	})

	t.Run("select should invalidate the cache when the value is invalid.", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id > $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego").AddRow("maria")

		ctx := cache.WithCachePolicy(
			cache.WithCache(context.Background()),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "select-policy-key",
			},
		)

		fake.CacheBackend.Value = "[}"
		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		err := fake.Memory().Select(ctx, &[]Row{}, query, 1)

		assert.NoError(t, err)
		assert.True(t, fake.MemoryMetric.LastMemoryInvalid)
	})

	t.Run("select should fetch from DB when cache does not exists", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		query := "SELECT * FROM anything WHERE id > $1"
		rows := sqlmock.NewRows([]string{"name"}).AddRow("diego").AddRow("maria")

		ctx := cache.WithCachePolicy(
			cache.WithCache(context.Background()),
			cache.CachePolicy{
				TTL: 1 * time.Minute,
				Key: "select-policy-key",
			},
		)

		fake.DBMock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
		fake.CacheBackend.Err = redis.Nil

		err := fake.Memory().Select(ctx, &[]Row{}, query, 1)

		assert.NoError(t, err)
		assert.NotNil(t, fake.MemoryMetric.LastMemoryMissLatency)
		assert.Equal(t, "select-policy-key", fake.MemoryMetric.LastMemoryMissKey)
	})
}
