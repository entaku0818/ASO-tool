'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/AuthContext'

const FEATURES = [
  {
    icon: '📊',
    title: 'キーワード順位の自動追跡',
    desc: '登録したキーワードの App Store 順位を毎日自動で取得。最大90日のトレンドグラフで推移を一目で把握。複数キーワードを同時比較できます。',
  },
  {
    icon: '🎯',
    title: '競合キーワードギャップ分析',
    desc: '競合アプリが20位以内・自社が圏外のキーワードを自動検出。見落としていた狙い目キーワードをリストアップします。',
  },
  {
    icon: '📝',
    title: 'App Storeメタデータ管理',
    desc: 'タイトル・サブタイトル・説明文・キーワードをロケール別・バージョン別に管理。変更履歴も残せます。',
  },
  {
    icon: '⭐',
    title: 'レビュー・評価の追跡',
    desc: '各国ストアのレビューと評価の推移を一元管理。新規レビューをまとめて確認できます。',
  },
  {
    icon: '🔍',
    title: '人気キーワード調査',
    desc: 'Apple Search Ads の人気スコアをもとに注目キーワードを発見。国別・ジャンル別フィルターで絞り込み可能。',
  },
  {
    icon: '⚡',
    title: 'macOS ネイティブアプリ',
    desc: 'Dock やメニューバーから即アクセス。ブラウザ不要でサクサク動く、プロ仕様の作業環境。',
  },
]

const PAIN_POINTS = [
  { problem: 'キーワード順位を毎回手動で確認している', solution: '毎日自動取得。朝開くだけで最新情報が揃っている' },
  { problem: '競合が上位表示されているキーワードを把握できていない', solution: 'ギャップ分析機能が差分を自動リストアップ' },
  { problem: 'メタデータの変更履歴が Notion や Excel に分散している', solution: '専用のバージョン管理でロケール別に一元管理' },
]

const FAQS = [
  {
    q: '購入後はどうすれば使えますか？',
    a: '購入完了後、入力したメールアドレスにライセンスキーが届きます。macOS アプリを起動してキーとメールアドレスを入力するとすぐにアクティベートできます。',
  },
  {
    q: 'アプリはどこからダウンロードできますか？',
    a: '購入後にダウンロードリンクをメールでお送りします。macOS 13 Ventura 以降に対応しています。',
  },
  {
    q: '解約・返金はできますか？',
    a: 'いつでも解約できます。解約後はその年の有効期限まで引き続き利用できます。初回購入後7日以内であれば全額返金対応しています。',
  },
  {
    q: '1ライセンスで複数台使えますか？',
    a: '1ライセンスは1アカウント（1人）に対して発行されます。複数台のMacにインストールして同一アカウントで利用可能です。',
  },
  {
    q: '無料で試せる機能はありますか？',
    a: '人気キーワードページ（/popular-keywords）は登録なしで無料でご利用いただけます。App Store のキーワードサジェストやランキングデータを確認できます。',
  },
]

