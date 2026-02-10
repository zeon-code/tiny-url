package service

import (
	"context"
	"time"

	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/pkg/observability"
	"github.com/zeon-code/tiny-url/internal/repository"
)

type URLService interface {
	Create(context.Context, string) (*model.URL, error)
	List(context.Context, int, string, *int64) ([]model.URL, error)
	GetByID(context.Context, int64) (*model.URL, error)
}

type UrlSvc struct {
	repo     repository.URLRepository
	cacheKey cache.CacheKey
	logger   observability.Logger
}

func NewUrlService(repositories repository.Repositories, observer observability.Observer) URLService {
	return UrlSvc{
		repo:     repositories.Url,
		cacheKey: cache.NewCacheKey("url", "service"),
		logger:   observer.Logger().WithGroup("url-service"),
	}
}

func (s UrlSvc) Create(ctx context.Context, target string) (*model.URL, error) {
	return s.repo.Create(ctx, target)
}

func (s UrlSvc) List(ctx context.Context, limit int, direction string, cursor *int64) ([]model.URL, error) {
	return s.repo.List(
		cache.WithCachePolicy(
			ctx,
			cache.CachePolicy{
				TTL: 5 * time.Minute,
				Key: s.cacheKey.With("list", direction, cursor).String(),
			},
		),
		limit,
		direction,
		cursor,
	)
}

func (s UrlSvc) GetByID(ctx context.Context, id int64) (*model.URL, error) {
	return s.repo.GetByID(
		cache.WithCachePolicy(
			ctx,
			cache.CachePolicy{
				TTL: 5 * time.Minute,
				Key: s.cacheKey.With("id", id).String(),
			},
		),
		id,
	)
}
