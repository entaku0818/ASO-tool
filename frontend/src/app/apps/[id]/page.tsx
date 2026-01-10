'use client'

import { useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { useApp } from '@/hooks/useApp'
import { useKeywords, KeywordWithRanking } from '@/hooks/useKeywords'
import { useRankings } from '@/hooks/useRankings'
import { RankingChart } from '@/components/RankingChart'
import { createKeyword, deleteKeyword } from '@/lib/api'
import { ReviewsSection } from '@/components/ReviewsSection'
import { CompetitorSection } from '@/components/CompetitorSection'
import { AnalyticsSection } from '@/components/AnalyticsSection'

function KeywordRow({
  keyword,
  isSelected,
  onSelect,
  onDelete,
}: {
  keyword: KeywordWithRanking
  isSelected: boolean
  onSelect: () => void
  onDelete: () => void
}) {
  const rankColor = keyword.latestRank === null
    ? 'text-gray-500'
    : keyword.latestRank <= 10
    ? 'text-green-600'
    : keyword.latestRank <= 50
    ? 'text-yellow-600'
    : 'text-red-600'

  return (
    <tr
      className={`border-b cursor-pointer ${isSelected ? 'bg-blue-50' : 'hover:bg-gray-50'}`}
      onClick={onSelect}
    >
      <td className="py-3 px-4">{keyword.keyword}</td>
      <td className="py-3 px-4">{keyword.country}</td>
      <td className={`py-3 px-4 font-bold ${rankColor}`}>
        {keyword.latestRank === null ? '圏外' : `${keyword.latestRank}位`}
      </td>
      <td className="py-3 px-4">
        <button
          onClick={(e) => {
            e.stopPropagation()
            onDelete()
          }}
          className="text-red-500 hover:text-red-700 text-sm"
        >
          削除
        </button>
      </td>
    </tr>
  )
}

function RankingChartSection({ appId, keywordId, keywordName }: { appId: string; keywordId: string; keywordName: string }) {
  const { rankings, isLoading } = useRankings(appId, keywordId)

  return (
    <div className="bg-white rounded-lg shadow p-4 mb-6">
      <h3 className="text-lg font-semibold mb-4">「{keywordName}」の順位推移</h3>
      {isLoading ? (
        <p className="text-gray-500">読み込み中...</p>
      ) : (
        <RankingChart rankings={rankings} />
      )}
    </div>
  )
}

export default function AppDetailPage() {
  const params = useParams()
  const appId = params.id as string
  const { app, isLoading: appLoading, error: appError } = useApp(appId)
  const { keywords, isLoading: keywordsLoading, refetch } = useKeywords(appId)
  const [selectedKeyword, setSelectedKeyword] = useState<KeywordWithRanking | null>(null)
  const [newKeyword, setNewKeyword] = useState('')
  const [isAdding, setIsAdding] = useState(false)

  const handleAddKeyword = async () => {
    if (!newKeyword.trim()) return
    setIsAdding(true)
    try {
      await createKeyword(appId, newKeyword.trim())
      setNewKeyword('')
      refetch()
    } catch (e) {
      alert('キーワードの追加に失敗しました')
    } finally {
      setIsAdding(false)
    }
  }

  const handleDeleteKeyword = async (keywordId: string) => {
    if (!confirm('このキーワードを削除しますか？')) return
    try {
      await deleteKeyword(appId, keywordId)
      if (selectedKeyword?.id === keywordId) {
        setSelectedKeyword(null)
      }
      refetch()
    } catch (e) {
      alert('キーワードの削除に失敗しました')
    }
  }

  if (appLoading) {
    return (
      <div>
        <Link href="/" className="text-blue-600 hover:underline mb-4 inline-block">
          &larr; ダッシュボードに戻る
        </Link>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  if (appError || !app) {
    return (
      <div>
        <Link href="/" className="text-blue-600 hover:underline mb-4 inline-block">
          &larr; ダッシュボードに戻る
        </Link>
        <p className="text-red-600">アプリが見つかりません</p>
      </div>
    )
  }

  const sortedKeywords = [...keywords].sort((a, b) => {
    if (a.latestRank === null) return 1
    if (b.latestRank === null) return -1
    return a.latestRank - b.latestRank
  })

  return (
    <div>
      <Link href="/" className="text-blue-600 hover:underline mb-4 inline-block">
        &larr; ダッシュボードに戻る
      </Link>

      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold">{app.name}</h2>
            <p className="text-gray-500">{app.bundle_id}</p>
          </div>
          <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
            {app.platform.toUpperCase()}
          </span>
        </div>
        {app.store_url && (
          <a
            href={app.store_url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline text-sm mt-2 inline-block"
          >
            ストアページを開く
          </a>
        )}
      </div>

      {app.platform === 'ios' && (
        <div className="mb-6">
          <AnalyticsSection appId={appId} platform={app.platform} />
        </div>
      )}

      {selectedKeyword && (
        <RankingChartSection
          appId={appId}
          keywordId={selectedKeyword.id}
          keywordName={selectedKeyword.keyword}
        />
      )}

      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold">キーワード順位 ({keywords.length}件)</h3>
            <p className="text-sm text-gray-500">クリックで順位推移を表示</p>
          </div>
        </div>

        <div className="p-4 border-b bg-gray-50">
          <div className="flex gap-2">
            <input
              type="text"
              value={newKeyword}
              onChange={(e) => setNewKeyword(e.target.value)}
              placeholder="新しいキーワード"
              className="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              onKeyDown={(e) => e.key === 'Enter' && handleAddKeyword()}
            />
            <button
              onClick={handleAddKeyword}
              disabled={isAdding || !newKeyword.trim()}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {isAdding ? '追加中...' : '追加'}
            </button>
          </div>
        </div>

        {keywordsLoading ? (
          <p className="p-4 text-gray-500">読み込み中...</p>
        ) : keywords.length === 0 ? (
          <p className="p-4 text-gray-500">キーワードが登録されていません</p>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600">キーワード</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">国</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">順位</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600"></th>
              </tr>
            </thead>
            <tbody>
              {sortedKeywords.map((keyword) => (
                <KeywordRow
                  key={keyword.id}
                  keyword={keyword}
                  isSelected={selectedKeyword?.id === keyword.id}
                  onSelect={() => setSelectedKeyword(keyword)}
                  onDelete={() => handleDeleteKeyword(keyword.id)}
                />
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="mt-6">
        <CompetitorSection
          appId={appId}
          platform={app.platform}
          selectedKeywordId={selectedKeyword?.id}
        />
      </div>

      <div className="mt-6">
        <ReviewsSection appId={appId} />
      </div>
    </div>
  )
}
