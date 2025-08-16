# handlers.go

Defines HTTP handler functions for the backend API, especially the `/analyze` endpoint.

## Key Endpoint: `/analyze` (POST)
- Accepts JSON: `{ filename, code, language? }`
- Modular pipeline:
	1. **Language detection:** Uses static and ML-based methods to detect the programming language, with confidence score.
	2. **Framework detection:** Uses static rules and ML (Naive Bayes) to detect frameworks/libraries.
	3. **Security analysis:** Calls static analyzer to find deprecated APIs, insecure patterns, and CVEs.
	4. **Purpose guess:** Uses rules/ML to guess what the code does (e.g., API handler, CLI, data science, etc.).
	5. **Complexity & style:** Computes cyclomatic complexity and style score.
	6. **Similar code search:** Calls GitHub client to find similar code/repos, with caching and timeouts.
- Returns a rich JSON response with all the above info, plus error handling and logging.

## Design Highlights
- **Separation of concerns:** Each analysis step is handled by a dedicated function in `internal/`.
- **Timeouts:** Uses context timeouts to avoid slow GitHub API calls blocking the server.
- **Extensible:** Easy to add more analysis steps or enrich the response.

## Why this matters
This file shows how to build a robust, extensible API handler in Go, with a clear, modular pipeline for code analysis. It demonstrates how to combine static analysis, ML, and external API calls in a single endpoint.
