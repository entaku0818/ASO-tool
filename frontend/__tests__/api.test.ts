import { App, Keyword, TrackedKeyword, SearchAdsCredentials, KeywordPopularitySuggestion } from '@/lib/api'

describe('API types', () => {
  it('should have correct App type structure', () => {
    const app: App = {
      id: '123',
      name: 'Test App',
      bundle_id: 'com.test.app',
      platform: 'ios',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    expect(app.id).toBe('123')
    expect(app.platform).toBe('ios')
  })

  it('should have correct Keyword type structure', () => {
    const keyword: Keyword = {
      id: '456',
      app_id: '123',
      keyword: 'test keyword',
      country: 'jp',
      created_at: '2024-01-01T00:00:00Z',
    }
    expect(keyword.keyword).toBe('test keyword')
  })

  it('should accept optional popularity fields on Keyword', () => {
    const keyword: Keyword = {
      id: '456',
      app_id: '123',
      keyword: 'calculator',
      country: 'jp',
      created_at: '2024-01-01T00:00:00Z',
      popularity_score: 4,
      popularity_fetched_at: '2024-01-01T00:00:00Z',
    }
    expect(keyword.popularity_score).toBe(4)
    expect(keyword.popularity_fetched_at).toBe('2024-01-01T00:00:00Z')
  })

  it('should have correct TrackedKeyword type structure', () => {
    const trackedKeyword: TrackedKeyword = {
      id: '789',
      keyword: 'tracked keyword',
      country: 'jp',
      platform: 'ios',
      created_at: '2024-01-01T00:00:00Z',
    }
    expect(trackedKeyword.keyword).toBe('tracked keyword')
    expect(trackedKeyword.platform).toBe('ios')
  })

  it('should have correct SearchAdsCredentials type structure', () => {
    const creds: SearchAdsCredentials = {
      id: 'cred-id',
      app_id: 'app-id',
      client_id: 'CLIENT_ID',
      team_id: 'TEAM_ID',
      key_id: 'KEY_ID',
      adam_id: 123456789,
      is_valid: true,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    expect(creds.client_id).toBe('CLIENT_ID')
    expect(creds.adam_id).toBe(123456789)
    expect(creds.is_valid).toBe(true)
  })

  it('should accept optional org_id on SearchAdsCredentials', () => {
    const creds: SearchAdsCredentials = {
      id: 'cred-id',
      app_id: 'app-id',
      client_id: 'CLIENT_ID',
      team_id: 'TEAM_ID',
      key_id: 'KEY_ID',
      adam_id: 123456789,
      is_valid: true,
      org_id: 'ORG_ID',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    expect(creds.org_id).toBe('ORG_ID')
  })

  it('should have correct KeywordPopularitySuggestion type structure', () => {
    const suggestion: KeywordPopularitySuggestion = {
      text: 'calculator',
      popularityScore: 5,
    }
    expect(suggestion.text).toBe('calculator')
    expect(suggestion.popularityScore).toBe(5)
  })
})

describe('Popularity score helpers', () => {
  function getDifficultyLabel(score?: number): string {
    if (score === undefined || score === null) return '−'
    if (score <= 1) return '低'
    if (score <= 3) return '中'
    return '高'
  }

  it('returns "−" for undefined score', () => {
    expect(getDifficultyLabel(undefined)).toBe('−')
  })

  it('returns "低" for score 0', () => {
    expect(getDifficultyLabel(0)).toBe('低')
  })

  it('returns "低" for score 1', () => {
    expect(getDifficultyLabel(1)).toBe('低')
  })

  it('returns "中" for score 2', () => {
    expect(getDifficultyLabel(2)).toBe('中')
  })

  it('returns "中" for score 3', () => {
    expect(getDifficultyLabel(3)).toBe('中')
  })

  it('returns "高" for score 4', () => {
    expect(getDifficultyLabel(4)).toBe('高')
  })

  it('returns "高" for score 5', () => {
    expect(getDifficultyLabel(5)).toBe('高')
  })
})
