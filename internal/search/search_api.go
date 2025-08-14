package search

import (
	"encoding/json"
	"net/http"
	"strings"
)

// SearchAPIHandler provides a fast /search endpoint using the n-gram index
func SearchAPIHandler(idx *NGramIndex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if len(query) < idx.n {
			http.Error(w, "query too short", http.StatusBadRequest)
			return
		}
		candidates, err := idx.SearchNGram(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Optionally: scan file contents for exact match, rank, etc.
		var results []FileMeta
		for _, meta := range candidates {
			// For demo, just return all candidates
			if strings.Contains(strings.ToLower(meta.Path), strings.ToLower(query)) {
				results = append(results, meta)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
