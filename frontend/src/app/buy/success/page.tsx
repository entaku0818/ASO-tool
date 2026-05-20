'use client'

import Link from 'next/link'

const DOWNLOAD_URL = 'https://github.com/entaku0818/ASO-tool/releases/latest/download/ASO-tool.dmg'

const STEPS = [
  {
    num: 1,
    title: 'メールを確認する',
    desc: 'ライセンスキー（ASOT-XXXX-XXXX-XXXX）が届いています。数分かかる場合は迷惑メールフォルダもご確認ください。',
    action: null,
  },
  {
    num: 2,
    title: 'macOSアプリをダウンロード',
    desc: 'macOS 13 Ventura 以降に対応しています。',
    action: { label: 'ASO-tool.dmg をダウンロード', href: DOWNLOAD_URL },
  },
  {
    num: 3,
    title: 'ライセンスを有効化',
    desc: 'アプリを起動して、メールアドレスとライセンスキーを入力してアクティベートしてください。',
    action: null,
  },
]

export default function BuySuccessPage() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-[#0d0f14] flex items-center justify-center px-4 py-12">
      <div className="max-w-lg w-full">

        {/* Hero */}
        <div className="text-center mb-8">
          <div className="text-6xl mb-4">🎉</div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
            ご購入ありがとうございます！
          </h1>
          <p className="text-gray-500 dark:text-gray-400 text-sm">
            ライセンスキーをご購入時のメールアドレスに送信しました
          </p>
        </div>

        {/* Steps */}
        <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow dark:shadow-black/20 border border-gray-200 dark:border-gray-700 p-6 mb-4">
          <h2 className="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-5">
            セットアップ手順
          </h2>
          <ol className="space-y-5">
            {STEPS.map((step) => (
              <li key={step.num} className="flex gap-4">
                <span className="flex-shrink-0 w-7 h-7 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 rounded-full flex items-center justify-center font-bold text-sm">
                  {step.num}
                </span>
                <div>
                  <p className="font-semibold text-gray-900 dark:text-gray-100 text-sm mb-0.5">{step.title}</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400 leading-relaxed">{step.desc}</p>
                  {step.action && (
                    <a
                      href={step.action.href}
                      className="inline-flex items-center gap-1.5 mt-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition-colors"
                    >
                      ↓ {step.action.label}
                    </a>
                  )}
                </div>
              </li>
            ))}
          </ol>
        </div>

        {/* Web dashboard cross-sell */}
        <div className="bg-blue-50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-800 rounded-2xl p-5 mb-6">
          <p className="text-sm font-semibold text-blue-800 dark:text-blue-300 mb-1">
            メールを待つ間に…
          </p>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
            Webダッシュボードでもキーワード追跡や競合分析が使えます。アカウント登録は無料です。
          </p>
          <Link
            href="/signup"
            className="inline-block px-4 py-2 bg-white dark:bg-[#12161e] border border-blue-200 dark:border-blue-700 text-blue-700 dark:text-blue-300 text-sm font-semibold rounded-lg hover:bg-blue-50 dark:hover:bg-blue-900/20 transition-colors"
          >
            無料Webアカウントを作成 →
          </Link>
        </div>

        <p className="text-center text-xs text-gray-400 dark:text-gray-500">
          問題が発生した場合はサポートまでお問い合わせください
        </p>
      </div>
    </div>
  )
}
