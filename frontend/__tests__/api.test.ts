import { App, Keyword, TrackedKeyword } from '@/lib/api'

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
})
