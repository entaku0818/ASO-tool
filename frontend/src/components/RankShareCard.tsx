'use client'

import { useEffect, useState } from 'react'
import type { RisingKeyword } from '@/lib/api'

// ── 5パターンのシェアテキスト ────────────────────────────────
function buildShareTexts(kw: RisingKeyword): string[] {
  const { keyword, current_rank, previous_rank, improvement } = kw
  return [
    // パターン1（デフォルト）
    `ASO Tool でキーワード「${keyword}」が #${previous_rank} → #${current_rank} に上昇しました 🎉\n+${improvement}位アップ！ App Store 最適化の成果が出ています 📈\n\n#ASO #AppStoreOptimization #アプリ開発`,
    // パターン2
    `「${keyword}」が急上昇！ #${current_rank}位 を達成しました 🚀\n地道な ASO 施策が実を結びました✨\n\n#AppStoreOptimization #iOS #アプリ開発`,
    // パターン3
    `「${keyword}」の検索順位が ${improvement}位アップ 📈\n#${previous_rank}位 → #${current_rank}位！ ASO の力を実感しています\n\n#ASO #AppStore`,
    // パターン4
    `App Store で「${keyword}」が #${current_rank}位 に！\nこつこつ続けたキーワード最適化の成果です 🎯\n\n#アプリ開発 #ASO #AppStore`,
    // パターン5
    `今週の ASO 成果：「${keyword}」が ${improvement}ランクアップ ⬆️\n#${current_rank}位 まで上昇！\n\n#AppStore #AppStoreOptimization`,
  ]
}

// ── Share card preview (canvas-free, CSS only) ──────────────
function ShareCardPreview({ kw, text }: { kw: RisingKeyword; text: string }) {
  return (
    <div
      className="w-full rounded-xl overflow-hidden select-none"
      style={{ background: 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 50%, #ec4899 100%)' }}
    >
      <div className="px-6 py-5 text-white">
        <div className="text-3xl mb-2">🎉</div>
        <p className="font-bold text-lg leading-snug mb-3">
          キーワード順位が上昇しました！
        </p>
        <hr className="border-white/30 mb-3" />
        <p className="text-white/80 text-sm mb-1">「{kw.keyword}」</p>
        <div className="flex items-end gap-3">
          <span className="text-white/50 font-bold text-3xl">#{kw.previous_rank}</span>
          <span className="text-white/40 text-xl mb-1">──→──</span>
          <span className="text-white font-black text-4xl">#{kw.current_rank}</span>
          <span className="text-green-300 font-semibold text-base mb-1">+{kw.improvement}位 ↑</span>
        </div>
        <hr className="border-white/30 mt-3 mb-2" />
        <p className="text-white/50 text-xs">App Store Optimization · powered by ASO Tool</p>
      </div>
    </div>
  )
}

// ── Props ────────────────────────────────────────────────────
interface RankShareCardProps {
  keyword: RisingKeyword
  onClose: () => void
}

// ── Main component ───────────────────────────────────────────
export function RankShareCard({ keyword, onClose }: RankShareCardProps) {
  const patterns = buildShareTexts(keyword)
  const [selectedIdx, setSelectedIdx] = useState(0)
  const [copied, setCopied] = useState(false)

  const selectedText = patterns[selectedIdx]

  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [onClose])

  const handleTwitterShare = () => {
    const url = `https://twitter.com/intent/tweet?text=${encodeURIComponent(selectedText)}`
    window.open(url, '_blank', 'noopener,noreferrer')
  }

  const handleCopyText = async () => {
    try {
      await navigator.clipboard.writeText(selectedText)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // clipboard unavailable — silently ignore
    }
  }

  return (
    <div
      className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      <div
        className="bg-white dark:bg-[#12161e] rounded-2xl shadow-2xl dark:shadow-black/20 w-full max-w-md overflow-hidden"
        onClick={e => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b">
          <h2 className="font-semibold text-gray-900 dark:text-gray-100">シェアする</h2>
          <button onClick={onClose} className="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 transition-colors">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="p-6 space-y-5">
          {/* Card preview */}
          <ShareCardPreview kw={keyword} text={selectedText} />

          {/* Pattern selector */}
          <div>
            <p className="text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">シェアテキストを選ぶ</p>
            <div className="space-y-2 max-h-40 overflow-y-auto pr-1">
              {patterns.map((p, i) => (
                <label
                  key={i}
                  className={`flex items-start gap-2 p-2.5 rounded-lg border cursor-pointer transition-colors text-xs ${
                    selectedIdx === i
                      ? 'border-blue-400 bg-blue-50 dark:bg-blue-900/20'
                      : 'border-gray-100 dark:border-gray-800 hover:border-blue-200 dark:hover:border-blue-700'
                  }`}
                >
                  <input
                    type="radio"
                    name="share-pattern"
                    checked={selectedIdx === i}
                    onChange={() => setSelectedIdx(i)}
                    className="mt-0.5 accent-blue-600 flex-shrink-0"
                  />
                  <span className="text-gray-700 dark:text-gray-300 whitespace-pre-line leading-relaxed">
                    {p.split('\n')[0]}{p.includes('\n') ? ' …' : ''}
                  </span>
                  {i === 0 && (
                    <span className="ml-auto flex-shrink-0 text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-1.5 py-0.5 rounded-full font-medium">
                      推奨
                    </span>
                  )}
                </label>
              ))}
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex flex-col gap-2">
            <button
              onClick={handleTwitterShare}
              className="flex items-center justify-center gap-2 w-full py-2.5 rounded-xl font-semibold text-white text-sm transition-colors"
              style={{ background: 'linear-gradient(135deg, #1d9bf0, #1a8cd8)' }}
            >
              <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.741l7.73-8.835L1.254 2.25H8.08l4.253 5.622 5.91-5.622Zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
              </svg>
              X (Twitter) でシェア
            </button>
            <button
              onClick={handleCopyText}
              className="flex items-center justify-center gap-2 w-full py-2.5 rounded-xl font-semibold text-sm border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-[#0d0f14] transition-colors"
            >
              {copied ? '✅ コピーしました' : '📋 テキストをコピー'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
