import Link from 'next/link'

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-[#0d0f14] dark:to-[#12161e]">
      <div className="max-w-3xl mx-auto px-4 py-16">
        <Link href="/" className="text-blue-600 dark:text-blue-400 hover:underline mb-8 inline-block">
          &larr; トップに戻る
        </Link>

        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">利用規約</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-8">最終更新日: 準備中</p>

        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4 mb-8 text-sm text-yellow-800 dark:text-yellow-300">
          このページは準備中です。正式な利用規約の文面は法務レビュー後に掲載されます。
        </div>

        <div className="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 space-y-6">
          <p>
            本利用規約（以下「本規約」といいます）は、ASO Tool（以下「当サービス」といいます）の
            提供条件および当サービスの利用に関する当サービス運営者（以下「当社」といいます）と
            登録ユーザー（以下「ユーザー」といいます）との間の権利義務関係を定めるものです。
          </p>
          <p>
            正式な条項確定までの間、本サービスのご利用は本ページの内容を暫定的な取り扱いとして
            ご案内するものであり、内容は今後変更される可能性があります。
          </p>
        </div>
      </div>
    </div>
  )
}
