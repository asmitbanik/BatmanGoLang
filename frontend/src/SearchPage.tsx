import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

type CodeSearchResult = {
  filename: string;
  repo: string;
  line: number;
  snippet: string;
  url: string;
};

// Dummy search function (replace with real backend call)
async function searchCode(query: string, source: string): Promise<CodeSearchResult[]> {
  // source: 'github' | 'local'
  // TODO: Call your backend endpoint with ?q=query&source=github|local
  // Return format: [{ filename, repo, line, snippet, url }]
  return [];
}

export default function CodeSearchPage() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<CodeSearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showPopup, setShowPopup] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [source, setSource] = useState<'github' | 'local'>('github');
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
        const res = await searchCode(query, source);
        setResults(res || []);
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
          window.open(results[selectedIndex].url, '_blank');
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
                key={r.url}
                style={{
                  padding: '12px 16px',
                  background: i === selectedIndex ? '#f5f7fa' : undefined,
                  cursor: 'pointer'
                }}
                onMouseEnter={() => setSelectedIndex(i)}
                onClick={() => window.open(r.url, '_blank')}
              >
                <div style={{ fontSize: 14, color: '#555' }}>
                  <strong>{r.filename}</strong> in <span style={{ color: '#888' }}>{r.repo}</span> — line {r.line}
                </div>
                <pre style={{ margin: '4px 0 0 0', fontSize: 13, background: 'none', whiteSpace: 'pre-wrap' }}>
                  {highlightMatch(r.snippet, query)}
                </pre>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function highlightMatch(text: string, query: string) {
  if (!query) return text;
  const regex = new RegExp(`(${escapeRegExp(query)})`, 'gi');
  const parts = text.split(regex);
  return parts.map((part, i) =>
    regex.test(part) ? <mark key={i} style={{ background: '#ffe066' }}>{part}</mark> : part
  );
}
function escapeRegExp(string: string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
