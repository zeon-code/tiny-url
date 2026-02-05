package test

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type FakeRedis struct {
	Err   error
	Value any

	LastGetKey        string
	LastDelKey        []string
	LastIncrKey       string
	LastSetKey        string
	LastSetValue      any
	LastSetExpiration time.Duration
}

func (r *FakeRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	r.LastGetKey = key

	v, _ := r.Value.(string)
	return redis.NewStringResult(v, r.Err)
}

func (r *FakeRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	r.LastDelKey = keys

	v, _ := r.Value.(int64)
	return redis.NewIntResult(v, r.Err)
}

func (r *FakeRedis) Incr(ctx context.Context, key string) *redis.IntCmd {
	r.LastIncrKey = key

	v, _ := r.Value.(int64)
	return redis.NewIntResult(v, r.Err)
}

func (r *FakeRedis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	r.LastSetKey = key
	r.LastSetValue = value
	r.LastSetExpiration = expiration

	v, _ := r.Value.(bool)
	return redis.NewBoolResult(v, r.Err)
}

func NewFakeRedisBackend() *FakeRedis {
	return &FakeRedis{}
}
