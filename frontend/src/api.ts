import axios from 'axios'

export type AnalyzeRequest = {
  filename: string
  code: string
}

export type AnalyzeResponse = {
  language: string
  framework: string
  similar_repos: string[]
}

export async function analyzeCode(body: AnalyzeRequest): Promise<AnalyzeResponse> {
  const res = await axios.post<AnalyzeResponse>('/api/analyze', body, {
    headers: { 'Content-Type': 'application/json' },
    timeout: 20000
  })
  return res.data
}
