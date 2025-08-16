# logger.go

Implements logging for the backend using Uber's Zap logger.

## Key Concepts
- **NewLogger:** Returns a production-ready Zap logger (JSON, stack traces, etc.)
- **GinZapMiddleware:** Gin middleware that logs every HTTP request (method, path, status, client IP)
- **Best practices:**
	- Structured logging (fields, not just strings)
	- Logs both requests and errors
	- Logger is passed to all major components for consistent logging

## Why this matters
Good logging is essential for debugging, monitoring, and production reliability. This file shows how to integrate a modern logger with a Go web server.
