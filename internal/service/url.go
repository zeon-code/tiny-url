package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/zeon-code/tiny-url/internal/model"
	"github.com/zeon-code/tiny-url/internal/pkg/cache"
	"github.com/zeon-code/tiny-url/internal/repository"
)

type URLService interface {
	Create(ctx context.Context, target string) (*model.URL, error)
	List(ctx context.Context, limit int, direction string, cursor *int64) ([]model.URL, error)
	GetByID(ctx context.Context, id int64) (*model.URL, error)
}

type UrlSvc struct {
	repo     repository.URLRepository
	cacheKey cache.CacheKey
	logger   *slog.Logger
}

func NewUrlService(repositories repository.Repositories, logger *slog.Logger) URLService {
	return UrlSvc{
		repo:     repositories.Url,
		cacheKey: cache.NewCacheKey("url", "service"),
		logger:   logger,
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
