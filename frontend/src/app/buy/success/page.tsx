'use client'

import Link from 'next/link'

export default function BuySuccessPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-green-50 to-white dark:from-green-900/20 dark:to-[#12161e] flex items-center justify-center px-4">
      <div className="max-w-md w-full text-center">
        <div className="text-6xl mb-6">🎉</div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-3">ご購入ありがとうございます！</h1>
        <p className="text-gray-600 dark:text-gray-400 mb-2">
          ライセンスキーをご購入時のメールアドレスに送信しました。
        </p>
        <p className="text-gray-500 dark:text-gray-400 text-sm mb-8">
          数分以内に届かない場合は迷惑メールフォルダをご確認ください。
        </p>

        <div className="bg-white dark:bg-[#12161e] rounded-xl border border-gray-200 dark:border-gray-700 p-6 mb-8 text-left space-y-4 text-sm">
          <h2 className="font-semibold text-gray-900 dark:text-gray-100">次のステップ</h2>
          <ol className="space-y-3 text-gray-600 dark:text-gray-400">
            <li className="flex gap-3">
              <span className="flex-shrink-0 w-6 h-6 bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-full flex items-center justify-center font-bold text-xs">1</span>
              <span>メールに記載のライセンスキー（ASOT-XXXX-XXXX-XXXX）を確認する</span>
            </li>
            <li className="flex gap-3">
              <span className="flex-shrink-0 w-6 h-6 bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-full flex items-center justify-center font-bold text-xs">2</span>
              <span>ASO-tool macOSアプリをダウンロードしてインストールする</span>
            </li>
            <li className="flex gap-3">
              <span className="flex-shrink-0 w-6 h-6 bg-blue-100 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-full flex items-center justify-center font-bold text-xs">3</span>
              <span>アプリを起動してライセンスキーとメールアドレスを入力してアクティベート</span>
            </li>
          </ol>
        </div>

        <Link href="/" className="text-blue-600 dark:text-blue-400 hover:underline text-sm">
          ← ダッシュボードに戻る
        </Link>
      </div>
    </div>
  )
}
