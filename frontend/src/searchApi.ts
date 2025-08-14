import axios from 'axios'

export type SearchResult = {
  repo: string
  path: string
  language?: string
  updatedAt?: number
}

export async function searchCode(query: string): Promise<SearchResult[]> {
  if (!query || query.length < 3) return []
  const res = await axios.get('/search', { params: { q: query } })
  return res.data
}
