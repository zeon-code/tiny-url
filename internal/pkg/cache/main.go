package cache

import (
	"context"
	"time"
)

type cacheKey struct{}

type CachePolicy struct {
	TTL time.Duration
	Key string
}

type Cache struct {
	IsEnabled bool
	Policy    CachePolicy
}

func WithCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, cacheKey{}, Cache{
		IsEnabled: true,
	})
}

func WithCachePolicy(ctx context.Context, policy CachePolicy) context.Context {
	cache, _ := ctx.Value(cacheKey{}).(Cache)
	cache.Policy = policy
	return context.WithValue(ctx, cacheKey{}, cache)
}

func CacheFromContext(ctx context.Context) Cache {
	cache, _ := ctx.Value(cacheKey{}).(Cache)
	return cache
}
