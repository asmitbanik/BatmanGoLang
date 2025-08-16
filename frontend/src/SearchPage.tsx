import React, { useState, useRef, useEffect } from 'react';
// Needed for JSX.Element type
import type { ReactElement } from 'react';
import { useNavigate } from 'react-router-dom';


type CodeSearchResult = {
  repo: string;
  path: string;
  language?: string;
  line: string;
  lineNumber: number;
  matchRanges: [number, number][];
};

async function searchCode(query: string, source: string, filters: {repo?: string, language?: string, path?: string} = {}, limit = 20, offset = 0): Promise<{results: CodeSearchResult[], totalCount: number, hasMore: boolean}> {
  if (!query || query.length < 3) return {results: [], totalCount: 0, hasMore: false};
  if (source === 'github') {
    // TODO: Implement GitHub search integration
    return {results: [], totalCount: 0, hasMore: false};
  }
  const params: any = { q: query, limit, offset };
  if (filters.repo) params.repo = filters.repo;
  if (filters.language) params.language = filters.language;
  if (filters.path) params.path = filters.path;
  const res = await fetch(`/search?${new URLSearchParams(params).toString()}`);
  if (!res.ok) throw new Error('Search failed');
  return await res.json();
}

export default function CodeSearchPage() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<CodeSearchResult[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showPopup, setShowPopup] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [source, setSource] = useState<'github' | 'local'>('local');
  const inputRef = useRef<HTMLInputElement>(null);
  const popupRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  // Debounced search
  useEffect(() => {
    if (!query) {
      setResults([]);
      setShowPopup(false);
      return;
    }
    setLoading(true);
    setError(null);
    const timeout = setTimeout(async () => {
      try {
        const {results: res, totalCount, hasMore} = await searchCode(query, source);
        setResults(res || []);
        setTotalCount(totalCount);
        setHasMore(hasMore);
        setShowPopup(true);
      } catch (e) {
        setError('Search failed');
      } finally {
        setLoading(false);
      }
    }, 300);
    return () => clearTimeout(timeout);
  }, [query, source]);

  // Keyboard navigation
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (!showPopup) return;
      if (e.key === 'ArrowDown') {
        setSelectedIndex(i => Math.min(i + 1, results.length - 1));
        e.preventDefault();
      } else if (e.key === 'ArrowUp') {
        setSelectedIndex(i => Math.max(i - 1, 0));
        e.preventDefault();
      } else if (e.key === 'Enter') {
        if (results[selectedIndex]) {
          const r = results[selectedIndex];
          window.open(`#file=${encodeURIComponent(r.path)}&line=${r.lineNumber}`, '_blank');
        }
      } else if (e.key === 'Escape') {
        setShowPopup(false);
      }
    }
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [showPopup, results, selectedIndex]);

  // Click outside to close
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (
        popupRef.current &&
        !(popupRef.current as HTMLDivElement).contains(e.target as Node) &&
        inputRef.current &&
        !(inputRef.current as HTMLInputElement).contains(e.target as Node)
      ) {
        setShowPopup(false);
      }
    }
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, []);

  return (
    <div style={{ minHeight: '100vh', background: '#f6f7fb' }}>
      <div style={{ maxWidth: 600, margin: '60px auto 0 auto', padding: 24 }}>
        <button onClick={() => navigate('/')} style={{ marginBottom: 32, background: 'none', border: 'none', color: '#0070f3', fontSize: 16, cursor: 'pointer' }}>&larr; Back to Editor</button>
        <h2 style={{ textAlign: 'center', marginBottom: 32 }}>Global Code Search</h2>
        <div style={{ display: 'flex', justifyContent: 'center', marginBottom: 16 }}>
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={e => setQuery(e.target.value)}
            onFocus={() => query && setShowPopup(true)}
            placeholder="Search code across all repositories..."
            style={{
              width: '100%',
              maxWidth: 480,
              padding: '18px 20px',
              fontSize: 20,
              borderRadius: 10,
              border: '1px solid #ccc',
              boxShadow: '0 2px 8px rgba(0,0,0,0.04)'
            }}
          />
          <select value={source} onChange={e => setSource(e.target.value as 'github' | 'local')} style={{ marginLeft: 12, padding: 10, fontSize: 16, borderRadius: 8 }}>
            <option value="github">GitHub</option>
            <option value="local">Local Index</option>
          </select>
        </div>
        {showPopup && (
          <div
            ref={popupRef}
            style={{
              position: 'absolute',
              left: '50%',
              transform: 'translateX(-50%)',
              width: '100%',
              maxWidth: 600,
              background: '#fff',
              border: '1px solid #eee',
              borderRadius: 10,
              boxShadow: '0 8px 32px rgba(0,0,0,0.08)',
              zIndex: 100,
              maxHeight: 400,
              overflowY: 'auto',
              marginTop: 8
            }}
          >
            {loading && <div style={{ padding: 16 }}>Searching…</div>}
            {error && <div style={{ padding: 16, color: 'red' }}>{error}</div>}
            {!loading && !error && results.length === 0 && (
              <div style={{ padding: 16, color: '#888' }}>No results found</div>
            )}
            {results.map((r, i) => (
              <div
                key={r.repo + r.path + r.lineNumber}
                style={{
                  padding: '12px 16px',
                  background: i === selectedIndex ? '#f5f7fa' : undefined,
                  cursor: 'pointer'
                }}
                onMouseEnter={() => setSelectedIndex(i)}
                onClick={() => window.open(`#file=${encodeURIComponent(r.path)}&line=${r.lineNumber}`, '_blank')}
              >
                <div style={{ fontSize: 14, color: '#555' }}>
                  <strong>{r.path.split('/').pop()}</strong> in <span style={{ color: '#888' }}>{r.repo}</span> — line {r.lineNumber}
                </div>
                <pre style={{ margin: '4px 0 0 0', fontSize: 13, background: 'none', whiteSpace: 'pre-wrap' }}>
                  {highlightMatch(r.line, query, r.matchRanges)}
                </pre>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function highlightMatch(text: string, query: string, matchRanges?: [number, number][]) {
  if (!query || !matchRanges || matchRanges.length === 0) return text;
  const out: (string | ReactElement)[] = [];
  let last = 0;
  matchRanges.forEach(([start, end], i) => {
    if (start > last) out.push(text.slice(last, start));
    out.push(<mark key={i} style={{ background: '#ffe066' }}>{text.slice(start, end)}</mark>);
    last = end;
  });
  if (last < text.length) out.push(text.slice(last));
  return out;
}
