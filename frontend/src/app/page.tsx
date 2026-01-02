'use client'

import { useState } from 'react'
import { useApps } from '@/hooks/useApps'
import { useKeywords, KeywordWithRanking } from '@/hooks/useKeywords'
import { useSelectedApp } from '@/hooks/useSelectedApp'
import { useRankingHistory } from '@/hooks/useRankingHistory'
import { App, Ranking } from '@/lib/api'

function AppSidebar({
  apps,
  selectedApp,
  onSelectApp
}: {
  apps: App[]
  selectedApp: App | null
  onSelectApp: (app: App) => void
}) {
  return (
    <div className="w-56 bg-gray-900 text-white h-full overflow-y-auto">
      <div className="p-4 border-b border-gray-700">
        <h2 className="text-sm font-semibold text-gray-400">Apps</h2>
      </div>
      <div className="py-2">
        {apps.map((app) => (
          <button
            key={app.id}
            onClick={() => onSelectApp(app)}
            className={`w-full px-4 py-3 flex items-center gap-3 hover:bg-gray-800 transition-colors ${
              selectedApp?.id === app.id ? 'bg-gray-800 border-l-2 border-blue-500' : ''
            }`}
          >
            <div className="w-10 h-10 bg-gray-700 rounded-xl flex items-center justify-center text-lg">
              {app.platform === 'ios' ? '📱' : '🤖'}
            </div>
            <div className="text-left overflow-hidden">
              <div className="text-sm font-medium truncate">{app.name}</div>
              <div className="text-xs text-gray-400">
                {app.platform === 'ios' ? 'iPhone' : 'Android'}
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  )
}

function CountryFlag({ country }: { country: string }) {
  const flags: Record<string, string> = {
    jp: '🇯🇵',
    us: '🇺🇸',
    de: '🇩🇪',
    cn: '🇨🇳',
    kr: '🇰🇷',
    gb: '🇬🇧',
    fr: '🇫🇷',
  }
  return <span>{flags[country.toLowerCase()] || '🌐'}</span>
}

function RankBadge({ rank }: { rank: number | null }) {
  if (rank === null) {
    return <span className="text-gray-400"># -</span>
  }

  const color = rank <= 10
    ? 'text-green-500'
    : rank <= 50
    ? 'text-yellow-500'
    : rank <= 100
    ? 'text-orange-500'
    : 'text-red-500'

  return <span className={`font-bold ${color}`}># {rank}</span>
}

function DifficultyBar({ value }: { value: number }) {
  const color = value <= 30
    ? 'bg-green-500'
    : value <= 60
    ? 'bg-yellow-500'
    : 'bg-red-500'

  return (
    <div className="flex items-center gap-2">
      <span className="text-sm w-8">{value}</span>
      <div className="w-20 h-2 bg-gray-200 rounded-full overflow-hidden">
        <div className={`h-full ${color}`} style={{ width: `${value}%` }} />
      </div>
    </div>
  )
}

function ChartIcon({ onClick }: { onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className="ml-2 p-1 text-gray-400 hover:text-blue-500 transition-colors"
      title="順位履歴を表示"
    >
      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4v16" />
      </svg>
    </button>
  )
}

function RankingChart({ rankings }: { rankings: Ranking[] }) {
  if (rankings.length === 0) {
    return <p className="text-gray-500 text-center py-8">データがありません</p>
  }

  const validRankings = rankings.filter(r => r.rank !== null)
  if (validRankings.length === 0) {
    return <p className="text-gray-500 text-center py-8">順位データがありません</p>
  }

  const maxRank = Math.max(...validRankings.map(r => r.rank!))
  const minRank = Math.min(...validRankings.map(r => r.rank!))
  const chartHeight = 200
  const chartWidth = 400
  const padding = 40

  const getY = (rank: number) => {
    const range = maxRank - minRank || 1
    return padding + ((rank - minRank) / range) * (chartHeight - padding * 2)
  }

  const getX = (index: number) => {
    return padding + (index / (validRankings.length - 1 || 1)) * (chartWidth - padding * 2)
  }

  const points = validRankings.map((r, i) => `${getX(i)},${getY(r.rank!)}`).join(' ')

  return (
    <div className="overflow-x-auto">
      <svg width={chartWidth} height={chartHeight} className="mx-auto">
        {/* Y axis labels */}
        <text x={padding - 5} y={padding} textAnchor="end" className="text-xs fill-gray-500">{minRank}</text>
        <text x={padding - 5} y={chartHeight - padding} textAnchor="end" className="text-xs fill-gray-500">{maxRank}</text>

        {/* Grid lines */}
        <line x1={padding} y1={padding} x2={padding} y2={chartHeight - padding} stroke="#e5e7eb" />
        <line x1={padding} y1={chartHeight - padding} x2={chartWidth - padding} y2={chartHeight - padding} stroke="#e5e7eb" />

        {/* Line chart */}
        <polyline
          fill="none"
          stroke="#3b82f6"
          strokeWidth="2"
          points={points}
        />

        {/* Data points */}
        {validRankings.map((r, i) => (
          <circle
            key={r.id}
            cx={getX(i)}
            cy={getY(r.rank!)}
            r="4"
            fill="#3b82f6"
          />
        ))}
      </svg>

      {/* Date labels */}
      <div className="flex justify-between px-10 text-xs text-gray-500 mt-2">
        <span>{new Date(validRankings[0].recorded_at).toLocaleDateString('ja-JP')}</span>
        <span>{new Date(validRankings[validRankings.length - 1].recorded_at).toLocaleDateString('ja-JP')}</span>
      </div>
    </div>
  )
}

function RankingModal({
  isOpen,
  onClose,
  keyword,
  rankings,
  isLoading
}: {
  isOpen: boolean
  onClose: () => void
  keyword: string
  rankings: Ranking[]
  isLoading: boolean
}) {
  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white rounded-lg p-6 max-w-lg w-full mx-4" onClick={e => e.stopPropagation()}>
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-bold">順位履歴: {keyword}</h3>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {isLoading ? (
          <p className="text-gray-500 text-center py-8">読み込み中...</p>
        ) : (
          <RankingChart rankings={rankings} />
        )}

        {/* Recent rankings table */}
        {!isLoading && rankings.length > 0 && (
          <div className="mt-4 max-h-40 overflow-y-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-gray-500">
                  <th className="text-left py-2">日付</th>
                  <th className="text-right py-2">順位</th>
                </tr>
              </thead>
              <tbody>
                {[...rankings].reverse().slice(0, 10).map(r => (
                  <tr key={r.id} className="border-b border-gray-100">
                    <td className="py-2">{new Date(r.recorded_at).toLocaleDateString('ja-JP')}</td>
                    <td className="text-right py-2">{r.rank ?? '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

function KeywordTable({
  keywords,
  isLoading,
  appId,
  onShowHistory
}: {
  keywords: KeywordWithRanking[]
  isLoading: boolean
  appId: string
  onShowHistory: (keyword: KeywordWithRanking) => void
}) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-gray-500">読み込み中...</p>
      </div>
    )
  }

  if (keywords.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-gray-500">キーワードが登録されていません</p>
      </div>
    )
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-gray-200 text-left text-sm text-gray-500">
            <th className="pb-3 font-medium">Keyword</th>
            <th className="pb-3 font-medium">Store</th>
            <th className="pb-3 font-medium">Position</th>
            <th className="pb-3 font-medium">Difficulty</th>
          </tr>
        </thead>
        <tbody>
          {keywords.map((keyword) => (
            <tr key={keyword.id} className="border-b border-gray-100 hover:bg-gray-50">
              <td className="py-4">
                <span className="font-medium">{keyword.keyword}</span>
              </td>
              <td className="py-4">
                <div className="flex items-center gap-2">
                  <CountryFlag country={keyword.country} />
                  <span className="text-sm text-gray-500">{keyword.country.toUpperCase()}</span>
                </div>
              </td>
              <td className="py-4">
                <div className="flex items-center">
                  <RankBadge rank={keyword.latestRank} />
                  <ChartIcon onClick={() => onShowHistory(keyword)} />
                </div>
              </td>
              <td className="py-4">
                <DifficultyBar value={Math.floor(Math.random() * 100)} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function Home() {
  const { apps, isLoading: appsLoading, error: appsError } = useApps()
  const { selectedApp, setSelectedApp } = useSelectedApp(apps, appsLoading)
  const { keywords, isLoading: keywordsLoading } = useKeywords(selectedApp?.id || '')
  const { rankings, isLoading: rankingsLoading, fetchHistory, clearHistory } = useRankingHistory()

  const [modalKeyword, setModalKeyword] = useState<KeywordWithRanking | null>(null)

  const handleShowHistory = (keyword: KeywordWithRanking) => {
    setModalKeyword(keyword)
    if (selectedApp) {
      fetchHistory(selectedApp.id, keyword.id)
    }
  }

  const handleCloseModal = () => {
    setModalKeyword(null)
    clearHistory()
  }

  if (appsLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  if (appsError) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-red-600">エラー: {appsError.message}</p>
      </div>
    )
  }

  if (apps.length === 0) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-gray-600">アプリが登録されていません</p>
      </div>
    )
  }

  return (
    <div className="flex h-[calc(100vh-73px)] -mx-4 -mt-6">
      <AppSidebar
        apps={apps}
        selectedApp={selectedApp}
        onSelectApp={setSelectedApp}
      />
      <div className="flex-1 overflow-y-auto">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h1 className="text-xl font-bold">{selectedApp?.name}</h1>
              <p className="text-sm text-gray-500">{selectedApp?.bundle_id}</p>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-500">
                {keywords.length} keywords
              </span>
            </div>
          </div>
          <KeywordTable
            keywords={keywords}
            isLoading={keywordsLoading}
            appId={selectedApp?.id || ''}
            onShowHistory={handleShowHistory}
          />
        </div>
      </div>

      <RankingModal
        isOpen={modalKeyword !== null}
        onClose={handleCloseModal}
        keyword={modalKeyword?.keyword || ''}
        rankings={rankings}
        isLoading={rankingsLoading}
      />
    </div>
  )
}
