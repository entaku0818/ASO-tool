const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://aso-api-671942133800.asia-northeast1.run.app'
const DB_MANAGER_URL = process.env.NEXT_PUBLIC_DB_MANAGER_URL || 'https://db-manager-te5er5txcq-an.a.run.app'

// Wake up the database when accessed
export async function wakeDatabase(): Promise<{ status: string; message: string } | null> {
  try {
    const response = await fetch(`${DB_MANAGER_URL}?action=start`, {
      method: 'GET',
    })
    return response.json()
  } catch {
    // Silently fail - DB might already be running
    return null
  }
}

export type App = {
  id: string
  name: string
  bundle_id: string
  platform: 'ios' | 'android'
  store_url?: string
  user_id?: string
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

function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('token')
}

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = getAuthToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options?.headers as Record<string, string>),
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  })

  if (response.status === 401) {
    // Clear token and redirect to login
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    throw new Error('Unauthorized')
  }

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

export type CreateAppRequest = {
  name: string
  bundle_id: string
  platform: 'ios' | 'android'
  store_url?: string
}

export async function createApp(data: CreateAppRequest): Promise<App> {
  return fetchApi<App>('/api/apps', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function updateApp(id: string, data: Partial<CreateAppRequest>): Promise<App> {
  return fetchApi<App>(`/api/apps/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export async function deleteApp(id: string): Promise<void> {
  await fetchApi<void>(`/api/apps/${id}`, {
    method: 'DELETE',
  })
}

export async function getKeywords(appId: string): Promise<Keyword[]> {
  return fetchApi<Keyword[]>(`/api/apps/${appId}/keywords`)
}

export async function getRankings(appId: string, keywordId: string): Promise<Ranking[]> {
  return fetchApi<Ranking[]>(`/api/apps/${appId}/keywords/${keywordId}/rankings`)
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

// Competitors
export type Competitor = {
  id: string
  app_id: string
  competitor_bundle_id: string
  competitor_name: string
  platform: 'ios' | 'android'
  notes?: string
  created_at: string
}

export type CompetitorRankData = {
  competitor_id: string
  competitor_name: string
  bundle_id: string
  rank: number | null
}

export type CompetitorComparison = {
  keyword: string
  keyword_id: string
  app_rank: number | null
  competitors: CompetitorRankData[]
  recorded_at: string
}

export type CreateCompetitorRequest = {
  competitor_bundle_id: string
  competitor_name: string
  platform: 'ios' | 'android'
  notes?: string
}

export async function getCompetitors(appId: string): Promise<Competitor[]> {
  return fetchApi<Competitor[]>(`/api/apps/${appId}/competitors`)
}

export async function createCompetitor(appId: string, data: CreateCompetitorRequest): Promise<Competitor> {
  return fetchApi<Competitor>(`/api/apps/${appId}/competitors`, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function deleteCompetitor(appId: string, competitorId: string): Promise<void> {
  await fetchApi<void>(`/api/apps/${appId}/competitors/${competitorId}`, {
    method: 'DELETE',
  })
}

export async function updateCompetitorRankings(appId: string): Promise<{ updated: number }> {
  return fetchApi(`/api/apps/${appId}/competitors/update-rankings`, { method: 'POST' })
}

export async function getCompetitorComparison(appId: string, keywordId: string): Promise<CompetitorComparison> {
  return fetchApi<CompetitorComparison>(`/api/apps/${appId}/keywords/${keywordId}/comparison`)
}

// App Store Connect Analytics
export type ASCCredentials = {
  id: string
  app_id: string
  issuer_id: string
  key_id: string
  is_valid: boolean
  last_validated_at?: string
  created_at: string
  updated_at: string
}

export type CreateASCCredentialsRequest = {
  issuer_id: string
  key_id: string
  private_key: string // Base64 encoded
}

export type Analytics = {
  id: string
  app_id: string
  date: string
  impressions: number
  downloads: number
  page_views: number
  conversion_rate?: number
}

export type AnalyticsWithCorrelation = Analytics & {
  ranking_change?: number
  new_version?: string
  review_count_change?: number
  average_rating_change?: number
}

export type AnalyticsSummary = {
  total_impressions: number
  total_downloads: number
  average_conversion: number
  impressions_change_percent: number
  downloads_change_percent: number
}

export async function getASCCredentials(appId: string): Promise<ASCCredentials | null> {
  try {
    return await fetchApi<ASCCredentials>(`/api/apps/${appId}/asc-credentials`)
  } catch {
    return null
  }
}

export async function setASCCredentials(appId: string, data: CreateASCCredentialsRequest): Promise<{ message: string }> {
  return fetchApi(`/api/apps/${appId}/asc-credentials`, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export async function deleteASCCredentials(appId: string): Promise<void> {
  await fetchApi<void>(`/api/apps/${appId}/asc-credentials`, {
    method: 'DELETE',
  })
}

export async function validateASCCredentials(appId: string): Promise<{ valid: boolean; message: string }> {
  return fetchApi(`/api/apps/${appId}/asc-credentials/validate`, { method: 'POST' })
}

export async function getAnalytics(appId: string, days: number = 30): Promise<Analytics[]> {
  return fetchApi<Analytics[]>(`/api/apps/${appId}/analytics?days=${days}`)
}

export async function getAnalyticsWithCorrelation(appId: string, days: number = 30): Promise<AnalyticsWithCorrelation[]> {
  return fetchApi<AnalyticsWithCorrelation[]>(`/api/apps/${appId}/analytics/correlation?days=${days}`)
}

export async function getAnalyticsSummary(appId: string, days: number = 30): Promise<AnalyticsSummary> {
  return fetchApi<AnalyticsSummary>(`/api/apps/${appId}/analytics/summary?days=${days}`)
}

export async function triggerAnalyticsFetch(appId: string): Promise<{ message: string; days_stored: number }> {
  return fetchApi(`/api/apps/${appId}/analytics/fetch`, { method: 'POST' })
}

// Public API (no authentication required)
export type AppRankingEntry = {
  rank: number
  name: string
  developer: string
  icon_url: string
  category: string
  store_url: string
  app_id: string
  price: string
  release_date: string
}

export async function getAppRankings(
  country: string = 'jp',
  rankingType: string = 'topfreeapplications',
  genreID: string = '',
  limit: number = 100,
): Promise<AppRankingEntry[]> {
  const params = new URLSearchParams({ country, ranking_type: rankingType, limit: String(limit) })
  if (genreID) params.set('genre_id', genreID)
  const response = await fetch(`${API_BASE_URL}/api/public/app-rankings?${params}`)
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`)
  }
  return response.json()
}

export type PopularKeyword = {
  keyword: string
  country: string
  platform: string
  results_count: number
  tracking_count: number
}

export async function getPopularKeywords(country: string = 'jp', platform: string = 'ios', limit: number = 100): Promise<PopularKeyword[]> {
  const response = await fetch(`${API_BASE_URL}/api/public/popular-keywords?country=${country}&platform=${platform}&limit=${limit}`)
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`)
  }
  return response.json()
}
