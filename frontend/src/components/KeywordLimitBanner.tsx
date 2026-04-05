'use client'

import { useState, useEffect } from 'react'

function lsKey(appId: string) {
  return `keyword_limit_banner_dismissed_${appId}`
}

function isDismissed(appId: string): boolean {
  try { return localStorage.getItem(lsKey(appId)) === 'true' } catch { return false }
}

interface KeywordLimitBannerProps {
  appId: string
  count: number
  limit?: number
  onUpgrade: () => void
}

export function KeywordLimitBanner({ appId, count, limit = 10, onUpgrade }: KeywordLimitBannerProps) {
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    // 7〜9件は dismiss 済みなら非表示。10件（上限）は dismiss 不可なので常に表示
    const atLimit = count >= limit
    if (atLimit) {
      setVisible(true)
    } else if (count >= 7 && !isDismissed(appId)) {
      setVisible(true)
    }
  }, [count, limit, appId])

  if (!visible) return null

  const remaining = limit - count
  const atLimit = remaining <= 0

  const handleDismiss = () => {
    try { localStorage.setItem(lsKey(appId), 'true') } catch {}
    setVisible(false)
  }

  return (
    <div className={`flex items-start justify-between gap-3 px-4 py-3 rounded-xl border mb-4 ${
      atLimit
        ? 'bg-red-50 border-red-200'
        : 'bg-amber-50 border-amber-200'
    }`}>
      <div className="flex items-start gap-2 min-w-0">
        <span className="flex-shrink-0 text-base mt-0.5">{atLimit ? '🔒' : '⚠️'}</span>
        <div className="min-w-0">
          {atLimit ? (
            <>
              <p className="text-sm font-semibold text-red-700">
                キーワード登録が上限（{limit}件）に達しています
              </p>
              <p className="text-xs text-red-600 mt-0.5">
                Pro プランにアップグレードすると無制限に登録できます。
              </p>
            </>
          ) : (
            <>
              <p className="text-sm font-semibold text-amber-700">
                キーワードを残り {remaining} 件登録できます（{count} / {limit}件）
              </p>
              <p className="text-xs text-amber-600 mt-0.5">
                Freeプランの上限は{limit}件です。Pro プランで無制限になります。
              </p>
            </>
          )}
        </div>
      </div>
      <div className="flex items-center gap-2 flex-shrink-0">
        {atLimit ? (
          // 上限到達: Pro CTA のみ、dismiss 不可
          <button
            onClick={onUpgrade}
            className="text-xs font-semibold text-white bg-red-500 hover:bg-red-600 px-3 py-1.5 rounded-lg transition-colors whitespace-nowrap"
          >
            Pro へアップグレード
          </button>
        ) : (
          // 警告（7〜9件）: dismiss 可
          <button
            onClick={handleDismiss}
            className="text-gray-400 hover:text-gray-600 transition-colors flex-shrink-0"
            aria-label="閉じる"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>
    </div>
  )
}
