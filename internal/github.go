package internal

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"

    "go.uber.org/zap"
)

// GitHubClient wraps interactions with GitHub Search API
type GitHubClient struct {
    token string
    cache CacheClient
    logger *zap.Logger
    httpClient *http.Client
}

type codeSearchResponse struct {
    TotalCount int `json:"total_count"`
    IncompleteResults bool `json:"incomplete_results"`
    Items []struct {
        Name string `json:"name"`
        Path string `json:"path"`
        HTMLURL string `json:"html_url"`
        Repository struct {
            FullName string `json:"full_name"`
            HTMLURL string `json:"html_url"`
        } `json:"repository"`
    } `json:"items"`
}

func NewGitHubClient(token string, cache CacheClient, logger *zap.Logger) *GitHubClient {
    return &GitHubClient{
        token: token,
        cache: cache,
        logger: logger,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
}

// buildQuery creates a safe GitHub search query. We limit length to avoid rate-limits.
func buildQuery(snippet string) string {
    // quick cleanup: remove newlines and trim long content
    cleaned := strings.ReplaceAll(snippet, "\n", " ")
    cleaned = strings.ReplaceAll(cleaned, "\t", " ")
    cleaned = strings.TrimSpace(cleaned)
    if len(cleaned) > 400 {
        cleaned = cleaned[:400]
    }
    // URL-encode
    return url.QueryEscape(cleaned)
}

// SearchCode searches GitHub code and returns top repository URLs
func (g *GitHubClient) SearchCode(ctx context.Context, snippet string) ([]string, error) {
    if g.cache != nil {
        if v, err := g.cache.Get("gh_cache_" + snippet); err == nil {
            // cached JSON (comma separated)
            return strings.Split(v, ","), nil
        }
    }

    query := buildQuery(snippet)
    // We'll request the search/code endpoint and limit to 5 results
    url := fmt.Sprintf("https://api.github.com/search/code?q=%s+in:file&per_page=5", query)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    if g.token != "" {
        req.Header.Set("Authorization", "token "+g.token)
    }
    req.Header.Set("Accept", "application/vnd.github.v3+json")

    resp, err := g.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusUnauthorized {
        return nil, fmt.Errorf("github: unauthorized — check token and scopes")
    }
    if resp.StatusCode == http.StatusForbidden {
        // possibly rate-limited
        body, _ := io.ReadAll(resp.Body)
        g.logger.Sugar().Warnw("github forbidden", "status", resp.StatusCode, "body", string(body))
        return nil, fmt.Errorf("github: forbidden — rate limit or permissions")
    }
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("github search error: status=%d body=%s", resp.StatusCode, string(body))
    }

    var out codeSearchResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, err
    }

    var repos []string
    seen := map[string]bool{}
    for _, item := range out.Items {
        repoURL := item.Repository.HTMLURL
        if repoURL == "" {
            repoURL = item.HTMLURL
        }
        if !seen[repoURL] {
            repos = append(repos, repoURL)
            seen[repoURL] = true
        }
    }

    // cache result for a short period
    if g.cache != nil && len(repos) > 0 {
        _ = g.cache.Set("gh_cache_"+snippet, strings.Join(repos, ","), 10*time.Minute)
    }

    return repos, nil
}
