'use client'

import { useState, useEffect } from 'react'
import {
  Competitor,
  KeywordGap,
  getCompetitors,
  createCompetitor,
  deleteCompetitor,
  updateCompetitorRankings,
  getCompetitorComparison,
  getKeywordGap,
  CompetitorComparison,
} from '@/lib/api'

type CompetitorSectionProps = {
  appId: string
  platform: 'ios' | 'android'
  selectedKeywordId?: string
}

export function CompetitorSection({ appId, platform, selectedKeywordId }: CompetitorSectionProps) {
  const [competitors, setCompetitors] = useState<Competitor[]>([])
  const [comparison, setComparison] = useState<CompetitorComparison | null>(null)
  const [gaps, setGaps] = useState<KeywordGap[]>([])
  const [showGap, setShowGap] = useState(false)
  const [isLoadingGap, setIsLoadingGap] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [isUpdating, setIsUpdating] = useState(false)
  const [showAddForm, setShowAddForm] = useState(false)
  const [newCompetitor, setNewCompetitor] = useState({ name: '', bundleId: '', notes: '' })
  const [isAdding, setIsAdding] = useState(false)

  const fetchCompetitors = async () => {
    try {
      const data = await getCompetitors(appId)
      setCompetitors(data)
    } catch (e) {
      console.error('Failed to fetch competitors:', e)
    }
  }

  const fetchComparison = async () => {
    if (!selectedKeywordId) {
      setComparison(null)
      return
    }
    try {
      const data = await getCompetitorComparison(appId, selectedKeywordId)
      setComparison(data)
    } catch (e) {
      console.error('Failed to fetch comparison:', e)
      setComparison(null)
    }
  }

  useEffect(() => {
    setIsLoading(true)
    Promise.all([fetchCompetitors(), fetchComparison()]).finally(() => setIsLoading(false))
  }, [appId, selectedKeywordId])

  const handleAdd = async () => {
    if (!newCompetitor.name.trim() || !newCompetitor.bundleId.trim()) return
    setIsAdding(true)
    try {
      await createCompetitor(appId, {
        competitor_name: newCompetitor.name.trim(),
        competitor_bundle_id: newCompetitor.bundleId.trim(),
        platform,
        notes: newCompetitor.notes.trim() || undefined,
      })
      setNewCompetitor({ name: '', bundleId: '', notes: '' })
      setShowAddForm(false)
      fetchCompetitors()
    } catch (e) {
      alert('競合アプリの追加に失敗しました')
    } finally {
      setIsAdding(false)
    }
  }

  const handleDelete = async (competitorId: string) => {
    if (!confirm('この競合アプリを削除しますか？')) return
    try {
      await deleteCompetitor(appId, competitorId)
      fetchCompetitors()
      if (selectedKeywordId) {
        fetchComparison()
      }
    } catch (e) {
      alert('削除に失敗しました')
    }
  }

  const handleShowGap = async () => {
    if (showGap) {
      setShowGap(false)
      return
    }
    setIsLoadingGap(true)
    try {
      const data = await getKeywordGap(appId)
      setGaps(data)
      setShowGap(true)
    } catch (e) {
      console.error('Failed to fetch keyword gap:', e)
    } finally {
      setIsLoadingGap(false)
    }
  }

  const handleUpdateRankings = async () => {
    setIsUpdating(true)
    try {
      const result = await updateCompetitorRankings(appId)
      alert(`${result.updated}件の順位を更新しました`)
      if (selectedKeywordId) {
        fetchComparison()
      }
    } catch (e) {
      alert('順位の更新に失敗しました')
    } finally {
      setIsUpdating(false)
    }
  }

  const getRankColor = (rank: number | null) => {
    if (rank === null) return 'text-gray-500 dark:text-gray-400'
    if (rank <= 10) return 'text-green-600 dark:text-green-400'
    if (rank <= 50) return 'text-yellow-600 dark:text-yellow-400'
    return 'text-red-600 dark:text-red-400'
  }

  const formatRank = (rank: number | null) => {
    return rank === null ? '圏外' : `${rank}位`
  }

  if (isLoading) {
    return <div className="text-gray-500 dark:text-gray-400">読み込み中...</div>
  }

  return (
    <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20">
      <div className="p-4 border-b flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">競合アプリ ({competitors.length}件)</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">競合アプリの順位を比較</p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleShowGap}
            disabled={isLoadingGap || competitors.length === 0}
            className="px-3 py-1.5 text-sm bg-purple-600 text-white rounded hover:bg-purple-700 disabled:opacity-50 flex items-center gap-1"
          >
            {isLoadingGap ? '取得中...' : showGap ? 'ギャップ非表示' : 'キーワードギャップ'}
            {!showGap && gaps.length > 0 && (
              <span className="inline-flex items-center justify-center w-5 h-5 text-xs bg-white dark:bg-[#12161e] text-purple-700 rounded-full font-bold">
                {gaps.length}
              </span>
            )}
          </button>
          <button
            onClick={handleUpdateRankings}
            disabled={isUpdating || competitors.length === 0}
            className="px-3 py-1.5 text-sm bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
          >
            {isUpdating ? '更新中...' : '順位更新'}
          </button>
          <button
            onClick={() => setShowAddForm(!showAddForm)}
            className="px-3 py-1.5 text-sm bg-blue-600 dark:bg-blue-500 text-white rounded hover:bg-blue-700 dark:hover:bg-blue-600"
          >
            {showAddForm ? 'キャンセル' : '追加'}
          </button>
        </div>
      </div>

      {showAddForm && (
        <div className="p-4 bg-gray-50 dark:bg-[#0d0f14] border-b">
          <div className="space-y-3">
            <input
              type="text"
              value={newCompetitor.name}
              onChange={(e) => setNewCompetitor({ ...newCompetitor, name: e.target.value })}
              placeholder="アプリ名"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
            />
            <input
              type="text"
              value={newCompetitor.bundleId}
              onChange={(e) => setNewCompetitor({ ...newCompetitor, bundleId: e.target.value })}
              placeholder="Bundle ID (例: com.example.app)"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
            />
            <input
              type="text"
              value={newCompetitor.notes}
              onChange={(e) => setNewCompetitor({ ...newCompetitor, notes: e.target.value })}
              placeholder="メモ (任意)"
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
            />
            <button
              onClick={handleAdd}
              disabled={isAdding || !newCompetitor.name.trim() || !newCompetitor.bundleId.trim()}
              className="w-full px-4 py-2 bg-blue-600 dark:bg-blue-500 text-white rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600 disabled:opacity-50"
            >
              {isAdding ? '追加中...' : '競合アプリを追加'}
            </button>
          </div>
        </div>
      )}

      {showGap && (
        <div className="p-4 border-b bg-purple-50">
          <h4 className="font-medium mb-1 text-purple-800">キーワードギャップ</h4>
          <p className="text-xs text-purple-600 mb-3">競合が20位以内・自社が30位以下またはランク外のキーワード</p>
          {gaps.length === 0 ? (
            <p className="text-sm text-gray-500 dark:text-gray-400">ギャップキーワードはありません</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-left text-xs text-purple-700">
                    <th className="pb-2 pr-4">キーワード</th>
                    <th className="pb-2 pr-4">国</th>
                    <th className="pb-2 pr-4">競合アプリ</th>
                    <th className="pb-2 pr-4">競合順位</th>
                    <th className="pb-2">自社順位</th>
                  </tr>
                </thead>
                <tbody>
                  {gaps.map((g) => (
                    <tr key={g.keyword_id} className="border-t border-purple-100">
                      <td className="py-1.5 pr-4 font-medium">{g.keyword}</td>
                      <td className="py-1.5 pr-4 text-gray-500 dark:text-gray-400">{g.country}</td>
                      <td className="py-1.5 pr-4 text-gray-600 dark:text-gray-400">{g.competitor_name}</td>
                      <td className="py-1.5 pr-4 text-green-700 dark:text-green-300 font-bold">{g.competitor_rank}位</td>
                      <td className="py-1.5 text-red-600 dark:text-red-400 font-bold">
                        {g.our_rank === null ? '圏外' : `${g.our_rank}位`}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {competitors.length === 0 ? (
        <p className="p-4 text-gray-500 dark:text-gray-400">競合アプリが登録されていません</p>
      ) : (
        <>
          {comparison && selectedKeywordId && (
            <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border-b">
              <h4 className="font-medium mb-3">「{comparison.keyword}」の順位比較</h4>
              <div className="space-y-2">
                <div className="flex items-center justify-between p-2 bg-white dark:bg-[#12161e] rounded border dark:border-gray-700">
                  <span className="font-medium">あなたのアプリ</span>
                  <span className={`font-bold ${getRankColor(comparison.app_rank)}`}>
                    {formatRank(comparison.app_rank)}
                  </span>
                </div>
                {comparison.competitors.map((c) => (
                  <div key={c.competitor_id} className="flex items-center justify-between p-2 bg-white dark:bg-[#12161e] rounded border dark:border-gray-700">
                    <span>{c.competitor_name}</span>
                    <span className={`font-bold ${getRankColor(c.rank)}`}>
                      {formatRank(c.rank)}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}

          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-[#0d0f14]">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">アプリ名</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">Bundle ID</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">メモ</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400"></th>
              </tr>
            </thead>
            <tbody>
              {competitors.map((competitor) => (
                <tr key={competitor.id} className="border-b hover:bg-gray-50 dark:hover:bg-[#0d0f14]">
                  <td className="py-3 px-4">{competitor.competitor_name}</td>
                  <td className="py-3 px-4 text-sm text-gray-500 dark:text-gray-400">{competitor.competitor_bundle_id}</td>
                  <td className="py-3 px-4 text-sm text-gray-500 dark:text-gray-400">{competitor.notes || '-'}</td>
                  <td className="py-3 px-4">
                    <button
                      onClick={() => handleDelete(competitor.id)}
                      className="text-red-500 hover:text-red-700 text-sm"
                    >
                      削除
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  )
}
