import type { Metadata } from 'next'
import Link from 'next/link'
import './globals.css'

export const metadata: Metadata = {
  title: 'ASO Tool',
  description: 'App Store Optimization Tool',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body className="bg-gray-50 min-h-screen">
        <header className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
            <Link href="/" className="text-xl font-bold text-gray-900 hover:text-gray-700">
              ASO Tool
            </Link>
            <nav className="flex gap-4">
              <Link href="/" className="text-gray-600 hover:text-gray-900">
                ダッシュボード
              </Link>
              <Link href="/rankings" className="text-gray-600 hover:text-gray-900">
                ランキング
              </Link>
              <Link href="/keywords" className="text-gray-600 hover:text-gray-900">
                キーワード検索
              </Link>
            </nav>
          </div>
        </header>
        <main className="max-w-7xl mx-auto px-4 py-6">
          {children}
        </main>
      </body>
    </html>
  )
}
