'use client'

import { useEffect, useState } from 'react'
import { createCheckoutSession } from '@/lib/api'

const PRO_FEATURES = [
  'キーワード 無制限',
  '競合キーワード逆引き — 競合アプリの上位キーワードをリスト化',
  'CSVエクスポート',
  '多言語スクリーンショット 一括生成（6言語）',
  '背景プリセット 全種類',
]

interface UpgradeModalProps {
  isOpen: boolean
  onClose: () => void
  /** Feature name that triggered the modal — generates default subhead */
  triggerFeature?: string
  /** Overrides the auto-generated subhead (e.g. pattern D copy) */
  subhead?: string
}

export function UpgradeModal({ isOpen, onClose, triggerFeature, subhead }: UpgradeModalProps) {
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (!isOpen) return
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [isOpen, onClose])

  const handleUpgrade = async (planType: 'monthly' | 'yearly') => {
    setIsLoading(true)
    try {
      const { url } = await createCheckoutSession(planType)
      window.location.href = url
    } catch (err) {
      console.error('Failed to create checkout session:', err)
      setIsLoading(false)
    }
  }

  if (!isOpen) return null

  const resolvedSubhead = subhead ?? (triggerFeature ? `「${triggerFeature}」は Pro プランの機能です` : undefined)

  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      onClick={onClose}
    >
      <div
        className="bg-white rounded-2xl shadow-2xl p-8 max-w-sm w-full mx-4"
        onClick={e => e.stopPropagation()}
      >
        {/* Header */}
        <div className="text-center mb-6">
          <div className="text-4xl mb-3">🚀</div>
          <h2 className="text-xl font-bold text-gray-900">Pro にアップグレード</h2>
          {resolvedSubhead && (
            <p className="text-sm text-gray-500 mt-1">{resolvedSubhead}</p>
          )}
        </div>

        {/* Feature list */}
        <ul className="space-y-2 mb-6">
          {PRO_FEATURES.map(f => (
            <li key={f} className="flex items-center gap-2 text-sm text-gray-700">
              <span className="text-green-500 font-bold flex-shrink-0">✓</span>
              {f}
            </li>
          ))}
        </ul>

        {/* Pricing */}
        <div className="text-center mb-6 p-4 bg-blue-50 rounded-xl">
          <p className="text-2xl font-bold text-gray-900">
            ¥1,980<span className="text-sm font-normal text-gray-500">/月</span>
          </p>
          <p className="text-xs text-gray-500 mt-1">¥19,800/年（17% OFF・2ヶ月分お得）</p>
          <p className="text-xs text-gray-400 mt-1">いつでもキャンセル可能。</p>
        </div>

        {/* CTA */}
        <button
          onClick={() => handleUpgrade('monthly')}
          disabled={isLoading}
          className="w-full py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 transition-colors mb-3 disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {isLoading ? '処理中...' : '7日間無料で試す'}
        </button>
        <button
          onClick={onClose}
          className="w-full py-2 text-sm text-gray-400 hover:text-gray-600 transition-colors"
        >
          今は無料のまま使う
        </button>
      </div>
    </div>
  )
}
