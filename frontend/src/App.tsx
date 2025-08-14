

import React, { useState } from 'react'
import CodeEditor from './components/CodeEditor'
import { analyzeCode } from './api'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { InstantSearch } from './InstantSearch'

const LANGUAGES = [
  'auto', 'go', 'python', 'javascript', 'typescript', 'java', 'c', 'cpp', 'rust', 'ruby', 'php', 'csharp', 'kotlin', 'swift', 'scala', 'html', 'css', 'json', 'bash', 'r', 'dart', 'elixir', 'haskell', 'perl', 'lua', 'matlab', 'objective-c', 'sql', 'xml', 'yaml', 'powershell', 'shell', 'assembly', 'fortran', 'groovy', 'julia', 'lisp', 'prolog', 'scheme', 'visualbasic', 'verilog', 'vhdl', 'coffeescript', 'fsharp', 'ocaml', 'clojure', 'erlang', 'elm', 'nim', 'crystal', 'reason', 'vala', 'zig', 'solidity', 'graphql', 'dockerfile', 'makefile', 'cmake', 'ini', 'toml', 'protobuf', 'tsx', 'jsx'
]


function App() {
  const [code, setCode] = useState<string>('package main\n\nimport "fmt"\n\nfunc main() {\n  fmt.Println("Hello, world")\n}\n')
  const [filename, setFilename] = useState('main.go')
  const [selectedLang, setSelectedLang] = useState('auto')
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<any | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')

  async function handleAnalyze() {
    setError(null)
    setLoading(true)
    setResult(null)
    try {
      const res = await analyzeCode({ filename, code, language: selectedLang !== 'auto' ? selectedLang : undefined })
      setResult(res)
    } catch (e: any) {
      setError(e?.message || 'Request failed')
    } finally {
      setLoading(false)
    }
  }


  return (
    <div className="app">
      <div className="header">
        <div>
          <h2>Shazam for Code — Live Editor</h2>
          <div>Backend: <strong>POST /analyze</strong></div>
        </div>
        <div style={{minWidth:320}}>
          <input
            type="text"
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
            placeholder="Instant codebase search (grep.app style)"
            style={{padding:8, width:'100%', borderRadius:4, border:'1px solid #ccc'}}
          />
          <InstantSearch query={searchQuery} />
        </div>
      </div>

      <div className="controls" style={{display:'flex',gap:8,alignItems:'center'}}>
        <input value={filename} onChange={(e) => setFilename(e.target.value)} style={{padding:8}} placeholder="Filename (optional)" />
        <select value={selectedLang} onChange={e => setSelectedLang(e.target.value)} style={{padding:8}}>
          {LANGUAGES.map(l => <option key={l} value={l}>{l === 'auto' ? 'Auto Detect' : l}</option>)}
        </select>
        <button onClick={handleAnalyze} disabled={loading} style={{padding:'8px 12px'}}>
          {loading ? 'Analyzing...' : 'Analyze'}
        </button>
      </div>

      <CodeEditor code={code} setCode={setCode} language={selectedLang === 'auto' ? undefined : selectedLang} />

      <div style={{height: 16}} />

      <div className="result">
        <h3>Result</h3>
        {error && <div style={{color:'red'}}>Error: {error}</div>}
        {!result && !error && <div>Click Analyze to see results</div>}
        {result && (
          <div>
            <p><strong>Language:</strong> {result.language?.name} {result.language?.confidence && <span style={{color:'#888'}}>({(result.language.confidence*100).toFixed(0)}%)</span>}</p>
            <p><strong>Framework:</strong> {result.framework?.name} {result.framework?.confidence && <span style={{color:'#888'}}>({(result.framework.confidence*100).toFixed(0)}%)</span>}</p>

            {result.security && result.security.length > 0 && (
              <div style={{margin:'8px 0'}}>
                <strong>Security Warnings:</strong>
                <ul>
                  {result.security.map((s: any, i: number) => (
                    <li key={i} style={{color: s.severity === 'high' ? 'red' : s.severity === 'medium' ? 'orange' : 'inherit'}}>
                      {s.issue} {s.cve && <span>({s.cve})</span>}
                    </li>
                  ))}
                </ul>
              </div>
            )}

            {result.purpose_guess && (
              <div style={{margin:'8px 0'}}>
                <strong>Purpose Guess:</strong> <span>{result.purpose_guess}</span>
              </div>
            )}

            {result.complexity && (
              <div style={{margin:'8px 0'}}>
                <strong>Complexity:</strong> Cyclomatic {result.complexity.cyclomatic} | Style Score {(result.complexity.style_score*100).toFixed(0)}%
              </div>
            )}

            <p><strong>Similar Repos:</strong></p>
            <ul>
              {result.similar_repos?.slice(0, 10).map((r: any) => (
                <li key={r.url}>
                  <a href={r.url} target="_blank" rel="noreferrer">{r.repo}</a>
                  {typeof r.stars === 'number' && r.stars > 0 && <span style={{marginLeft:8, color:'#888'}}>★ {r.stars}</span>}
                </li>
              ))}
            </ul>

            <h4>Preview</h4>
            <SyntaxHighlighter language={selectedLang === 'auto' ? result.language?.name?.toLowerCase() : selectedLang}>
              {code}
            </SyntaxHighlighter>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
