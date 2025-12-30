const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://aso-api-671942133800.asia-northeast1.run.app'

export type App = {
  id: string
  name: string
  bundle_id: string
  platform: 'ios' | 'android'
  store_url?: string
  created_at: string
  updated_at: string
}

export type Keyword = {
  id: string
  app_id: string
  keyword: string
  country: string
  created_at: string
}

export type Ranking = {
  id: string
  keyword_id: string
  rank: number | null
  recorded_at: string
}

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`)
  }

  return response.json()
}

export async function getApps(): Promise<App[]> {
  return fetchApi<App[]>('/api/apps')
}

export async function getApp(id: string): Promise<App> {
  return fetchApi<App>(`/api/apps/${id}`)
}

export async function getKeywords(appId: string): Promise<Keyword[]> {
  return fetchApi<Keyword[]>(`/api/apps/${appId}/keywords`)
}

export async function getRankings(keywordId: string): Promise<Ranking[]> {
  return fetchApi<Ranking[]>(`/api/apps/_/keywords/${keywordId}/rankings`)
}

export async function getLatestRanking(appId: string, keywordId: string): Promise<Ranking | null> {
  try {
    return await fetchApi<Ranking>(`/api/apps/${appId}/keywords/${keywordId}/rankings/latest`)
  } catch {
    return null
  }
}

export async function createKeyword(appId: string, keyword: string, country: string = 'jp'): Promise<Keyword> {
  return fetchApi<Keyword>(`/api/apps/${appId}/keywords`, {
    method: 'POST',
    body: JSON.stringify({ keyword, country }),
  })
}

export async function deleteKeyword(appId: string, keywordId: string): Promise<void> {
  await fetchApi<void>(`/api/apps/${appId}/keywords/${keywordId}`, {
    method: 'DELETE',
  })
}

export type Review = {
  id: string
  app_id: string
  review_id: string
  author: string
  rating: number
  title: string
  content: string
  version: string
  posted_at: string
  created_at: string
}

export type ReviewStats = {
  total_count: number
  average_rating: number
  rating_counts: { [key: number]: number }
}

export async function getReviews(appId: string): Promise<Review[]> {
  return fetchApi<Review[]>(`/api/apps/${appId}/reviews`)
}

export async function getReviewStats(appId: string): Promise<ReviewStats> {
  return fetchApi<ReviewStats>(`/api/apps/${appId}/reviews/stats`)
}

// Tracked Keywords
export type TrackedKeyword = {
  id: string
  keyword: string
  country: string
  platform: string
  created_at: string
}

export type SearchResult = {
  id: string
  tracked_keyword_id: string
  rank: number
  app_name: string
  bundle_id: string
  developer: string
  recorded_at: string
}

export async function getTrackedKeywords(): Promise<TrackedKeyword[]> {
  return fetchApi<TrackedKeyword[]>('/api/tracked-keywords')
}

export async function getTrackedKeyword(id: string): Promise<TrackedKeyword> {
  return fetchApi<TrackedKeyword>(`/api/tracked-keywords/${id}`)
}

export async function getSearchResults(id: string): Promise<SearchResult[]> {
  return fetchApi<SearchResult[]>(`/api/tracked-keywords/${id}/results`)
}

export async function triggerTrackedKeywordUpdate(id: string): Promise<{ message: string; results_stored: number }> {
  return fetchApi(`/api/tracked-keywords/${id}/trigger`, { method: 'POST' })
}

export async function triggerAllTrackedKeywordUpdates(): Promise<{ message: string; keywords_processed: number; results_stored: number }> {
  return fetchApi('/api/tracked-keywords/trigger', { method: 'POST' })
}
