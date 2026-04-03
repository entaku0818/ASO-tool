'use client'

import { useState, useEffect } from 'react'
import { useApps } from '@/hooks/useApps'
import { useKeywords, KeywordWithRanking } from '@/hooks/useKeywords'
import { useSelectedApp } from '@/hooks/useSelectedApp'
import { useRankingHistory } from '@/hooks/useRankingHistory'
import { App, Ranking, createApp, deleteApp, CreateAppRequest } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'
import { UpgradeModal } from '@/components/UpgradeModal'
import { useUpgradeModal } from '@/hooks/useUpgradeModal'
import { OnboardingFlow, shouldShowOnboarding } from '@/components/OnboardingFlow'

function AppSidebar({
  apps,
  selectedApp,
  onSelectApp,
  onAddApp,
  onDeleteApp
}: {
  apps: App[]
  selectedApp: App | null
  onSelectApp: (app: App) => void
  onAddApp: () => void
  onDeleteApp: (app: App) => void
}) {
  return (
    <div className="w-56 bg-gray-900 text-white h-full overflow-y-auto">
      <div className="p-4 border-b border-gray-700 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-gray-400">Apps</h2>
        <button
          onClick={onAddApp}
          className="text-gray-400 hover:text-white transition-colors"
          title="アプリを追加"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
      </div>
      <div className="py-2">
        {apps.map((app) => (
          <div
            key={app.id}
            className={`group w-full px-4 py-3 flex items-center gap-3 hover:bg-gray-800 transition-colors cursor-pointer ${
              selectedApp?.id === app.id ? 'bg-gray-800 border-l-2 border-blue-500' : ''
            }`}
            onClick={() => onSelectApp(app)}
          >
            <div className="w-10 h-10 bg-gray-700 rounded-xl flex items-center justify-center text-lg">
              {app.platform === 'ios' ? '📱' : '🤖'}
            </div>
            <div className="text-left overflow-hidden flex-1">
              <div className="text-sm font-medium truncate">{app.name}</div>
              <div className="text-xs text-gray-400">
                {app.platform === 'ios' ? 'iPhone' : 'Android'}
              </div>
            </div>
            <button
              onClick={(e) => {
                e.stopPropagation()
                onDeleteApp(app)
              }}
              className="opacity-0 group-hover:opacity-100 text-gray-500 hover:text-red-400 transition-opacity"
              title="削除"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
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

function AddAppModal({
  isOpen,
  onClose,
  onSubmit,
  isSubmitting
}: {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: CreateAppRequest) => void
  isSubmitting: boolean
}) {
  const [name, setName] = useState('')
  const [bundleId, setBundleId] = useState('')
  const [platform, setPlatform] = useState<'ios' | 'android'>('ios')
  const [storeUrl, setStoreUrl] = useState('')

  if (!isOpen) return null

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim() || !bundleId.trim()) return
    onSubmit({
      name: name.trim(),
      bundle_id: bundleId.trim(),
      platform,
      store_url: storeUrl.trim() || undefined
    })
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4" onClick={e => e.stopPropagation()}>
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-bold">アプリを追加</h3>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">アプリ名</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="My App"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Bundle ID</label>
            <input
              type="text"
              value={bundleId}
              onChange={(e) => setBundleId(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="com.example.myapp"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">プラットフォーム</label>
            <select
              value={platform}
              onChange={(e) => setPlatform(e.target.value as 'ios' | 'android')}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="ios">iOS</option>
              <option value="android">Android</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">ストアURL (任意)</label>
            <input
              type="url"
              value={storeUrl}
              onChange={(e) => setStoreUrl(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="https://apps.apple.com/..."
            />
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-600 hover:text-gray-800"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !name.trim() || !bundleId.trim()}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {isSubmitting ? '追加中...' : '追加'}
            </button>
          </div>
        </form>
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
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function Home() {
  const { apps, isLoading: appsLoading, error: appsError, refetch } = useApps()
  const { selectedApp, setSelectedApp } = useSelectedApp(apps, appsLoading)
  const { keywords, isLoading: keywordsLoading } = useKeywords(selectedApp?.id || '')
  const { rankings, isLoading: rankingsLoading, fetchHistory, clearHistory } = useRankingHistory()
  const { user } = useAuth()
  const isPro = user?.is_pro ?? false
  const upgradeModal = useUpgradeModal()

  const [modalKeyword, setModalKeyword] = useState<KeywordWithRanking | null>(null)
  const [showAddModal, setShowAddModal] = useState(false)
  const [isAddingApp, setIsAddingApp] = useState(false)
  const [showOnboarding, setShowOnboarding] = useState(false)

  useEffect(() => {
    if (!appsLoading && user) {
      setShowOnboarding(shouldShowOnboarding(user.created_at, apps.length))
    }
  }, [appsLoading, user, apps.length])

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

  const handleAddApp = async (data: CreateAppRequest) => {
    setIsAddingApp(true)
    try {
      await createApp(data)
      setShowAddModal(false)
      refetch()
    } catch (e) {
      if (e instanceof Error && e.message.includes('402')) {
        setShowAddModal(false)
        upgradeModal.open('app_limit')
      } else {
        alert('アプリの追加に失敗しました')
      }
    } finally {
      setIsAddingApp(false)
    }
  }

  const handleDeleteApp = async (app: App) => {
    if (!confirm(`「${app.name}」を削除しますか？`)) return
    try {
      await deleteApp(app.id)
      if (selectedApp?.id === app.id) {
        setSelectedApp(apps.find(a => a.id !== app.id) || null)
      }
      refetch()
    } catch (e) {
      alert('アプリの削除に失敗しました')
    }
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
      <div className="flex flex-col items-center justify-center h-screen gap-4">
        <p className="text-gray-600">アプリが登録されていません</p>
        <button
          onClick={() => setShowAddModal(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          アプリを追加
        </button>
        <AddAppModal
          isOpen={showAddModal}
          onClose={() => setShowAddModal(false)}
          onSubmit={handleAddApp}
          isSubmitting={isAddingApp}
        />
      </div>
    )
  }

  return (
    <div className="flex h-[calc(100vh-73px)] -mx-4 -mt-6">
      <AppSidebar
        apps={apps}
        selectedApp={selectedApp}
        onSelectApp={setSelectedApp}
        onAddApp={() => {
          // Pattern D: Free users can register only 1 app
          if (!isPro && apps.length >= 1) {
            upgradeModal.open('app_limit')
            return
          }
          setShowAddModal(true)
        }}
        onDeleteApp={handleDeleteApp}
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

      <AddAppModal
        isOpen={showAddModal}
        onClose={() => setShowAddModal(false)}
        onSubmit={handleAddApp}
        isSubmitting={isAddingApp}
      />

      {/* Pattern D: triggered when Free user tries to register 2nd app */}
      <UpgradeModal
        isOpen={upgradeModal.isOpen}
        onClose={upgradeModal.close}
        subhead="複数のアプリをまとめて管理する / Proプランならアプリを無制限に登録できます"
      />

      {showOnboarding && user && (
        <OnboardingFlow
          userCreatedAt={user.created_at}
          appsCount={apps.length}
          onComplete={(newApp) => {
            setShowOnboarding(false)
            refetch()
          }}
          onDismiss={() => setShowOnboarding(false)}
        />
      )}
    </div>
  )
}
