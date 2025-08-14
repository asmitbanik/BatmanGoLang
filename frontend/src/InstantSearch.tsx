import React, { useState, useEffect } from 'react'
import { searchCode, SearchResult } from './searchApi'

export function InstantSearch({ query }: { query: string }) {
  const [results, setResults] = useState<SearchResult[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!query || query.length < 3) {
      setResults([])
      return
    }
    setLoading(true)
    const timeout = setTimeout(() => {
      searchCode(query).then(setResults).finally(() => setLoading(false))
    }, 400) // debounce
    return () => clearTimeout(timeout)
  }, [query])

  return (
    <div className="instant-search-results" style={{marginTop:8}}>
      {loading && <div>Searching...</div>}
      {!loading && results.length > 0 && (
        <ul style={{background:'#f8f8f8',padding:8,borderRadius:4}}>
          {results.map(r => (
            <li key={r.repo + r.path}>
              <span style={{fontWeight:'bold'}}>{r.repo}</span> / {r.path}
              {r.language && <span style={{color:'#888',marginLeft:8}}>{r.language}</span>}
            </li>
          ))}
        </ul>
      )}
      {!loading && query.length >= 3 && results.length === 0 && <div>No matches found.</div>}
    </div>
  )
}
