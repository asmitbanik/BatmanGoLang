# main.go

This is the entry point for the backend server. It wires together all major components and manages the server lifecycle.

## Startup Flow
1. Loads environment variables (via `godotenv`)
2. Loads configuration using `config.LoadConfig()` (reads port, GitHub token, Redis, etc.)
3. Initializes the ML model for framework detection
4. Sets up logging (Zap logger, with environment-based config)
5. Initializes cache (prefers Redis, falls back to in-memory)
6. Creates a GitHub API client (with caching and logging)
7. Sets up the Gin HTTP server:
	- Logging middleware (Zap)
	- Recovery middleware (panic handling)
	- Rate limiting middleware (per-IP, configurable)
	- Registers all API routes (see `api/handlers.go`)
8. Starts the HTTP server in a goroutine
9. Waits for shutdown signal, then gracefully shuts down the server

## Key Code Concepts
- **Graceful shutdown:** Uses context and timeout to ensure the server closes cleanly.
- **Middleware:** Logging, recovery, and rate limiting are all handled as Gin middleware.
- **Dependency injection:** Passes config, logger, cache, and GitHub client to handlers for testability and modularity.

## Why this matters
This file demonstrates best practices for Go web servers: clean startup/shutdown, robust error handling, modular design, and production-ready middleware. It is the glue that connects all backend logic.
