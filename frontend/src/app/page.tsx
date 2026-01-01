'use client'

import { useApps } from '@/hooks/useApps'
import { useKeywords, KeywordWithRanking } from '@/hooks/useKeywords'
import { useSelectedApp } from '@/hooks/useSelectedApp'
import { App } from '@/lib/api'

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

function KeywordTable({ keywords, isLoading }: { keywords: KeywordWithRanking[], isLoading: boolean }) {
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
                <RankBadge rank={keyword.latestRank} />
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
          <KeywordTable keywords={keywords} isLoading={keywordsLoading} />
        </div>
      </div>
    </div>
  )
}
