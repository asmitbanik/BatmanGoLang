# Shazam for Code — Production-ready Go Backend

## Overview
This service analyzes a code snippet and returns detected language, guessed framework, and similar repositories from GitHub.

## Quickstart (local)

1. Copy `.env.example` to `.env` and set `GITHUB_TOKEN`.
2. Start with Docker Compose:

```bash
docker-compose up --build
```

3. Or run locally:

```bash
export GITHUB_TOKEN=ghp_xxx
go run main.go
```

## API

POST `/analyze`

Request body:

```json
{
  "filename": "main.go",
  "code": "package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"Hello\") }"
}
```

Response example:

```json
{
  "language": "Go",
  "framework": "Gin",
  "similar_repos": ["https://github.com/user/repo1", "https://github.com/user/repo2"]
}
```

## Production notes

* Use Redis; set `REDIS_URL`.
* Run with TLS behind a reverse proxy.
* Use Kubernetes manifests to scale.
