import Link from 'next/link'

const ROWS: { label: string; value: string }[] = [
  { label: '販売事業者', value: '準備中' },
  { label: '運営統括責任者', value: '準備中' },
  { label: '所在地', value: '準備中（請求があった場合、遅滞なく開示します）' },
  { label: '連絡先電話番号', value: '準備中（請求があった場合、遅滞なく開示します）' },
  { label: '連絡先メールアドレス', value: '準備中' },
  { label: '販売価格', value: '各商品ページに記載の価格（消費税込）' },
  { label: '商品代金以外の必要料金', value: '準備中' },
  { label: 'お支払い方法', value: 'クレジットカード決済（Stripe）' },
  { label: 'お支払い時期', value: 'ご注文確定時' },
  { label: '商品の引き渡し時期', value: '決済完了後、即時にご利用いただけます' },
  { label: '返品・キャンセルについて', value: '準備中' },
]

export default function TokushohoPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-[#0d0f14] dark:to-[#12161e]">
      <div className="max-w-3xl mx-auto px-4 py-16">
        <Link href="/" className="text-blue-600 dark:text-blue-400 hover:underline mb-8 inline-block">
          &larr; トップに戻る
        </Link>

        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">特定商取引法に基づく表記</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-8">最終更新日: 準備中</p>

        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4 mb-8 text-sm text-yellow-800 dark:text-yellow-300">
          このページは準備中です。「準備中」となっている項目は事業者情報の確定後に正式な内容へ差し替えます。
        </div>

        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 overflow-hidden">
          <table className="w-full text-sm">
            <tbody>
              {ROWS.map((row) => (
                <tr key={row.label} className="border-b dark:border-gray-700 last:border-b-0">
                  <th className="text-left align-top py-3 px-4 w-1/3 font-medium text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-[#0d0f14]">
                    {row.label}
                  </th>
                  <td className="py-3 px-4 text-gray-800 dark:text-gray-200">{row.value}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