function FaqItem({ q, a }: { q: string; a: string }) {
  const [open, setOpen] = useState(false)
  return (
    <div className="border-b border-gray-200 dark:border-gray-700">
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center justify-between py-5 text-left text-gray-900 dark:text-gray-100 font-medium hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
      >
        <span>{q}</span>
        <span className={`ml-4 flex-shrink-0 text-xl transition-transform ${open ? 'rotate-45' : ''}`}>+</span>
      </button>
      {open && (
        <p className="pb-5 text-gray-600 dark:text-gray-400 text-sm leading-relaxed">{a}</p>
      )}
    </div>
  )
}

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
    <div className="-mx-4 -my-6">

      {/* ── Hero ── */}
      <section className="bg-gradient-to-br from-slate-900 via-slate-800 to-blue-950 text-white py-28 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <span className="inline-block px-4 py-1.5 bg-blue-500/20 border border-blue-400/30 text-blue-300 rounded-full text-sm font-medium mb-6">
            iOS デベロッパー向け ASO ツール
          </span>
          <h1 className="text-5xl sm:text-6xl font-bold leading-tight mb-6">
            キーワード最適化を、<br className="hidden sm:block" />
            もっとスマートに。
          </h1>
          <p className="text-xl text-slate-300 mb-10 max-w-2xl mx-auto leading-relaxed">
            順位追跡・競合ギャップ分析・メタデータ管理をひとつに。
            macOS ネイティブアプリで毎朝の ASO 作業を10分で完結。
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/buy"
              className="w-full sm:w-auto px-8 py-4 bg-blue-500 hover:bg-blue-400 text-white rounded-xl font-bold text-lg transition-colors shadow-lg shadow-blue-500/30"
            >
              今すぐ購入 ¥9,800/年 →
            </Link>
            <Link
              href="/popular-keywords"
              className="w-full sm:w-auto px-8 py-4 border border-white/30 text-white hover:bg-white/10 rounded-xl font-semibold text-lg transition-colors"
            >
              無料ツールを試す
            </Link>
          </div>
          <p className="mt-4 text-sm text-slate-400">7日間返金保証 · 月額費用なし · クレジットカード決済</p>
        </div>
      </section>

      {/* ── Trust bar ── */}
      <section className="bg-white dark:bg-[#12161e] border-b border-gray-200 dark:border-gray-800 py-5 px-4">
        <div className="max-w-4xl mx-auto">
          <div className="flex flex-wrap justify-center gap-x-10 gap-y-3 text-sm text-gray-500 dark:text-gray-400">
            {[
              { icon: '🔄', text: '毎日自動更新' },
              { icon: '💳', text: '買い切り（年額）' },
              { icon: '🔒', text: 'Stripe 決済' },
              { icon: '↩️', text: '7日間返金保証' },
              { icon: '🍎', text: 'macOS 13+対応' },
            ].map(({ icon, text }) => (
              <div key={text} className="flex items-center gap-1.5">
                <span>{icon}</span>
                <span>{text}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Pain → Solution ── */}
      <section className="bg-gray-50 dark:bg-[#0d0f14] py-20 px-4">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-gray-100 mb-4">
            こんな悩み、ありませんか？
          </h2>
          <p className="text-center text-gray-500 dark:text-gray-400 mb-12">ASO-tool はこれらをまとめて解決します</p>
          <div className="space-y-4">
            {PAIN_POINTS.map(({ problem, solution }) => (
              <div key={problem} className="bg-white dark:bg-[#12161e] rounded-xl p-6 shadow-sm dark:shadow-black/20 border border-gray-100 dark:border-gray-800 grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-start gap-3">
                  <span className="mt-0.5 flex-shrink-0 w-5 h-5 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center text-red-600 dark:text-red-400 text-xs font-bold">✗</span>
                  <p className="text-gray-700 dark:text-gray-300">{problem}</p>
                </div>
                <div className="flex items-start gap-3">
                  <span className="mt-0.5 flex-shrink-0 w-5 h-5 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center text-green-600 dark:text-green-400 text-xs font-bold">✓</span>
                  <p className="text-gray-700 dark:text-gray-300 font-medium">{solution}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section className="bg-white dark:bg-[#12161e] py-20 px-4">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-gray-100 mb-4">
            すべての機能
          </h2>
          <p className="text-center text-gray-500 dark:text-gray-400 mb-12">ASO に必要なものが1つのアプリに</p>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {FEATURES.map((f) => (
              <div
                key={f.title}
                className="rounded-xl p-6 border border-gray-100 dark:border-gray-800 hover:border-blue-200 dark:hover:border-blue-800 hover:shadow-md transition-all"
              >
                <div className="text-3xl mb-3">{f.icon}</div>
                <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">{f.title}</h3>
                <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Pricing ── */}
      <section className="bg-gray-50 dark:bg-[#0d0f14] py-20 px-4">
        <div className="max-w-md mx-auto text-center">
          <h2 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-4">シンプルな料金</h2>
          <p className="text-gray-500 dark:text-gray-400 mb-10">プランは1つ。全機能を使えます</p>
          <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow-xl dark:shadow-black/40 border border-gray-200 dark:border-gray-700 p-8">
            <div className="text-5xl font-bold text-gray-900 dark:text-gray-100 mb-1">¥9,800</div>
            <div className="text-gray-500 dark:text-gray-400 mb-6">年額 · 毎年自動更新</div>
            <ul className="space-y-3 text-sm text-left text-gray-600 dark:text-gray-400 mb-8">
              {[
                '全機能アンロック',
                '1ライセンス = 1アカウント',
                '複数台Mac にインストール可',
                'macOS 13 Ventura 以降対応',
                '購入後すぐにメールでキーを送付',
                'いつでも解約可・残期間は利用可',
                '7日間全額返金保証',
              ].map((item) => (
                <li key={item} className="flex items-center gap-2">
                  <span className="text-green-500 dark:text-green-400 flex-shrink-0">✓</span>
                  {item}
                </li>
              ))}
            </ul>
            <Link
              href="/buy"
              className="block w-full py-4 bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-400 text-white rounded-xl font-bold text-lg text-center transition-colors"
            >
              今すぐ購入する →
            </Link>
            <p className="text-xs text-gray-400 dark:text-gray-500 mt-3">決済は Stripe で安全に処理されます</p>
          </div>
        </div>
      </section>

      {/* ── FAQ ── */}
      <section className="bg-white dark:bg-[#12161e] py-20 px-4">
        <div className="max-w-2xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-gray-100 mb-12">よくある質問</h2>
          <div>
            {FAQS.map((faq) => (
              <FaqItem key={faq.q} q={faq.q} a={faq.a} />
            ))}
          </div>
        </div>
      </section>

      {/* ── Bottom CTA ── */}
      <section className="bg-gradient-to-br from-blue-600 to-blue-800 py-24 px-4 text-center text-white">
        <div className="max-w-2xl mx-auto">
          <h2 className="text-4xl font-bold mb-4">ASO 作業を今日から変えよう</h2>
          <p className="text-blue-100 mb-10 text-lg">7日間返金保証付き。まずは試してみてください。</p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/buy"
              className="w-full sm:w-auto px-10 py-4 bg-white text-blue-700 hover:bg-blue-50 rounded-xl font-bold text-lg transition-colors shadow-lg"
            >
              今すぐ購入 ¥9,800/年 →
            </Link>
            <Link
              href="/popular-keywords"
              className="w-full sm:w-auto px-10 py-4 border border-white/40 text-white hover:bg-white/10 rounded-xl font-semibold text-lg transition-colors"
            >
              無料ツールを試す
            </Link>
          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="bg-slate-900 text-slate-400 py-10 px-4">
        <div className="max-w-5xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-sm">
          <span className="font-semibold text-white">ASO Tool</span>
          <div className="flex gap-6">
            <Link href="/popular-keywords" className="hover:text-white transition-colors">人気キーワード</Link>
            <Link href="/buy" className="hover:text-white transition-colors">購入</Link>
            <Link href="/login" className="hover:text-white transition-colors">ログイン</Link>
          </div>
          <span>© 2025 ASO Tool</span>
        </div>
      </footer>
    </div>
  )
}
