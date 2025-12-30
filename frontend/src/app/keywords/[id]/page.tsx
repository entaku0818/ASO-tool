'use client'

import { useState, useEffect, use } from 'react'
import Link from 'next/link'
import { TrackedKeyword, SearchResult, getTrackedKeyword, getSearchResults, triggerTrackedKeywordUpdate } from '@/lib/api'

export default function KeywordDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params)
  const [keyword, setKeyword] = useState<TrackedKeyword | null>(null)
  const [results, setResults] = useState<SearchResult[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isUpdating, setIsUpdating] = useState(false)

  useEffect(() => {
    fetchData()
  }, [resolvedParams.id])

  async function fetchData() {
    setIsLoading(true)
    try {
      const [keywordData, resultsData] = await Promise.all([
        getTrackedKeyword(resolvedParams.id),
        getSearchResults(resolvedParams.id)
      ])
      setKeyword(keywordData)
      setResults(resultsData || [])
    } catch (error) {
      console.error('Failed to fetch data:', error)
    } finally {
      setIsLoading(false)
    }
  }

  async function handleUpdate() {
    setIsUpdating(true)
    try {
      await triggerTrackedKeywordUpdate(resolvedParams.id)
      await fetchData()
    } finally {
      setIsUpdating(false)
    }
  }

  const getRankColor = (rank: number) => {
    if (rank <= 10) return 'text-green-600'
    if (rank <= 50) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getRankBgColor = (rank: number) => {
    if (rank <= 10) return 'bg-green-50'
    if (rank <= 50) return 'bg-yellow-50'
    return 'bg-red-50'
  }

  if (isLoading) {
    return (
      <div>
        <Link href="/keywords" className="text-blue-600 hover:underline mb-4 inline-block">
          &larr; キーワード一覧へ
        </Link>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  if (!keyword) {
    return (
      <div>
        <Link href="/keywords" className="text-blue-600 hover:underline mb-4 inline-block">
          &larr; キーワード一覧へ
        </Link>
        <p className="text-red-600">キーワードが見つかりませんでした</p>
      </div>
    )
  }

  return (
    <div>
      <Link href="/keywords" className="text-blue-600 hover:underline mb-4 inline-block">
        &larr; キーワード一覧へ
      </Link>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold">&quot;{keyword.keyword}&quot; の検索結果</h2>
          <p className="text-gray-500 text-sm mt-1">
            {keyword.platform === 'ios' ? 'iOS' : 'Android'} / {keyword.country.toUpperCase()}
          </p>
        </div>
        <button
          onClick={handleUpdate}
          disabled={isUpdating}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {isUpdating ? '更新中...' : 'ランキング更新'}
        </button>
      </div>

      {results.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-6 text-center">
          <p className="text-gray-600">検索結果がありません</p>
          <p className="text-sm text-gray-500 mt-2">
            「ランキング更新」ボタンをクリックして取得してください
          </p>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="p-3 border-b bg-gray-50">
            <p className="text-sm text-gray-500">
              最終更新: {new Date(results[0]?.recorded_at).toLocaleString('ja-JP')}
            </p>
          </div>
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600 w-20">順位</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">アプリ名</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">デベロッパー</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">Bundle ID</th>
              </tr>
            </thead>
            <tbody>
              {results.sort((a, b) => a.rank - b.rank).map((result) => (
                <tr
                  key={result.id}
                  className={`border-b ${getRankBgColor(result.rank)}`}
                >
                  <td className={`py-3 px-4 font-bold text-lg ${getRankColor(result.rank)}`}>
                    {result.rank}位
                  </td>
                  <td className="py-3 px-4 font-medium">{result.app_name}</td>
                  <td className="py-3 px-4 text-gray-600">{result.developer || '-'}</td>
                  <td className="py-3 px-4 text-gray-500 text-sm font-mono">{result.bundle_id || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
