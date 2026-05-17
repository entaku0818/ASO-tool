'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { TrackedKeyword, SearchResult, getTrackedKeyword, getSearchResults, triggerTrackedKeywordUpdate } from '@/lib/api'

export default function KeywordDetailPage() {
  const params = useParams()
  const id = params.id as string
  const [keyword, setKeyword] = useState<TrackedKeyword | null>(null)
  const [results, setResults] = useState<SearchResult[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isUpdating, setIsUpdating] = useState(false)

  useEffect(() => {
    if (id) {
      fetchData()
    }
  }, [id])

  async function fetchData() {
    setIsLoading(true)
    try {
      const [keywordData, resultsData] = await Promise.all([
        getTrackedKeyword(id),
        getSearchResults(id)
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
      await triggerTrackedKeywordUpdate(id)
      await fetchData()
    } finally {
      setIsUpdating(false)
    }
  }

  const getRankColor = (rank: number) => {
    if (rank <= 10) return 'text-green-600 dark:text-green-400'
    if (rank <= 50) return 'text-yellow-600 dark:text-yellow-400'
    return 'text-red-600 dark:text-red-400'
  }

  const getRankBgColor = (rank: number) => {
    if (rank <= 10) return 'bg-green-50 dark:bg-green-900/20'
    if (rank <= 50) return 'bg-yellow-50 dark:bg-yellow-900/20'
    return 'bg-red-50 dark:bg-red-900/20'
  }

  if (isLoading) {
    return (
      <div>
        <Link href="/keywords" className="text-blue-600 dark:text-blue-400 hover:underline mb-4 inline-block">
          &larr; キーワード一覧へ
        </Link>
        <p className="text-gray-600 dark:text-gray-400">読み込み中...</p>
      </div>
    )
  }

  if (!keyword) {
    return (
      <div>
        <Link href="/keywords" className="text-blue-600 dark:text-blue-400 hover:underline mb-4 inline-block">
          &larr; キーワード一覧へ
        </Link>
        <p className="text-red-600 dark:text-red-400">キーワードが見つかりませんでした</p>
      </div>
    )
  }

  return (
    <div>
      <Link href="/keywords" className="text-blue-600 dark:text-blue-400 hover:underline mb-4 inline-block">
        &larr; キーワード一覧へ
      </Link>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold">&quot;{keyword.keyword}&quot; の検索結果</h2>
          <p className="text-gray-500 dark:text-gray-400 text-sm mt-1">
            {keyword.platform === 'ios' ? 'iOS' : 'Android'} / {keyword.country.toUpperCase()}
          </p>
        </div>
        <button
          onClick={handleUpdate}
          disabled={isUpdating}
          className="bg-blue-600 dark:bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {isUpdating ? '更新中...' : 'ランキング更新'}
        </button>
      </div>

      {results.length === 0 ? (
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6 text-center">
          <p className="text-gray-600 dark:text-gray-400">検索結果がありません</p>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
            「ランキング更新」ボタンをクリックして取得してください
          </p>
        </div>
      ) : (
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
          <div className="p-3 border-b dark:border-gray-700 bg-gray-50 dark:bg-[#0d0f14]">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              最終更新: {new Date(results[0]?.recorded_at).toLocaleString('ja-JP')}
            </p>
          </div>
          <table className="w-full">
            <thead className="bg-gray-100 dark:bg-gray-800">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400 w-20">順位</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">アプリ名</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">デベロッパー</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">Bundle ID</th>
              </tr>
            </thead>
            <tbody>
              {results.sort((a, b) => a.rank - b.rank).map((result) => (
                <tr
                  key={result.id}
                  className={`border-b dark:border-gray-700 ${getRankBgColor(result.rank)}`}
                >
                  <td className={`py-3 px-4 font-bold text-lg ${getRankColor(result.rank)}`}>
                    {result.rank}位
                  </td>
                  <td className="py-3 px-4 font-medium">{result.app_name}</td>
                  <td className="py-3 px-4 text-gray-600 dark:text-gray-400">{result.developer || '-'}</td>
                  <td className="py-3 px-4 text-gray-500 dark:text-gray-400 text-sm font-mono">{result.bundle_id || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
