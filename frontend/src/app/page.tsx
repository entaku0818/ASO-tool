'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { useAuth } from '@/contexts/AuthContext'
import { PLAN_LIMITS, PRICES, formatYen } from '@/lib/plans'

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
    q: 'Webダッシュボードの無料プランでできることは？',
    a: `無料プランでは${PLAN_LIMITS.freeAppLimit}アプリ・${PLAN_LIMITS.freeKeywordLimit}キーワードまで登録でき、毎日の順位追跡・App Storeアナリティクス連携・キーワード提案などが使えます。登録不要の人気キーワードページも無料です。`,
  },
  {
    q: 'WebダッシュボードのProプランとmacOSアプリの違いは？',
    a: `Webダッシュボード（Pro）はブラウザで使えるSaaSで、${formatYen(PRICES.webProYearlyMonthlyEquivalent)}/月〜のサブスクリプションです。macOSアプリは${formatYen(PRICES.macosLicenseYearly)}/年の買い切りで、Dock・メニューバーからアクセスできるネイティブアプリです。機能は共通していますが、使い方の好みで選べます。`,
  },
  {
    q: 'macOSアプリ購入後はどうすれば使えますか？',
    a: '購入完了後、入力したメールアドレスにライセンスキーが届きます。macOS アプリを起動してキーとメールアドレスを入力するとすぐにアクティベートできます。',
  },
  {
    q: '解約・返金はできますか？',
    a: 'Webダッシュボード Proはいつでも解約でき、解約後は期間終了まで利用可能です。macOSアプリは初回購入後7日以内であれば全額返金対応しています。',
  },
  {
    q: '1ライセンスで複数台使えますか？',
    a: 'macOSアプリは1ライセンス（1アカウント）で複数台のMacにインストールして利用可能です。Webダッシュボードはブラウザからどこでもアクセスできます。',
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
            無料で始めて、必要に応じてアップグレード。
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/signup"
              className="w-full sm:w-auto px-8 py-4 bg-blue-500 hover:bg-blue-400 text-white rounded-xl font-bold text-lg transition-colors shadow-lg shadow-blue-500/30"
            >
              無料で始める →
            </Link>
            <Link
              href="/buy"
              className="w-full sm:w-auto px-8 py-4 border border-white/30 text-white hover:bg-white/10 rounded-xl font-semibold text-lg transition-colors"
            >
              macOSアプリを購入 {formatYen(PRICES.macosLicenseYearly)}/年
            </Link>
          </div>
          <p className="mt-4 text-sm text-slate-400">クレジットカード不要で無料スタート · Stripe 決済対応</p>
        </div>
      </section>

      {/* ── Trust bar ── */}
      <section className="bg-white dark:bg-[#12161e] border-b border-gray-200 dark:border-gray-800 py-5 px-4">
        <div className="max-w-4xl mx-auto">
          <div className="flex flex-wrap justify-center gap-x-10 gap-y-3 text-sm text-gray-500 dark:text-gray-400">
            {[
              { icon: '🆓', text: '無料プランあり' },
              { icon: '🔄', text: '毎日自動更新' },
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

      {/* ── Free Tool Spotlight ── */}
      <section className="bg-white dark:bg-[#12161e] py-20 px-4">
        <div className="max-w-5xl mx-auto">
          <div className="flex flex-col lg:flex-row items-center gap-12">
            <div className="flex-1">
              <span className="inline-block px-3 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 text-xs font-bold rounded-full mb-4 uppercase tracking-wide">
                登録不要・完全無料
              </span>
              <h2 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-4">
                まずは無料ツールで<br />キーワード調査を体験
              </h2>
              <p className="text-gray-500 dark:text-gray-400 mb-6 leading-relaxed">
                App Store のキーワードサジェストをリアルタイムで取得。
                アカウント登録不要で今すぐ使えます。
              </p>
              <ul className="space-y-2 mb-8">
                {[
                  'App Store 検索サジェストをリアルタイム取得',
                  '日本・アメリカ・中国など7カ国対応',
                  'App Store ランキングも無料で確認',
                ].map((item) => (
                  <li key={item} className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <span className="text-green-500 flex-shrink-0">✓</span>
                    {item}
                  </li>
                ))}
              </ul>
              <Link
                href="/popular-keywords"
                className="inline-flex items-center gap-2 px-6 py-3 bg-green-600 hover:bg-green-700 text-white rounded-xl font-bold transition-colors"
              >
                無料ツールを使ってみる →
              </Link>
            </div>
            <div className="flex-1 w-full">
              <div className="rounded-2xl overflow-hidden border border-gray-200 dark:border-gray-700 shadow-xl dark:shadow-black/40">
                <div className="bg-gray-100 dark:bg-gray-800 px-4 py-2 flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-red-400" />
                  <div className="w-3 h-3 rounded-full bg-yellow-400" />
                  <div className="w-3 h-3 rounded-full bg-green-400" />
                  <span className="ml-2 text-xs text-gray-400">aso-tool.entaku.app/popular-keywords</span>
                </div>
                <Image
                  src="/screenshot-keyword-search.png"
                  alt="キーワード検索機能のスクリーンショット"
                  width={760}
                  height={480}
                  className="w-full"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ── Pain → Solution ── */}
      <section className="bg-gray-50 dark:bg-[#0d0f14] py-20 px-4">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-gray-100 mb-4">
            こんな悩み、ありませんか？
          </h2>
          <p className="text-center text-gray-500 dark:text-gray-400 mb-12">ASO Compass はこれらをまとめて解決します</p>
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
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-4">料金プラン</h2>
          <p className="text-gray-500 dark:text-gray-400 mb-12">Webダッシュボードは無料スタート。macOSアプリは買い切り。</p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-left">

            {/* Web Plan */}
            <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow dark:shadow-black/20 border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col">
              <div className="p-8 flex-1">
                <p className="text-xs font-bold text-blue-600 dark:text-blue-400 uppercase tracking-wide mb-3">🌐 Webダッシュボード</p>
                <div className="mb-1">
                  <span className="text-4xl font-bold text-gray-900 dark:text-gray-100">¥0</span>
                  <span className="text-gray-500 dark:text-gray-400 ml-2">無料で始められます</span>
                </div>
                <p className="text-sm text-gray-400 dark:text-gray-500 mb-6">Proプランは {formatYen(PRICES.webProYearlyMonthlyEquivalent)}/月〜（年払い）</p>
                <ul className="space-y-2.5 text-sm text-gray-600 dark:text-gray-400 mb-8">
                  {[
                    { text: `${PLAN_LIMITS.freeAppLimit}アプリ・${PLAN_LIMITS.freeKeywordLimit}キーワードまで無料`, pro: false },
                    { text: 'キーワード順位の毎日追跡', pro: false },
                    { text: 'App Store アナリティクス連携', pro: false },
                    { text: '競合キーワードギャップ分析', pro: true },
                    { text: 'アプリ・キーワード数 無制限', pro: true },
                    { text: 'CSVエクスポート', pro: true },
                  ].map(({ text, pro }) => (
                    <li key={text} className="flex items-center gap-2">
                      <span className={`flex-shrink-0 ${pro ? 'text-purple-500' : 'text-green-500'}`}>✓</span>
                      <span>{text}</span>
                      {pro && <span className="ml-auto text-xs bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400 px-1.5 py-0.5 rounded font-medium">Pro</span>}
                    </li>
                  ))}
                </ul>
              </div>
              <div className="px-8 pb-8">
                <Link
                  href="/signup"
                  className="block w-full py-3.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold text-center transition-colors"
                >
                  無料で始める →
                </Link>
                <p className="text-xs text-gray-400 dark:text-gray-500 mt-2 text-center">クレジットカード不要</p>
              </div>
            </div>

            {/* macOS Plan */}
            <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow dark:shadow-black/20 border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col">
              <div className="p-8 flex-1">
                <p className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3">🍎 macOSネイティブアプリ</p>
                <div className="mb-1">
                  <span className="text-4xl font-bold text-gray-900 dark:text-gray-100">{formatYen(PRICES.macosLicenseYearly)}</span>
                  <span className="text-gray-500 dark:text-gray-400 ml-2">/ 年</span>
                </div>
                <p className="text-sm text-gray-400 dark:text-gray-500 mb-6">毎年自動更新 · 7日間返金保証</p>
                <ul className="space-y-2.5 text-sm text-gray-600 dark:text-gray-400 mb-8">
                  {[
                    '全機能アンロック',
                    'Dock・メニューバーから即アクセス',
                    '複数台Macにインストール可',
                    'macOS 13 Ventura 以降対応',
                    '購入後すぐメールでキーを送付',
                    'いつでも解約可（残期間は利用可）',
                  ].map((item) => (
                    <li key={item} className="flex items-center gap-2">
                      <span className="text-green-500 flex-shrink-0">✓</span>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
              <div className="px-8 pb-8">
                <Link
                  href="/buy"
                  className="block w-full py-3.5 bg-gray-800 dark:bg-gray-700 hover:bg-gray-700 dark:hover:bg-gray-600 text-white rounded-xl font-bold text-center transition-colors"
                >
                  今すぐ購入する →
                </Link>
                <p className="text-xs text-gray-400 dark:text-gray-500 mt-2 text-center">Stripe で安全に決済</p>
              </div>
            </div>
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
          <p className="text-blue-100 mb-10 text-lg">まずは無料で。必要になったらいつでもアップグレード。</p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/signup"
              className="w-full sm:w-auto px-10 py-4 bg-white text-blue-700 hover:bg-blue-50 rounded-xl font-bold text-lg transition-colors shadow-lg"
            >
              無料で始める →
            </Link>
            <Link
              href="/buy"
              className="w-full sm:w-auto px-10 py-4 border border-white/40 text-white hover:bg-white/10 rounded-xl font-semibold text-lg transition-colors"
            >
              macOSアプリを購入
            </Link>
          </div>
        </div>
      </section>

      {/* ── Footer ── */}
      <footer className="bg-slate-900 text-slate-400 py-10 px-4">
        <div className="max-w-5xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-sm">
          <span className="font-semibold text-white">ASO Compass</span>
          <div className="flex flex-wrap justify-center gap-x-6 gap-y-2">
            <Link href="/popular-keywords" className="hover:text-white transition-colors">人気キーワード</Link>
            <Link href="/buy" className="hover:text-white transition-colors">購入</Link>
            <Link href="/login" className="hover:text-white transition-colors">ログイン</Link>
            <Link href="/terms" className="hover:text-white transition-colors">利用規約</Link>
            <Link href="/privacy" className="hover:text-white transition-colors">プライバシーポリシー</Link>
            <Link href="/tokushoho" className="hover:text-white transition-colors">特定商取引法に基づく表記</Link>
          </div>
          <span>© 2026 ASO Compass</span>
        </div>
      </footer>
    </div>
  )
}
