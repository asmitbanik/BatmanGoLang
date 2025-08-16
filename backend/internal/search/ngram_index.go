// Indexer defines methods for indexing files
type Indexer interface {
	IndexFile(meta FileMeta, content []byte) error
	RemoveFile(path string) error
	ReindexAll() error
}

// SearchFilters for advanced search
type SearchFilters struct {
	Repo     string
	Language string
	Path     string
}

// Searcher defines methods for searching the index
type Searcher interface {
	Search(query string, filters SearchFilters, limit, offset int) ([]SearchMatch, int, error)
}

// Ensure NGramIndex implements Indexer and Searcher
var _ Indexer = (*NGramIndex)(nil)
var _ Searcher = (*NGramIndex)(nil)

// RemoveFile removes a file from the index (stub for now)
func (idx *NGramIndex) RemoveFile(path string) error {
	// TODO: implement removal logic
	return nil
}

// ReindexAll reindexes all files (stub for now)
func (idx *NGramIndex) ReindexAll() error {
	// TODO: implement full reindex
	return nil
}

// Search with filters, pagination, and ranking
func (idx *NGramIndex) Search(query string, filters SearchFilters, limit, offset int) ([]SearchMatch, int, error) {
	candidates, err := idx.SearchNGram(query)
	if err != nil {
		return nil, 0, err
	}
	var results []SearchMatch
	for _, meta := range candidates {
		// Filtering
		if filters.Repo != "" && meta.Repo != filters.Repo {
			continue
		}
		if filters.Language != "" && meta.Language != filters.Language {
			continue
		}
		if filters.Path != "" && !strings.Contains(meta.Path, filters.Path) {
			continue
		}
		absPath := meta.Path
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
			hay := line
			needle := query
			hay = strings.ToLower(line)
			needle = strings.ToLower(query)
			idxPos := 0
			for {
				i := strings.Index(hay[idxPos:], needle)
				if i == -1 {
					break
				}
				matchRanges = append(matchRanges, [2]int{idxPos + i, idxPos + i + len(needle)})
				idxPos += i + len(needle)
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
	// Ranking: sort by number of matches per file (descending)
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].MatchRanges) > len(results[j].MatchRanges)
	})
	total := len(results)
	// Pagination
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}
	return results[start:end], total, nil
}
package search

import (
	"bytes"
	"strings"
	"sync"
	"github.com/dgraph-io/badger/v3"
)

// NGramIndexer builds and queries n-gram (trigram) indexes for fast substring search
// This is the core of grep.app-style search

type FileMeta struct {
	Repo      string
	Path      string
	Language  string
	UpdatedAt int64
}

type NGramIndex struct {
	db   *badger.DB
	lock sync.RWMutex
	n    int // n-gram size (e.g., 3 for trigrams)
}

// NewNGramIndex opens or creates a BadgerDB-backed n-gram index
func NewNGramIndex(path string, n int) (*NGramIndex, error) {
	db, err := badger.Open(badger.DefaultOptions(path).WithLogger(nil))
	if err != nil {
		return nil, err
	}
	return &NGramIndex{db: db, n: n}, nil
}

// IndexFile indexes a file's content by extracting all n-grams and storing their positions
func (idx *NGramIndex) IndexFile(meta FileMeta, content []byte) error {
	idx.lock.Lock()
	defer idx.lock.Unlock()
	ngrams := extractNGrams(content, idx.n)
	return idx.db.Update(func(txn *badger.Txn) error {
		for ng, positions := range ngrams {
			key := []byte("ng:" + ng)
			val := encodeFileRef(meta, positions)
			if err := txn.Set(key, val); err != nil {
				return err
			}
		}
		return nil
	})
}

// SearchNGram returns candidate files containing all n-grams in the query
func (idx *NGramIndex) SearchNGram(query string) ([]FileMeta, error) {
	idx.lock.RLock()
	defer idx.lock.RUnlock()
	ngrams := extractNGrams([]byte(query), idx.n)
	candidates := make(map[string]FileMeta)
	for ng := range ngrams {
		err := idx.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			key := []byte("ng:" + ng)
			item, err := txn.Get(key)
			if err == nil {
				val, _ := item.ValueCopy(nil)
				meta := decodeFileRef(val)
				candidates[meta.Repo+meta.Path] = meta
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	var out []FileMeta
	for _, m := range candidates {
		out = append(out, m)
	}
	return out, nil
}

// extractNGrams returns a map of n-gram to positions in the content
func extractNGrams(content []byte, n int) map[string][]int {
	ngrams := make(map[string][]int)
	for i := 0; i <= len(content)-n; i++ {
		ng := string(bytes.ToLower(content[i : i+n]))
		ngrams[ng] = append(ngrams[ng], i)
	}
	return ngrams
}

// encodeFileRef and decodeFileRef are stubs for serializing file metadata and positions
func encodeFileRef(meta FileMeta, positions []int) []byte {
	// For demo: just store path (in production, use protobuf/gob/json)
	return []byte(meta.Repo + ":" + meta.Path)
}

func decodeFileRef(data []byte) FileMeta {
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) == 2 {
		return FileMeta{Repo: parts[0], Path: parts[1]}
	}
	return FileMeta{Path: string(data)}
}
