# Shazam for Code — Final Implementation

## Objective
Upload any code snippet → Detect programming language, framework (ML-based), and related GitHub repos.

## Structure
- `main.go` — Entry point, server setup
- `api/handlers.go` — API endpoint for `/analyze`
- `internal/` — Core logic: language/framework detection, ML, GitHub search, caching, logging, middleware
- `config/` — Config loader
- `frontend/` — React + Vite frontend (Monaco editor, API integration)

## How to Run

### Backend
1. `cd "new implementation"`
2. `go mod tidy`
3. `go run main.go`

### Frontend
1. `cd frontend`
2. `npm install`
3. `npm run dev`

### Usage
- Open `http://localhost:3000` in your browser
- Paste/upload code, set filename, click Analyze
- See detected language, framework, and similar GitHub repos

---

This folder contains only the final, non-duplicated, production-ready implementation for the Shazam for Code objective.
