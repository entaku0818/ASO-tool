import type { Metadata } from 'next'
import Link from 'next/link'

// ── Server Component — OGP tags are SSR'd by generateMetadata ────────────────
// URL例: /share/abc123?kw=ゲーム&from=45&to=12&app=MyApp

type Props = {
  params: { id: string }
  searchParams: { kw?: string; from?: string; to?: string; app?: string }
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const keyword  = searchParams.kw  ?? ''
  const fromRank = parseInt(searchParams.from ?? '0', 10)
  const toRank   = parseInt(searchParams.to   ?? '0', 10)
  const appName  = searchParams.app ?? 'アプリ'
  const improvement = fromRank - toRank

  const title = improvement > 0
    ? `${appName}「${keyword}」が #${toRank} に急上昇！`
    : `${appName} のキーワード順位 — ASO Tool`
  const description = improvement > 0
    ? `キーワード「${keyword}」で #${fromRank} → #${toRank}（+${improvement}位）。ASO Tool でキーワード順位を自動追跡。`
    : 'ASO Tool でキーワード順位を自動追跡・分析。無料プランで今すぐ始められます。'

  return {
    title,
    description,
    openGraph: {
      title,
      description,
      type: 'website',
      siteName: 'ASO Tool',
    },
    twitter: {
      card: 'summary',
      title,
      description,
    },
  }
}

export default function SharePage({ searchParams }: Props) {
  const keyword  = searchParams.kw  ?? ''
  const fromRank = parseInt(searchParams.from ?? '0', 10)
  const toRank   = parseInt(searchParams.to   ?? '0', 10)
  const appName  = searchParams.app ?? 'このアプリ'
  const improvement = fromRank - toRank

  const isValid = keyword && fromRank > 0 && toRank > 0 && improvement > 0

  if (!isValid) {
    return (
      <div className="max-w-lg mx-auto text-center py-24">
        <div className="text-5xl mb-4">🔍</div>
        <h2 className="text-xl font-bold text-gray-800 dark:text-gray-200 mb-2">URLが正しくありません</h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6 text-sm">
          このリンクは無効か、必要な情報が不足しています。
        </p>
        <Link href="/" className="text-blue-600 dark:text-blue-400 hover:underline text-sm">
          トップページへ →
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* ── ヒーロー: グラデーションカード ────────────────────── */}
      <div
        className="w-full rounded-2xl overflow-hidden shadow-2xl"
        style={{ background: 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 50%, #ec4899 100%)' }}
      >
        <div className="px-8 py-8 text-white">
          <div className="text-4xl mb-3">🎉</div>
          <p className="font-bold text-xl leading-snug mb-1">{appName}</p>
          <p className="text-white/70 text-sm mb-4">のキーワード順位が急上昇しました！</p>
          <hr className="border-white/30 mb-4" />
          <p className="text-white/80 text-base mb-2">キーワード「{keyword}」</p>
          <div className="flex items-end gap-4 mb-4">
            <div className="text-center">
              <p className="text-white/50 text-xs mb-1">以前</p>
              <span className="text-white/50 font-bold text-5xl">#{fromRank}</span>
            </div>
            <span className="text-white/40 text-2xl mb-3 flex-shrink-0">──→──</span>
            <div className="text-center">
              <p className="text-white/70 text-xs mb-1">現在</p>
              <span className="text-white font-black text-6xl">#{toRank}</span>
            </div>
            <span className="text-green-300 font-bold text-xl mb-3">+{improvement}位 ↑</span>
          </div>
          <hr className="border-white/30 mb-3" />
          <p className="text-white/50 text-xs">App Store Optimization · powered by ASO Tool</p>
        </div>
      </div>

      {/* ── 実績ハイライト ──────────────────────────────────── */}
      <div className="grid grid-cols-3 gap-4 mt-6">
        {[
          { label: '順位改善', value: `+${improvement}位`, sub: '上昇' },
          { label: '達成順位', value: `#${toRank}`, sub: keyword },
          { label: '開始順位', value: `#${fromRank}`, sub: '→ 達成' },
        ].map(item => (
          <div key={item.label} className="bg-white dark:bg-[#12161e] rounded-xl shadow-sm dark:shadow-black/20 border border-gray-100 dark:border-gray-800 p-4 text-center">
            <p className="text-xs text-gray-400 dark:text-gray-500 mb-1">{item.label}</p>
            <p className="text-2xl font-black text-gray-900 dark:text-gray-100">{item.value}</p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1 truncate">{item.sub}</p>
          </div>
        ))}
      </div>

      {/* ── CTA ────────────────────────────────────────────── */}
      <div className="mt-8 bg-white dark:bg-[#12161e] rounded-2xl shadow-sm dark:shadow-black/20 border border-gray-100 dark:border-gray-800 px-8 py-7 text-center">
        <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
          あなたのアプリも上位表示を狙えます
        </h2>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
          ASO Tool でキーワード順位を自動追跡・分析。<br />
          無料プランでも今すぐ始められます。
        </p>
        <Link
          href="/login"
          className="block w-full py-3.5 rounded-xl font-bold text-white text-base shadow-lg transition-opacity hover:opacity-90 text-center"
          style={{ background: 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 50%, #ec4899 100%)' }}
        >
          無料で始める →
        </Link>
        <p className="text-xs text-gray-400 dark:text-gray-500 mt-3">クレジットカード不要・登録30秒</p>
      </div>

      {/* ── 機能紹介 ────────────────────────────────────────── */}
      <div className="mt-6 grid grid-cols-1 sm:grid-cols-3 gap-4">
        {[
          { icon: '📈', title: 'キーワード順位追跡', desc: 'App Store の検索順位を毎日自動で記録' },
          { icon: '🔍', title: '競合分析', desc: '競合アプリのキーワード戦略を逆引き' },
          { icon: '✨', title: 'AI キャプション生成', desc: 'スクリーンショット用テキストを AI が自動生成' },
        ].map(f => (
          <div key={f.title} className="bg-white dark:bg-[#12161e] rounded-xl border border-gray-100 dark:border-gray-800 p-4">
            <div className="text-2xl mb-2">{f.icon}</div>
            <p className="font-semibold text-gray-900 dark:text-gray-100 text-sm mb-1">{f.title}</p>
            <p className="text-xs text-gray-500 dark:text-gray-400">{f.desc}</p>
          </div>
        ))}
      </div>

      {/* ── Xでシェア ───────────────────────────────────────── */}
      {(() => {
        const tweetText = `「${keyword}」で #${toRank} 位達成！#${fromRank}位から+${improvement}位改善 🎉 #ASO #AppStore`
        const shareUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(tweetText)}&url=${encodeURIComponent(`https://aso-tool.entaku.app/share/${searchParams.kw ?? ''}`)}`
        return (
          <div className="mt-6 text-center">
            <a
              href={shareUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-xl font-bold text-white text-sm shadow transition-opacity hover:opacity-90"
              style={{ background: '#000000' }}
            >
              <svg className="w-4 h-4 fill-current" viewBox="0 0 24 24" aria-hidden="true">
                <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.746l7.73-8.835L1.254 2.25H8.08l4.253 5.622L18.244 2.25zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
              </svg>
              もう一度シェアする
            </a>
          </div>
        )
      })()}

      {/* ── フッター ────────────────────────────────────────── */}
      <p className="text-center text-xs text-gray-400 dark:text-gray-500 mt-8 mb-4">
        © ASO Tool — App Store Optimization Platform
      </p>
    </div>
  )
}
