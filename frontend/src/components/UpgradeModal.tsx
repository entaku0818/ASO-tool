'use client'

import { useEffect, useState } from 'react'
import { createCheckoutSession } from '@/lib/api'

const COMPARISON_ROWS = [
  { label: 'アプリ登録',         free: '1件まで',   pro: '無制限' },
  { label: 'キーワード',         free: '10件まで',  pro: '無制限' },
  { label: '競合キーワード逆引き', free: false,       pro: true },
  { label: 'CSVエクスポート',     free: false,       pro: true },
  { label: '多言語スクショ一括生成', free: false,     pro: true },
]

interface UpgradeModalProps {
  isOpen: boolean
  onClose: () => void
  /** Feature name that triggered the modal — generates default subhead */
  triggerFeature?: string
  /** Overrides the auto-generated subhead */
  subhead?: string
}

export function UpgradeModal({ isOpen, onClose, triggerFeature, subhead }: UpgradeModalProps) {
  const [planType, setPlanType] = useState<'monthly' | 'yearly'>('yearly')
  const [isLoading, setIsLoading] = useState(false)
  const [checkoutError, setCheckoutError] = useState<string | null>(null)

  useEffect(() => {
    if (!isOpen) {
      setCheckoutError(null)
      return
    }
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [isOpen, onClose])

  const handleUpgrade = async () => {
    setCheckoutError(null)
    setIsLoading(true)
    try {
      const { url } = await createCheckoutSession(planType)
      window.location.href = url
    } catch (err) {
      const msg = err instanceof Error && err.message.includes('401')
        ? 'セッションが切れました。再ログインしてからお試しください'
        : err instanceof Error && err.message.includes('Failed to fetch')
        ? 'インターネット接続を確認してから再度お試しください'
        : '決済ページへの接続に失敗しました。しばらくしてから再度お試しください'
      setCheckoutError(msg)
      setIsLoading(false)
    }
  }

  if (!isOpen) return null

  const resolvedSubhead = subhead ?? (triggerFeature ? `「${triggerFeature}」を使って差をつけよう` : undefined)

  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      onClick={onClose}
    >
      <div
        className="bg-white dark:bg-[#12161e] rounded-2xl shadow-2xl dark:shadow-black/20 max-w-md w-full mx-4 overflow-hidden"
        onClick={e => e.stopPropagation()}
      >
        {/* Gradient bar */}
        <div className="h-2 bg-gradient-to-r from-blue-600 to-purple-600" />

        <div className="p-8">
          {/* Header */}
          <div className="text-center mb-6">
            <div className="text-2xl mb-2">✦</div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Pro にアップグレード</h2>
            {resolvedSubhead && (
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">{resolvedSubhead}</p>
            )}
          </div>

          {/* FREE / PRO comparison */}
          <div className="grid grid-cols-2 gap-3 mb-6">
            {/* Free column */}
            <div className="bg-gray-50 dark:bg-[#0d0f14] rounded-xl p-4">
              <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3 uppercase tracking-wide">FREE（現在）</p>
              <ul className="space-y-2">
                {COMPARISON_ROWS.map(row => (
                  <li key={row.label} className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-400">
                    {typeof row.free === 'string' ? (
                      <>
                        <span className="text-gray-400 dark:text-gray-500 flex-shrink-0">—</span>
                        <span>{row.label}: {row.free}</span>
                      </>
                    ) : (
                      <>
                        <span className="text-red-400 flex-shrink-0">✗</span>
                        <span>{row.label}</span>
                      </>
                    )}
                  </li>
                ))}
              </ul>
            </div>

            {/* Pro column */}
            <div className="bg-blue-50 dark:bg-blue-900/20 ring-1 ring-blue-200 dark:ring-blue-700 rounded-xl p-4">
              <p className="text-xs font-semibold text-blue-700 dark:text-blue-300 mb-3 uppercase tracking-wide">✦ PRO</p>
              <ul className="space-y-2">
                {COMPARISON_ROWS.map(row => (
                  <li key={row.label} className="flex items-center gap-2 text-xs text-gray-700 dark:text-gray-300">
                    {typeof row.pro === 'string' ? (
                      <>
                        <span className="text-green-500 flex-shrink-0">✓</span>
                        <span>{row.label}: {row.pro}</span>
                      </>
                    ) : (
                      <>
                        <span className="text-green-500 flex-shrink-0">✓</span>
                        <span>{row.label}</span>
                      </>
                    )}
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Plan selector */}
          <div className="border border-gray-200 dark:border-gray-700 rounded-xl p-4 mb-5 space-y-3">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="planType"
                value="monthly"
                checked={planType === 'monthly'}
                onChange={() => setPlanType('monthly')}
                className="accent-blue-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">
                月払い <span className="font-semibold">¥1,980</span>
                <span className="text-gray-400 dark:text-gray-500"> / 月</span>
              </span>
            </label>
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="planType"
                value="yearly"
                checked={planType === 'yearly'}
                onChange={() => setPlanType('yearly')}
                className="accent-blue-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300 flex items-center gap-2">
                年払い <span className="font-semibold">¥1,650</span>
                <span className="text-gray-400 dark:text-gray-500"> / 月</span>
                <span className="bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400 text-xs font-semibold px-2 py-0.5 rounded-full">人気 17%OFF</span>
              </span>
            </label>
          </div>

          {/* Checkout error banner */}
          {checkoutError && (
            <div className="mb-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-xl p-4 text-sm">
              <p className="text-red-700 dark:text-red-300"><span className="font-bold">✗</span> {checkoutError}</p>
            </div>
          )}

          {/* CTA */}
          <button
            onClick={handleUpgrade}
            disabled={isLoading}
            className="w-full py-3 bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-xl font-semibold hover:from-blue-700 hover:to-blue-800 transition-all mb-3 disabled:opacity-60 disabled:cursor-not-allowed"
          >
            {isLoading ? '処理中...' : checkoutError ? 'もう一度試す' : 'Proプランを始める'}
          </button>
          <button
            onClick={onClose}
            className="w-full py-2 text-xs text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 transition-colors"
          >
            今は無料のまま使う（いつでも変更可）
          </button>
        </div>
      </div>
    </div>
  )
}
