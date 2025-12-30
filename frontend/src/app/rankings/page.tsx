'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { App, Keyword, getApps, getKeywords, getLatestRanking } from '@/lib/api'

type KeywordRanking = {
  app: App
  keyword: Keyword
  rank: number | null
}

export default function RankingsPage() {
  const router = useRouter()
  const [rankings, setRankings] = useState<KeywordRanking[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    async function fetchAllRankings() {
      try {
        const apps = await getApps()
        const allRankings: KeywordRanking[] = []

        for (const app of apps) {
          const keywords = await getKeywords(app.id)
          for (const keyword of keywords) {
            const ranking = await getLatestRanking(app.id, keyword.id)
            allRankings.push({
              app,
              keyword,
              rank: ranking?.rank ?? null,
            })
          }
        }

        // Sort by rank (nulls last)
        allRankings.sort((a, b) => {
          if (a.rank === null) return 1
          if (b.rank === null) return -1
          return a.rank - b.rank
        })

        setRankings(allRankings)
      } finally {
        setIsLoading(false)
      }
    }

    fetchAllRankings()
  }, [])

  const getRankColor = (rank: number | null) => {
    if (rank === null) return 'text-gray-400'
    if (rank <= 10) return 'text-green-600'
    if (rank <= 50) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getRankBgColor = (rank: number | null) => {
    if (rank === null) return 'bg-gray-50'
    if (rank <= 10) return 'bg-green-50'
    if (rank <= 50) return 'bg-yellow-50'
    return 'bg-red-50'
  }

  if (isLoading) {
    return (
      <div>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">キーワードランキング</h2>
          <Link href="/" className="text-blue-600 hover:underline">
            ダッシュボードへ
          </Link>
        </div>
        <p className="text-gray-600">読み込み中...</p>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">キーワードランキング ({rankings.length}件)</h2>
        <Link href="/" className="text-blue-600 hover:underline">
          ダッシュボードへ
        </Link>
      </div>

      {rankings.length === 0 ? (
        <p className="text-gray-600">キーワードが登録されていません</p>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <p className="text-sm text-gray-500 p-3 border-b">クリックでアプリ詳細ページへ</p>
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="py-3 px-4 text-left font-medium text-gray-600">順位</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">キーワード</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">アプリ</th>
                <th className="py-3 px-4 text-left font-medium text-gray-600">国</th>
              </tr>
            </thead>
            <tbody>
              {rankings.map((item) => (
                <tr
                  key={`${item.keyword.id}`}
                  className={`border-b cursor-pointer ${getRankBgColor(item.rank)} hover:bg-blue-50`}
                  onClick={() => router.push(`/apps/${item.app.id}`)}
                >
                  <td className={`py-3 px-4 font-bold text-lg ${getRankColor(item.rank)}`}>
                    {item.rank === null ? '圏外' : `${item.rank}位`}
                  </td>
                  <td className="py-3 px-4 font-medium">{item.keyword.keyword}</td>
                  <td className="py-3 px-4 text-blue-600">{item.app.name}</td>
                  <td className="py-3 px-4 text-gray-500">{item.keyword.country.toUpperCase()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
