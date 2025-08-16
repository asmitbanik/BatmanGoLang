
# <img src="frontend/public/logo.svg" alt="BatmanGoLang Logo" width="64" height="64" style="vertical-align:middle;">  
# Batman GoLang Code Intelligence Platform

> Instantly analyze, search, and understand code across your entire codebase with blazing speed and modern UX.

---

## Overview
Batman GoLang is a full-stack code intelligence platform inspired by tools like grep.app and Sourcegraph. It enables:
- **Instant code analysis** (language, framework, security, purpose)
- **Global code search** (grep.app-style, with n-gram index)
- **GitHub similarity search** (find related public repos)
- **Modern, minimal UI** (React + Vite)

This project is designed for developers, teams, and interviewers who want to:
- Quickly understand unfamiliar code
- Search large codebases in milliseconds
- Find similar code and best practices from open source
- Demo code intelligence in interviews or onboarding

---

## Key Features

- **Live code analysis:**
	- Paste or upload code, auto-detect language and framework
	- Security warnings, complexity, and purpose guess
	- Similar repo suggestions from GitHub

- **Global code search:**
	- Grep.app-style search bar with live, ranked results
	- Fast n-gram index for substring and regex search
	- Filters: repo, language, file path, etc.
	- Pagination, keyboard navigation, and code highlighting

- **GitHub integration:**
	- Find similar code and related repositories
	- Toggle between local and GitHub search

- **Modern UI:**
	- React (Vite) frontend, instant feedback
	- Clean, minimal, mobile-friendly design
	- Smooth navigation between editor and search

---

## User Flow

### 1. Main Page
- Paste or upload code, enter filename, and click **Analyze**
- Instantly see:
	- Detected language and confidence
	- Framework and confidence
	- Security warnings (CVE, deprecated APIs)
	- Purpose guess (ML/LLM-based)
	- Complexity metrics (cyclomatic, style score)
	- Similar GitHub repos (with stars)
- Click **Global Code Search** to jump to the search page

### 2. Search Page
- Large, centered search bar (like grep.app)
- Type to see live, ranked results (file, repo, line, highlights)
- Keyboard navigation (up/down/enter/esc)
- Click a result to open file/line (customizable)
- Toggle between local index and GitHub search
- Click **Back** to return to the main editor

---

## Architecture

### Backend (Go)
- **Gin API** (`main.go`): Entry point, config, middleware
- **API Handlers** (`api/`): `/analyze`, `/search`, etc.
- **Core Logic** (`internal/`):
	- Language/framework detection (static + ML)
	- Security analysis, complexity, purpose guess
	- GitHub smart search (caching, filtering)
	- **NGramIndex**: Fast, persistent n-gram index for code search
	- Modular interfaces: `Indexer`, `Searcher` (extensible)
- **Config** (`config/`): Environment, Redis, etc.

### Frontend (React + Vite)
- **App.tsx**: Main code editor and analysis UI
- **SearchPage.tsx**: Grep.app-style global code search UI
- **RootRouter.tsx**: Routing between editor and search
- **API Layer**: `searchApi.ts` for backend integration
- **Components**: Code editor, instant search, result highlighting

### Documentation
- Each folder contains a `docs/` directory with per-file documentation
- See `backend/internal/search/docs/` for search/index details
- See `frontend/docs/` for UI and API docs

---

## API Endpoints

- `POST /analyze` — Analyze code (language, framework, security, etc.)
- `GET /search` — Search codebase (supports `q`, `repo`, `language`, `path`, `limit`, `offset`)
	- Returns: `{ results: [...], totalCount, hasMore }`

---

## Extensibility & Customization

- **Add new analyzers**: Plug in new ML models or static analyzers in `internal/`
- **Swap search backend**: Implement `Indexer`/`Searcher` for Bleve, Meilisearch, etc.
- **UI themes**: Customize styles in `frontend/src/styles.css`
- **Integrate with CI/CD**: Use API endpoints for automated code review

---

## Quickstart

### Prerequisites
- Go 1.20+
- Bun or Node.js 18+
- Redis (optional, for caching)
- GitHub API token (recommended)

### Backend
```sh
go mod tidy
$env:GITHUB_TOKEN="your_token_here" # PowerShell
go run main.go
```

### Frontend
```sh
cd frontend
bun install # or npm install
bun run dev # or npm run dev
```

### Environment Variables
- `GITHUB_TOKEN` — GitHub API token (required for production)
- `REDIS_ADDR` — Redis address (default: localhost:6379)

---

## Contribution Guide

1. Fork and clone the repo
2. Create a new branch for your feature or fix
3. Add/modify code and update docs as needed
4. Run tests and lint (if available)
5. Open a pull request with a clear description

---

## License
MIT
