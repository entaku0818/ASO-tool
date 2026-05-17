'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/AuthContext'

export default function HomePage() {
  const { user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && user) {
      router.replace('/dashboard')
    }
  }, [user, isLoading, router])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-gray-400 dark:text-gray-500">Loading...</p>
      </div>
    )
  }

  if (user) return null

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-[#0d0f14] dark:to-[#12161e]">
      <div className="max-w-4xl mx-auto px-4 py-20 text-center">
        <h1 className="text-5xl font-bold text-gray-900 dark:text-gray-100 mb-6">
          ASO Tool
        </h1>
        <p className="text-xl text-gray-600 dark:text-gray-400 mb-4">
          App Store キーワード最適化をデスクトップから
        </p>
        <p className="text-gray-500 dark:text-gray-400 mb-12 max-w-xl mx-auto">
          キーワード順位の追跡・競合分析・メタデータ管理を一元化。
          macOS ネイティブアプリで快適な ASO 作業を。
        </p>

        <div className="flex items-center justify-center gap-4">
          <Link
            href="/buy"
            className="px-8 py-3 bg-blue-600 dark:bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-700 transition-colors"
          >
            購入する ¥9,800/年
          </Link>
          <Link
            href="/login"
            className="px-8 py-3 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-semibold hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            ログイン
          </Link>
        </div>

        <div className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-6 text-left">
          {[
            { icon: '📊', title: 'キーワード順位追跡', desc: '複数キーワードの App Store 順位を毎日自動取得' },
            { icon: '🎯', title: '競合ギャップ分析', desc: '競合が上位・自社が圏外のキーワードを自動検出' },
            { icon: '📝', title: 'メタデータ管理', desc: 'タイトル・説明文をロケール/バージョン別に管理' },
          ].map((f) => (
            <div key={f.title} className="bg-white dark:bg-[#12161e] rounded-xl p-6 shadow-sm dark:shadow-black/20 border border-gray-100 dark:border-gray-800">
              <div className="text-3xl mb-3">{f.icon}</div>
              <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">{f.title}</h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">{f.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
