import axios from 'axios'

export type AnalyzeRequest = {
  filename: string
  code: string
  language?: string // optional, for user hint
}

export type AnalyzeResponse = {
  language: { name: string; confidence: number }
  framework: { name: string; confidence: number }
  security?: { issue: string; severity: string; cve?: string }[]
  purpose_guess?: string
  complexity?: { cyclomatic: number; style_score: number }
  similar_repos: { repo: string; stars: number; url: string }[]
}

export async function analyzeCode(body: AnalyzeRequest): Promise<AnalyzeResponse> {
  const res = await axios.post<AnalyzeResponse>('/api/analyze', body, {
    headers: { 'Content-Type': 'application/json' },
    timeout: 20000
  })
  return res.data
}
