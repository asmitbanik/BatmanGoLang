# github.go

Handles all GitHub API interactions, especially code search for similar repositories.

## Key Concepts
- **GitHubClient struct:** Wraps an HTTP client and stores the API token.
- **SearchCodeSmart:**
	- Builds a search query from the user's code (first 5 non-empty lines, language, framework)
	- Calls the GitHub Search API for code, parses results, and deduplicates repos
	- Returns a list of `RepoResult` (repo name, stars, URL)
	- Handles API errors, timeouts, and rate limits
- **Extensible:** Designed to fetch more metadata (stars, license, last commit) in the future

## Why this matters
This file is the bridge between your backend and the global open-source ecosystem. It enables the "find similar code/repos" feature, which is a key differentiator for the app.
