package internal

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string, ttl time.Duration)
}

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value      string
	expiresAt  time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{items: make(map[string]cacheItem)}
}

func (c *InMemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		return "", false
	}
	return item.value, true
}

func (c *InMemoryCache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{value: value, expiresAt: time.Now().Add(ttl)}
}

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache() *RedisCache {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisCache{client: client, ctx: context.Background()}
}

func (c *RedisCache) Get(key string) (string, bool) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func (c *RedisCache) Set(key, value string, ttl time.Duration) {
	c.client.Set(c.ctx, key, value, ttl)
}

// NewCache returns a Cache, preferring Redis if available
func NewCache() Cache {
	cache := NewRedisCache()
	testKey := "cache_test_key"
	cache.Set(testKey, "ok", 2*time.Second)
	if v, ok := cache.Get(testKey); ok && v == "ok" {
		return cache
	}
	return NewInMemoryCache()
}
