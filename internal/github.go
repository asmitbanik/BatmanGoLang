// RepoResult represents a similar repo with metadata
type RepoResult struct {
	Repo   string `json:"repo"`
	Stars  int    `json:"stars"`
	URL    string `json:"url"`
}

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitHubClient struct {
	httpClient *http.Client
	token      string
}

func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	return &GitHubClient{
		httpClient: &http.Client{Timeout: 8 * time.Second},
		token:      token,
	}
}

type githubSearchResponse struct {
	Items []struct {
		Repository struct {
			HTMLURL string `json:"html_url"`
			FullName string `json:"full_name"`
		} `json:"repository"`
		HTMLURL string `json:"html_url"`
	} `json:"items"`
}


// SearchCodeSmart searches GitHub for code using both code, language, and framework for best relevance
// Returns []RepoResult for richer metadata
func (gh *GitHubClient) SearchCodeSmart(ctx context.Context, code, language, framework string) ([]RepoResult, error) {
	// Use first 5 non-empty lines as query
	lines := strings.Split(code, "\n")
	var queryLines []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "//") && !strings.HasPrefix(l, "#") {
			queryLines = append(queryLines, l)
		}
		if len(queryLines) >= 5 {
			break
		}
	}
	q := strings.Join(queryLines, " ")
	q = strings.ReplaceAll(q, "\"", "") // remove quotes
	q = strings.ReplaceAll(q, "'", "")
	if len(q) > 200 {
		q = q[:200]
	}
	// Add language and framework to query if available
	if language != "" && language != "Unknown" {
		q += " language:" + language
	}
	if framework != "" && framework != "Unknown" {
		q += " " + framework
	}
	// Build GitHub search API URL
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=10", urlQueryEscape(q))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if gh.token != "" {
		req.Header.Set("Authorization", "token "+gh.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := gh.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s: %s", resp.Status, string(body))
	}
	var ghResp githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, err
	}
	repoSet := make(map[string]RepoResult)
	var repos []RepoResult
	for _, item := range ghResp.Items {
		url := item.Repository.HTMLURL
		if url == "" {
			continue
		}
		if _, exists := repoSet[url]; !exists {
			// TODO: Fetch stars, license, last commit for each repo (stubbed)
			repo := RepoResult{
				Repo: item.Repository.FullName,
				Stars: 0, // TODO: fetch real stars
				URL: url,
			}
			repoSet[url] = repo
			repos = append(repos, repo)
		}
		if len(repos) >= 10 {
			break
		}
	}
	return repos, nil
}

// urlQueryEscape escapes a string for use in a URL query
func urlQueryEscape(s string) string {
	s = strings.ReplaceAll(s, " ", "+")
	return s
}
