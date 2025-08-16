package search

import (
	"encoding/json"
	"net/http"
	"os"
	"bufio"
	"regexp"
	"strings"
	"path/filepath"
)

// SearchAPIHandler provides a fast /search endpoint using the n-gram index
type SearchMatch struct {
	Repo      string   `json:"repo"`
	Path      string   `json:"path"`
	Language  string   `json:"language,omitempty"`
	Line      string   `json:"line"`
	LineNumber int     `json:"lineNumber"`
	MatchRanges [][2]int `json:"matchRanges"` // [start,end) byte offsets in line
}

func SearchAPIHandler(idx Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if len(query) < 3 {
			http.Error(w, "query too short", http.StatusBadRequest)
			return
		}
		filters := SearchFilters{
			Repo:     r.URL.Query().Get("repo"),
			Language: r.URL.Query().Get("language"),
			Path:     r.URL.Query().Get("path"),
		}
		limit := 20
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			fmt.Sscanf(o, "%d", &offset)
		}
		results, total, err := idx.Search(query, filters, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": results,
			"totalCount": total,
			"hasMore": offset+limit < total,
		})
	}
}
