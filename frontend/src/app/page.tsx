'use client'

import { useApps } from '@/hooks/useApps'
import { useKeywords, KeywordWithRanking } from '@/hooks/useKeywords'

function KeywordList({ appId }: { appId: string }) {
  const { keywords, isLoading } = useKeywords(appId)

  if (isLoading) {
    return <p className="text-sm text-gray-500">読み込み中...</p>
  }

  const sortedKeywords = [...keywords].sort((a, b) => {
    if (a.latestRank === null) return 1
    if (b.latestRank === null) return -1
    return a.latestRank - b.latestRank
  })

  return (
    <div className="mt-4">
      <h4 className="text-sm font-medium text-gray-700 mb-2">キーワード順位</h4>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {sortedKeywords.slice(0, 8).map((keyword) => (
          <KeywordBadge key={keyword.id} keyword={keyword} />
        ))}
      </div>
      {keywords.length > 8 && (
        <p className="text-sm text-gray-500 mt-2">他 {keywords.length - 8} 件</p>
      )}
    </div>
  )
}

function KeywordBadge({ keyword }: { keyword: KeywordWithRanking }) {
  const rankColor = keyword.latestRank === null
    ? 'bg-gray-100 text-gray-500'
    : keyword.latestRank <= 10
    ? 'bg-green-100 text-green-800'
    : keyword.latestRank <= 50
    ? 'bg-yellow-100 text-yellow-800'
    : 'bg-red-100 text-red-800'

  return (
    <div className={`px-3 py-2 rounded-lg ${rankColor}`}>
      <div className="text-xs truncate">{keyword.keyword}</div>
      <div className="font-bold">
        {keyword.latestRank === null ? '圏外' : `${keyword.latestRank}位`}
      </div>
    </div>
  )
}

export default function Home() {
  const { apps, isLoading, error } = useApps()

  if (isLoading) {
    return (
      <div>
        <h2 className="text-2xl font-bold mb-6">ダッシュボード</h2>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div>
        <h2 className="text-2xl font-bold mb-6">ダッシュボード</h2>
        <p className="text-red-600">エラー: {error.message}</p>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">ダッシュボード</h2>

      {apps.length === 0 ? (
        <p className="text-gray-600">アプリが登録されていません</p>
      ) : (
        <div className="space-y-6">
          {apps.map((app) => (
            <div key={app.id} className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold">{app.name}</h3>
                  <p className="text-sm text-gray-500">{app.bundle_id}</p>
                </div>
                <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                  {app.platform.toUpperCase()}
                </span>
              </div>
              <KeywordList appId={app.id} />
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
