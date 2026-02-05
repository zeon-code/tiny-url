package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/pkg/test"
	"github.com/zeon-code/tiny-url/internal/service"
)

func TestUrlService(t *testing.T) {
	ctx := context.Background()

	t.Run("create url", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		svc := service.NewUrlService(fake.Repositories(), fake.Logger())

		fake.MockUrlCreate()
		url, err := svc.Create(ctx, "target")

		assert.NoError(t, err)
		assert.Equal(t, model.URL{ID: 1, Code: "1", Target: "target", CreatedAt: url.CreatedAt, UpdatedAt: url.UpdatedAt}, *url)
	})

	t.Run("list url", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		svc := service.NewUrlService(fake.Repositories(), fake.Logger())

		fake.MockUrlList()
		urls, err := svc.List(ctx, 5, ">", nil)

		assert.NoError(t, err)
		assert.Len(t, urls, 5)
	})

	t.Run("list paginated url", func(t *testing.T) {
		cursor := int64(1)
		fake := test.NewFakeDependencies()
		svc := service.NewUrlService(fake.Repositories(), fake.Logger())

		fake.MockPaginatedUrlList()
		urls, err := svc.List(ctx, 5, ">", &cursor)

		assert.NoError(t, err)
		assert.Len(t, urls, 5)
	})

	t.Run("list url from cache", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		svc := service.NewUrlService(fake.Repositories(), fake.Logger())

		fake.CacheBackend.Value = `[]`
		urls, err := svc.List(cache.WithCache(ctx), 5, ">", nil)

		assert.NoError(t, err)
		assert.Len(t, urls, 0)
	})

	t.Run("url get by id", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		svc := service.NewUrlService(fake.Repositories(), fake.Logger())

		fake.MockUrlGetById()
		url, err := svc.GetByID(ctx, int64(1))

		assert.NoError(t, err)
		assert.Equal(t, model.URL{ID: 1, Code: "1", Target: "target1", CreatedAt: url.CreatedAt, UpdatedAt: url.UpdatedAt}, *url)
	})

	t.Run("url get by id from cache", func(t *testing.T) {
		fake := test.NewFakeDependencies()
		svc := service.NewUrlService(fake.Repositories(), fake.Logger())

		fake.CacheBackend.Value = `{"id": 1, "Code": "1", "target": "target1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}`
		url, err := svc.GetByID(cache.WithCache(ctx), int64(1))

		assert.NoError(t, err)
		assert.Equal(t, model.URL{ID: 1, Code: "1", Target: "target1", CreatedAt: url.CreatedAt, UpdatedAt: url.UpdatedAt}, *url)
	})
}
