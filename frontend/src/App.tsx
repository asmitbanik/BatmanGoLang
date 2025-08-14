import React, { useState } from 'react'
import CodeEditor from './components/CodeEditor'
import { analyzeCode } from './api'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'

function App() {
  const [code, setCode] = useState<string>('package main\n\nimport "fmt"\n\nfunc main() {\n  fmt.Println("Hello, world")\n}\n')
  const [filename, setFilename] = useState('main.go')
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<any | null>(null)
  const [error, setError] = useState<string | null>(null)

  async function handleAnalyze() {
    setError(null)
    setLoading(true)
    setResult(null)
    try {
      const res = await analyzeCode({ filename, code })
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
        <h2>Shazam for Code — Live Editor</h2>
        <div>Backend: <strong>POST /analyze</strong></div>
      </div>

      <div className="controls">
        <input value={filename} onChange={(e) => setFilename(e.target.value)} style={{padding:8}} />
        <button onClick={handleAnalyze} disabled={loading} style={{padding:'8px 12px'}}>
          {loading ? 'Analyzing...' : 'Analyze'}
        </button>
      </div>

      <CodeEditor code={code} setCode={setCode} language="go" />

      <div style={{height: 16}} />

      <div className="result">
        <h3>Result</h3>
        {error && <div style={{color:'red'}}>Error: {error}</div>}
        {!result && !error && <div>Click Analyze to see results</div>}
        {result && (
          <div>
            <p><strong>Language:</strong> {result.language}</p>
            <p><strong>Framework:</strong> {result.framework}</p>
            <p><strong>Similar Repos:</strong></p>
            <ul>
              {result.similar_repos?.slice(0, 10).map((r: string) => (
                <li key={r}><a href={r} target="_blank" rel="noreferrer">{r}</a></li>
              ))}
            </ul>

            <h4>Preview</h4>
            <SyntaxHighlighter language="go">
              {code}
            </SyntaxHighlighter>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
