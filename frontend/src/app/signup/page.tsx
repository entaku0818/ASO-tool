'use client'

import { useState, FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/contexts/AuthContext'
import { register } from '@/lib/api'

const FREE_FEATURES = [
  '1アプリ・10キーワードまで',
  'キーワード順位の毎日自動追跡',
  'App Store アナリティクス連携',
  'おすすめキーワード提案',
]

const PRO_FEATURES = [
  'アプリ・キーワード数 無制限',
  '競合キーワードギャップ分析',
  'CSVエクスポート',
  '多言語スクリーンショット一括生成',
]

export default function SignupPage() {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()
  const { login } = useAuth()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')

    if (password.length < 8) {
      setError('パスワードは8文字以上で入力してください')
      return
    }

    setIsLoading(true)
    try {
      const { token } = await register(email, password, name)
      localStorage.setItem('auth_token', token)
      await login(email, password)
      router.push('/dashboard?new=1')
    } catch (err) {
      setError(err instanceof Error ? err.message : '登録に失敗しました')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-[#0d0f14] px-4 py-12">
      <div className="w-full max-w-4xl flex flex-col lg:flex-row gap-8 items-start">

        {/* Left: plan summary */}
        <div className="flex-1 w-full">
          <div className="mb-8">
            <span className="inline-block px-3 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 text-xs font-bold rounded-full uppercase tracking-wide mb-3">
              完全無料でスタート
            </span>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
              ASO作業を、もっとスマートに。
            </h2>
            <p className="text-gray-500 dark:text-gray-400 text-sm leading-relaxed">
              クレジットカード不要。今すぐ無料プランで始めて、
              必要なときにいつでもProへアップグレードできます。
            </p>
          </div>

          {/* Free plan */}
          <div className="bg-white dark:bg-[#12161e] rounded-xl border border-gray-200 dark:border-gray-700 p-5 mb-3">
            <p className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3">
              🆓 無料プラン — ずっと無料
            </p>
            <ul className="space-y-2">
              {FREE_FEATURES.map((f) => (
                <li key={f} className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                  <span className="text-green-500 flex-shrink-0">✓</span>
                  {f}
                </li>
              ))}
            </ul>
          </div>

          {/* Pro plan */}
          <div className="bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 rounded-xl border border-blue-100 dark:border-blue-800 p-5">
            <div className="flex items-center justify-between mb-3">
              <p className="text-xs font-bold text-blue-700 dark:text-blue-300 uppercase tracking-wide">
                ✦ Proプラン
              </p>
              <span className="text-xs bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 px-2 py-0.5 rounded-full font-medium">
                ¥1,650/月〜
              </span>
            </div>
            <ul className="space-y-2">
              {PRO_FEATURES.map((f) => (
                <li key={f} className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                  <span className="text-purple-500 flex-shrink-0">✓</span>
                  {f}
                </li>
              ))}
            </ul>
            <p className="text-xs text-gray-400 dark:text-gray-500 mt-3">登録後いつでもアップグレード可</p>
          </div>
        </div>

        {/* Right: signup form */}
        <div className="flex-1 w-full max-w-md lg:max-w-none">
          <div className="bg-white dark:bg-[#12161e] rounded-xl shadow dark:shadow-black/20 p-8">
            <div className="mb-6">
              <h1 className="text-2xl font-extrabold text-gray-900 dark:text-gray-100">
                無料アカウントを作成
              </h1>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                30秒で登録完了
              </p>
            </div>

            <form className="space-y-4" onSubmit={handleSubmit}>
              {error && (
                <div className="rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 p-3">
                  <p className="text-sm text-red-700 dark:text-red-400">{error}</p>
                </div>
              )}

              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  お名前
                </label>
                <input
                  id="name"
                  type="text"
                  autoComplete="name"
                  required
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent sm:text-sm"
                  placeholder="山田 太郎"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  メールアドレス
                </label>
                <input
                  id="email"
                  type="email"
                  autoComplete="email"
                  required
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent sm:text-sm"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  パスワード
                </label>
                <input
                  id="password"
                  type="password"
                  autoComplete="new-password"
                  required
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg placeholder-gray-400 dark:placeholder-gray-500 text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent sm:text-sm"
                  placeholder="8文字以上"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed mt-2"
              >
                {isLoading ? '登録中...' : '無料アカウントを作成'}
              </button>
            </form>

            <p className="text-center text-sm text-gray-500 dark:text-gray-400 mt-5">
              すでにアカウントをお持ちの方は{' '}
              <Link href="/login" className="text-blue-600 dark:text-blue-400 hover:underline font-medium">
                ログイン
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
