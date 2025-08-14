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

func SearchAPIHandler(idx *NGramIndex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if len(query) < idx.n {
			http.Error(w, "query too short", http.StatusBadRequest)
			return
		}
		caseSensitive := r.URL.Query().Get("case") == "1"
		useRegex := r.URL.Query().Get("regex") == "1"

		var re *regexp.Regexp
		var err error
		if useRegex {
			re, err = regexp.Compile(query)
			if err != nil {
				http.Error(w, "invalid regex: "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		candidates, err := idx.SearchNGram(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var results []SearchMatch
		for _, meta := range candidates {
			absPath := meta.Path
			// If Path is not absolute, try to resolve (customize as needed)
			if !filepath.IsAbs(absPath) {
				absPath, _ = filepath.Abs(absPath)
			}
			f, err := os.Open(absPath)
			if err != nil {
				continue
			}
			scanner := bufio.NewScanner(f)
			lineNum := 0
			for scanner.Scan() {
				line := scanner.Text()
				lineNum++
				var matchRanges [][2]int
				if useRegex && re != nil {
					locs := re.FindAllStringIndex(line, -1)
					for _, loc := range locs {
						matchRanges = append(matchRanges, [2]int{loc[0], loc[1]})
					}
				} else {
					// Substring search (case-insensitive by default)
					hay := line
					needle := query
					if !caseSensitive {
						hay = strings.ToLower(line)
						needle = strings.ToLower(query)
					}
					idx := 0
					for {
						i := strings.Index(hay[idx:], needle)
						if i == -1 {
							break
						}
						matchRanges = append(matchRanges, [2]int{idx + i, idx + i + len(needle)})
						idx += i + len(needle)
					}
				}
				if len(matchRanges) > 0 {
					results = append(results, SearchMatch{
						Repo: meta.Repo,
						Path: meta.Path,
						Language: meta.Language,
						Line: line,
						LineNumber: lineNum,
						MatchRanges: matchRanges,
					})
				}
			}
			f.Close()
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
