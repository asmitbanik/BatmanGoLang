# cache.go

Implements caching logic for the backend, with both in-memory and Redis-backed options.

## Key Concepts
- **Cache interface:** Defines `Get(key)`, `Set(key, value, ttl)` for pluggable cache backends.
- **InMemoryCache:** Thread-safe map with expiration, used as a fallback if Redis is unavailable.
- **RedisCache:** Uses `go-redis` to store cache entries in Redis, supports expiration.
- **NewCache():** Returns a Redis cache if available, else falls back to in-memory cache (auto-detects at startup).

## Why this matters
Caching is critical for performance, especially for expensive operations like GitHub API calls and code search. This design allows seamless switching between local and distributed cache, with no code changes needed elsewhere.
