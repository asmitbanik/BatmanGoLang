package tests

import (
    "context"
    "testing"
    "os"
    "time"

    // "github.com/rohanair/shazam-for-code/internal"
    "go.uber.org/zap"
)

func TestGitHubSearch(t *testing.T) {
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        t.Skip("no github token, skipping integration test")
    }
    logger := zap.NewExample()
    gh := internal.NewGitHubClient(token, nil, logger)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    repos, err := gh.SearchCode(ctx, "fmt.Println(\"Hello")")
    if err != nil {
        t.Fatalf("search error: %v", err)
    }
    if len(repos) == 0 {
        t.Fatalf("expected some repos")
    }
}
