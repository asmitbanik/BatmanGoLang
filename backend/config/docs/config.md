# config.go

Handles configuration loading for the backend. Centralizes all environment/config file parsing.

## Key Struct: `Config`
- `GitHubToken`: GitHub API token for authenticated requests
- `Port`: HTTP server port (default 8080, overridable by env)
- `RedisURL`: Redis connection string (optional, for caching)
- `Env`: Application environment (e.g., development, production)
- `RateLimitReqPerMin`: Per-IP rate limit for API requests

## How config is loaded
- Reads environment variables (with sensible defaults)
- Converts string env vars to int where needed (port, rate limit)
- Returns a pointer to a `Config` struct used throughout the backend

## Why this matters
Centralizing config makes the app easy to deploy, test, and scale. All settings are in one place, and the code is robust to missing/invalid env vars.
