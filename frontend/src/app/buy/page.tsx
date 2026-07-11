'use client'

import { useState } from 'react'
import { createLicenseCheckout } from '@/lib/api'
import { PRICES, formatYen } from '@/lib/plans'

const FEATURES = [
  { icon: '📊', title: 'ランキング推移チャート', desc: '最大90日のキーワード順位推移を複数キーワード同時比較' },
  { icon: '🎯', title: '競合キーワードギャップ', desc: '競合が上位・自社が圏外のキーワードを自動検出' },
  { icon: '📝', title: 'App Storeメタデータ管理', desc: 'タイトル・説明文・キーワードをロケール/バージョン別に管理' },
  { icon: '⚡', title: 'ネイティブmacOSアプリ', desc: 'メニューバーから素早くアクセス。オフライン閲覧対応' },
  { icon: '🔑', title: '買い切りライセンス', desc: '月額料金なし。1ライセンスで永続利用' },
  { icon: '🔄', title: '無料アップデート', desc: '将来の機能追加も追加費用なし' },
]

export default function BuyPage() {
  const [email, setEmail] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const handleCheckout = async () => {
    if (!email.trim() || !email.includes('@')) {
      setError('有効なメールアドレスを入力してください')
      return
    }
    setIsLoading(true)
    setError('')
    try {
      const { url } = await createLicenseCheckout(email.trim())
      window.location.href = url
    } catch (e) {
      setError('購入処理の開始に失敗しました。しばらく経ってから再試行してください。')
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-[#0d0f14] dark:to-[#12161e]">
      <div className="max-w-3xl mx-auto px-4 py-16">

        {/* Hero */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            ASO-tool macOS版
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-400">
            App Storeキーワード最適化をデスクトップから
          </p>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-12">
          {FEATURES.map((f) => (
            <div key={f.title} className="bg-white dark:bg-[#12161e] rounded-xl p-5 shadow-sm dark:shadow-black/20 border border-gray-100 dark:border-gray-800">
              <div className="text-2xl mb-2">{f.icon}</div>
              <h3 className="font-semibold text-gray-900 dark:text-gray-100 mb-1">{f.title}</h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">{f.desc}</p>
            </div>
          ))}
        </div>

        {/* Pricing card */}
        <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow-lg border border-gray-200 dark:border-gray-700 p-8 max-w-md mx-auto">
          <div className="text-center mb-6">
            <div className="text-5xl font-bold text-gray-900 dark:text-gray-100 mb-1">{formatYen(PRICES.macosLicenseYearly)}</div>
            <div className="text-gray-500 dark:text-gray-400">年額・毎年自動更新</div>
          </div>

          <ul className="space-y-2 mb-8 text-sm text-gray-600 dark:text-gray-400">
            <li className="flex items-center gap-2"><span className="text-green-500 dark:text-green-400">✓</span>1ライセンス = 1アカウント</li>
            <li className="flex items-center gap-2"><span className="text-green-500 dark:text-green-400">✓</span>macOS 13 Ventura以降対応</li>
            <li className="flex items-center gap-2"><span className="text-green-500 dark:text-green-400">✓</span>購入後すぐにメールでキーを送付</li>
            <li className="flex items-center gap-2"><span className="text-green-500 dark:text-green-400">✓</span>いつでも解約可能・残期間は利用可</li>
            <li className="flex items-center gap-2"><span className="text-green-500 dark:text-green-400">✓</span>クレジットカード決済（Stripe）</li>
          </ul>

          <div className="space-y-3">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleCheckout()}
              placeholder="購入者のメールアドレス"
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            />
            {error && <p className="text-red-500 dark:text-red-400 text-sm">{error}</p>}
            <button
              onClick={handleCheckout}
              disabled={isLoading}
              className="w-full py-3 bg-blue-600 dark:bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              {isLoading ? '処理中...' : '今すぐ購入 →'}
            </button>
          </div>

          <p className="text-xs text-gray-400 dark:text-gray-500 text-center mt-4">
            決済はStripeで安全に処理されます
          </p>
        </div>

        {/* FAQ */}
        <div className="mt-12 space-y-4 text-sm text-gray-600 dark:text-gray-400">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">よくある質問</h2>
          <div>
            <p className="font-medium text-gray-800 dark:text-gray-200">購入後はどうすれば使えますか？</p>
            <p>購入完了後、入力したメールアドレスにライセンスキーが届きます。macOSアプリを起動してキーとメールアドレスを入力するとアクティベートできます。</p>
          </div>
          <div>
            <p className="font-medium text-gray-800 dark:text-gray-200">アプリはどこからダウンロードできますか？</p>
            <p>購入後にダウンロードリンクをメールでお送りします。</p>
          </div>
          <div>
            <p className="font-medium text-gray-800 dark:text-gray-200">解約はできますか？</p>
            <p>いつでも解約できます。解約後はその年の有効期限まで引き続き利用できます。</p>
          </div>
          <div>
            <p className="font-medium text-gray-800 dark:text-gray-200">返金はできますか？</p>
            <p>初回購入後7日以内であれば全額返金対応しています。サポートまでご連絡ください。</p>
          </div>
        </div>
      </div>
    </div>
  )
}
