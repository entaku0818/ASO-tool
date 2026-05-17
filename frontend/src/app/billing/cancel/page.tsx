'use client'

import { useRouter } from 'next/navigation'

export default function BillingCancelPage() {
  const router = useRouter()

  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] gap-6 text-center">
      <div className="text-6xl">😔</div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">アップグレードをキャンセルしました</h1>
      <p className="text-gray-600 dark:text-gray-400 max-w-sm">
        引き続き無料プランをご利用いただけます。いつでもProプランにアップグレードできます。
      </p>
      <div className="flex gap-3">
        <button
          onClick={() => router.push('/')}
          className="px-6 py-3 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-xl font-semibold hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          ホームに戻る
        </button>
        <button
          onClick={() => router.back()}
          className="px-6 py-3 bg-blue-600 dark:bg-blue-500 text-white rounded-xl font-semibold hover:bg-blue-700 transition-colors"
        >
          もう一度試す
        </button>
      </div>
    </div>
  )
}
