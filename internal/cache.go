package internal

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v9"
    "github.com/patrickmn/go-cache"
)

// CacheClient abstracts caching implementation
type CacheClient interface {
    Get(key string) (string, error)
    Set(key, val string, ttl time.Duration) error
}

// RedisCache implementation
type RedisCache struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
    // accept single URL like redis://:password@host:port/0
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }
    client := redis.NewClient(opt)
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return &RedisCache{client: client, ctx: ctx}, nil
}

func (r *RedisCache) Get(key string) (string, error) {
    return r.client.Get(r.ctx, key).Result()
}

func (r *RedisCache) Set(key, val string, ttl time.Duration) error {
    return r.client.Set(r.ctx, key, val, ttl).Err()
}

// InMemoryCache fallback
type InMemoryCache struct {
    c *cache.Cache
}

func NewInMemoryCache(defaultTTL time.Duration) *InMemoryCache {
    return &InMemoryCache{c: cache.New(defaultTTL, 2*defaultTTL)}
}

func (m *InMemoryCache) Get(key string) (string, error) {
    if v, ok := m.c.Get(key); ok {
        return v.(string), nil
    }
    return "", fmt.Errorf("cache miss")
}

func (m *InMemoryCache) Set(key, val string, ttl time.Duration) error {
    m.c.Set(key, val, ttl)
    return nil
}

// NewCache initializes either redis cache or in-memory cache
func NewCache(cfg *config.Config) (CacheClient, error) {
    if cfg.RedisURL != "" {
        r, err := NewRedisCache(cfg.RedisURL)
        if err == nil {
            return r, nil
        }
    }
    // fallback
    return NewInMemoryCache(5 * time.Minute), nil
}
