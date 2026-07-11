import Link from 'next/link'

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-[#0d0f14] dark:to-[#12161e]">
      <div className="max-w-3xl mx-auto px-4 py-16">
        <Link href="/" className="text-blue-600 dark:text-blue-400 hover:underline mb-8 inline-block">
          &larr; トップに戻る
        </Link>

        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">プライバシーポリシー</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-8">最終更新日: 準備中</p>

        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4 mb-8 text-sm text-yellow-800 dark:text-yellow-300">
          このページは準備中です。正式なプライバシーポリシーの文面は法務レビュー後に掲載されます。
        </div>

        <div className="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 space-y-6">
          <p>
            ASO Tool（以下「当サービス」といいます）の運営者（以下「当社」といいます）は、
            ユーザーの個人情報の重要性を認識し、その保護の徹底を図るため、本プライバシーポリシー
            （以下「本ポリシー」といいます）を定めます。
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
