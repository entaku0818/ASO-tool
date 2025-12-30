'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { TrackedKeyword, getTrackedKeywords, triggerAllTrackedKeywordUpdates } from '@/lib/api'

export default function KeywordsPage() {
  const router = useRouter()
  const [keywords, setKeywords] = useState<TrackedKeyword[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isUpdating, setIsUpdating] = useState(false)

  useEffect(() => {
    fetchKeywords()
  }, [])

  async function fetchKeywords() {
    try {
      const data = await getTrackedKeywords()
      setKeywords(data || [])
    } finally {
      setIsLoading(false)
    }
  }

  async function handleTriggerUpdate() {
    setIsUpdating(true)
    try {
      await triggerAllTrackedKeywordUpdates()
      await fetchKeywords()
    } finally {
      setIsUpdating(false)
    }
  }

  if (isLoading) {
    return (
      <div>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">キーワード検索結果</h2>
        </div>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">キーワード検索結果 ({keywords.length}件)</h2>
        <button
          onClick={handleTriggerUpdate}
          disabled={isUpdating}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {isUpdating ? '更新中...' : '全て更新'}
        </button>
      </div>

      <p className="text-sm text-gray-500 mb-4">
        各キーワードで検索したときのApp Storeランキングを表示します。クリックで詳細へ。
      </p>

      {keywords.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-6 text-center">
          <p className="text-gray-600">キーワードが登録されていません</p>
          <p className="text-sm text-gray-500 mt-2">
            アプリにキーワードを追加すると自動的に登録されます
          </p>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600">キーワード</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">国</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">プラットフォーム</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">登録日</th>
              </tr>
            </thead>
            <tbody>
              {keywords.map((keyword) => (
                <tr
                  key={keyword.id}
                  className="border-b cursor-pointer hover:bg-blue-50"
                  onClick={() => router.push(`/keywords/${keyword.id}`)}
                >
                  <td className="py-3 px-4 font-medium text-blue-600">{keyword.keyword}</td>
                  <td className="py-3 px-4 text-gray-500">{keyword.country.toUpperCase()}</td>
                  <td className="py-3 px-4">
                    <span className={`px-2 py-1 rounded text-xs font-medium ${
                      keyword.platform === 'ios' ? 'bg-gray-100 text-gray-700' : 'bg-green-100 text-green-700'
                    }`}>
                      {keyword.platform === 'ios' ? 'iOS' : 'Android'}
                    </span>
                  </td>
                  <td className="py-3 px-4 text-gray-500 text-sm">
                    {new Date(keyword.created_at).toLocaleDateString('ja-JP')}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
