'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { Suspense } from 'react'
import Link from 'next/link'

const PRO_FEATURES = [
  {
    icon: '🔍',
    title: '競合キーワード逆引き',
    desc: 'アプリ詳細 → 競合セクションから分析',
    href: '/',
  },
  {
    icon: '📱',
    title: 'アプリ・キーワード無制限',
    desc: '複数アプリをまとめて管理',
    href: '/',
  },
  {
    icon: '📊',
    title: 'CSVエクスポート',
    desc: 'キーワード一覧からエクスポート可能',
    href: '/',
  },
]

function BillingSuccessContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const sessionId = searchParams.get('session_id')

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-[#0d0f14] flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        {/* Hero card */}
        <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow-lg dark:shadow-black/20 overflow-hidden mb-4">
          <div className="h-2 bg-gradient-to-r from-blue-600 to-purple-600" />
          <div className="p-8 text-center">
            <div className="text-5xl mb-4">🎉</div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
              Proプランへようこそ！
            </h1>
            <p className="text-gray-500 dark:text-gray-400 text-sm mb-6">
              お支払いが完了しました。すべてのPro機能がご利用いただけます。
            </p>

            <button
              onClick={() => router.push('/?upgraded=1')}
              className="w-full py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl font-semibold hover:from-blue-700 hover:to-purple-700 transition-all"
            >
              ダッシュボードへ →
            </button>
          </div>
        </div>

        {/* Pro features */}
        <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow dark:shadow-black/20 p-6">
          <p className="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wide mb-4">
            さっそく使ってみましょう
          </p>
          <div className="space-y-3">
            {PRO_FEATURES.map((f) => (
              <Link
                key={f.title}
                href={f.href}
                className="flex items-start gap-3 p-3 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors group"
              >
                <span className="text-2xl">{f.icon}</span>
                <div>
                  <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                    {f.title}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">{f.desc}</p>
                </div>
                <span className="ml-auto text-gray-300 dark:text-gray-600 group-hover:text-blue-400 transition-colors">›</span>
              </Link>
            ))}
          </div>
        </div>

        {sessionId && (
          <p className="text-center text-xs text-gray-300 dark:text-gray-700 mt-4">
            決済ID: {sessionId}
          </p>
        )}
      </div>
    </div>
  )
}

export default function BillingSuccessPage() {
  return (
    <Suspense fallback={
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-gray-500">読み込み中...</p>
      </div>
    }>
      <BillingSuccessContent />
    </Suspense>
  )
}
