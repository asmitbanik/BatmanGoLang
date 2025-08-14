# Shazam for Code

> Upload any code snippet → Detect programming language, framework, and related GitHub repos

## Features
- Paste or upload any code snippet
- Detects programming language (using enry)
- Detects framework (ML-based + keyword/regex)
- Finds similar code on GitHub and returns related repositories
- Fast, modern React frontend (Vite)

## Project Structure
- `main.go` — Go backend (Gin API)
- `internal/` — Core logic (detection, ML, GitHub, cache, etc.)
- `frontend/` — React (Vite) frontend

## Getting Started

### Prerequisites
- Go 1.20+
- Node.js 18+ (for frontend)
- Redis (optional, for caching)
- GitHub API token (recommended, for higher rate limits)

### Backend
```sh
go mod tidy
export GITHUB_TOKEN=your_token_here
go run main.go
```

### Frontend
```sh
cd frontend
npm install # or bun install
npm run dev # or bun run dev
```

### Usage
1. Open the frontend in your browser (default: http://localhost:5173)
2. Paste or upload code, enter filename, click Analyze
3. See detected language, framework, and similar GitHub repos

## Environment Variables
- `GITHUB_TOKEN` — GitHub API token (required for production)
- `REDIS_ADDR` — Redis address (default: localhost:6379)

## License
MIT
