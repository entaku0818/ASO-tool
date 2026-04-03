'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Suspense } from 'react'

function BillingSuccessContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const sessionId = searchParams.get('session_id')
  const [countdown, setCountdown] = useState(5)

  useEffect(() => {
    if (countdown <= 0) {
      router.push('/')
      return
    }
    const timer = setTimeout(() => setCountdown(c => c - 1), 1000)
    return () => clearTimeout(timer)
  }, [countdown, router])

  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] gap-6 text-center">
      <div className="text-6xl">🎉</div>
      <h1 className="text-2xl font-bold text-gray-900">Proプランへようこそ！</h1>
      <p className="text-gray-600 max-w-sm">
        お支払いが完了しました。すべてのPro機能がご利用いただけます。
      </p>
      {sessionId && (
        <p className="text-xs text-gray-400">セッションID: {sessionId}</p>
      )}
      <p className="text-sm text-gray-500">
        {countdown}秒後にトップページへ移動します...
      </p>
      <button
        onClick={() => router.push('/')}
        className="px-6 py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 transition-colors"
      >
        今すぐ使い始める
      </button>
    </div>
  )
}

export default function BillingSuccessPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center min-h-[60vh]"><p className="text-gray-500">読み込み中...</p></div>}>
      <BillingSuccessContent />
    </Suspense>
  )
}
