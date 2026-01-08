'use client'

import { useState, useEffect } from 'react'
import {
  Competitor,
  getCompetitors,
  createCompetitor,
  deleteCompetitor,
  updateCompetitorRankings,
  getCompetitorComparison,
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
    if (rank === null) return 'text-gray-500'
    if (rank <= 10) return 'text-green-600'
    if (rank <= 50) return 'text-yellow-600'
    return 'text-red-600'
  }

  const formatRank = (rank: number | null) => {
    return rank === null ? '圏外' : `${rank}位`
  }

  if (isLoading) {
    return <div className="text-gray-500">読み込み中...</div>
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-4 border-b flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">競合アプリ ({competitors.length}件)</h3>
          <p className="text-sm text-gray-500">競合アプリの順位を比較</p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleUpdateRankings}
            disabled={isUpdating || competitors.length === 0}
            className="px-3 py-1.5 text-sm bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
          >
            {isUpdating ? '更新中...' : '順位更新'}
          </button>
          <button
            onClick={() => setShowAddForm(!showAddForm)}
            className="px-3 py-1.5 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            {showAddForm ? 'キャンセル' : '追加'}
          </button>
        </div>
      </div>

      {showAddForm && (
        <div className="p-4 bg-gray-50 border-b">
          <div className="space-y-3">
            <input
              type="text"
              value={newCompetitor.name}
              onChange={(e) => setNewCompetitor({ ...newCompetitor, name: e.target.value })}
              placeholder="アプリ名"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              value={newCompetitor.bundleId}
              onChange={(e) => setNewCompetitor({ ...newCompetitor, bundleId: e.target.value })}
              placeholder="Bundle ID (例: com.example.app)"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              value={newCompetitor.notes}
              onChange={(e) => setNewCompetitor({ ...newCompetitor, notes: e.target.value })}
              placeholder="メモ (任意)"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={handleAdd}
              disabled={isAdding || !newCompetitor.name.trim() || !newCompetitor.bundleId.trim()}
              className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {isAdding ? '追加中...' : '競合アプリを追加'}
            </button>
          </div>
        </div>
      )}

      {competitors.length === 0 ? (
        <p className="p-4 text-gray-500">競合アプリが登録されていません</p>
      ) : (
        <>
          {comparison && selectedKeywordId && (
            <div className="p-4 bg-blue-50 border-b">
              <h4 className="font-medium mb-3">「{comparison.keyword}」の順位比較</h4>
              <div className="space-y-2">
                <div className="flex items-center justify-between p-2 bg-white rounded border">
                  <span className="font-medium">あなたのアプリ</span>
                  <span className={`font-bold ${getRankColor(comparison.app_rank)}`}>
                    {formatRank(comparison.app_rank)}
                  </span>
                </div>
                {comparison.competitors.map((c) => (
                  <div key={c.competitor_id} className="flex items-center justify-between p-2 bg-white rounded border">
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
            <thead className="bg-gray-50">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600">アプリ名</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">Bundle ID</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">メモ</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600"></th>
              </tr>
            </thead>
            <tbody>
              {competitors.map((competitor) => (
                <tr key={competitor.id} className="border-b hover:bg-gray-50">
                  <td className="py-3 px-4">{competitor.competitor_name}</td>
                  <td className="py-3 px-4 text-sm text-gray-500">{competitor.competitor_bundle_id}</td>
                  <td className="py-3 px-4 text-sm text-gray-500">{competitor.notes || '-'}</td>
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
