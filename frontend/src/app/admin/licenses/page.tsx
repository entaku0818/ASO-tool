'use client'

import { useState, useEffect, useCallback } from 'react'
import { getAdminLicenses, generateLicense, LicenseKey } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'
import { useRouter } from 'next/navigation'

export default function AdminLicensesPage() {
  const { user, isLoading: authLoading } = useAuth()
  const router = useRouter()
  const [licenses, setLicenses] = useState<LicenseKey[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [email, setEmail] = useState('')
  const [isGenerating, setIsGenerating] = useState(false)
  const [generatedKey, setGeneratedKey] = useState<string | null>(null)
  const [error, setError] = useState('')

  const fetchLicenses = useCallback(async () => {
    try {
      const data = await getAdminLicenses()
      setLicenses(data)
    } catch {
      setError('一覧の取得に失敗しました')
    }
  }, [])

  useEffect(() => {
    if (authLoading) return
    if (!user?.is_admin) {
      router.replace('/')
      return
    }
    setIsLoading(true)
    fetchLicenses().finally(() => setIsLoading(false))
  }, [user, authLoading, router, fetchLicenses])

  const handleGenerate = async () => {
    if (!email.trim() || !email.includes('@')) {
      setError('有効なメールアドレスを入力してください')
      return
    }
    setIsGenerating(true)
    setError('')
    setGeneratedKey(null)
    try {
      const lk = await generateLicense(email.trim())
      setGeneratedKey(lk.key)
      setEmail('')
      await fetchLicenses()
    } catch {
      setError('キーの発行に失敗しました')
    } finally {
      setIsGenerating(false)
    }
  }

  const activeCount = licenses.filter(l => l.is_active).length
  const totalCount = licenses.length

  if (authLoading || isLoading) {
    return <div className="p-8 text-gray-500 dark:text-gray-400">読み込み中...</div>
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-2">ライセンスキー管理</h1>
      <p className="text-gray-500 dark:text-gray-400 text-sm mb-6">管理者専用ページ</p>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4 mb-8">
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-4">
          <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">{totalCount}</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">発行済み合計</div>
        </div>
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-4">
          <div className="text-2xl font-bold text-green-600 dark:text-green-400">{activeCount}</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">アクティベート済み</div>
        </div>
        <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-4">
          <div className="text-2xl font-bold text-gray-600 dark:text-gray-400">{totalCount - activeCount}</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">未アクティベート</div>
        </div>
      </div>

      {/* Manual generation */}
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-6 mb-8">
        <h2 className="font-semibold mb-4">手動キー発行</h2>
        <div className="flex gap-3">
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleGenerate()}
            placeholder="購入者のメールアドレス"
            className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={handleGenerate}
            disabled={isGenerating}
            className="px-4 py-2 bg-blue-600 dark:bg-blue-500 text-white rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50"
          >
            {isGenerating ? '発行中...' : '発行'}
          </button>
        </div>
        {error && <p className="text-red-500 dark:text-red-400 text-sm mt-2">{error}</p>}
        {generatedKey && (
          <div className="mt-3 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 rounded-lg flex items-center justify-between">
            <span className="font-mono font-bold text-green-800 dark:text-green-300">{generatedKey}</span>
            <button
              onClick={() => navigator.clipboard.writeText(generatedKey)}
              className="text-xs text-green-700 dark:text-green-400 hover:underline"
            >
              コピー
            </button>
          </div>
        )}
      </div>

      {/* License list */}
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20">
        <div className="p-4 border-b dark:border-gray-700">
          <h2 className="font-semibold">全ライセンスキー ({totalCount}件)</h2>
        </div>
        {licenses.length === 0 ? (
          <p className="p-4 text-gray-500 dark:text-gray-400">キーがありません</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 dark:bg-[#0d0f14]">
                <tr>
                  <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">キー</th>
                  <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">メール</th>
                  <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">状態</th>
                  <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">アクティベート日</th>
                  <th className="py-3 px-4 text-left font-medium text-gray-600 dark:text-gray-400">発行日</th>
                </tr>
              </thead>
              <tbody>
                {licenses.map((lk) => (
                  <tr key={lk.id} className="border-t dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800">
                    <td className="py-3 px-4 font-mono text-xs">{lk.key}</td>
                    <td className="py-3 px-4 text-gray-700 dark:text-gray-300">{lk.email || '-'}</td>
                    <td className="py-3 px-4">
                      {lk.is_active ? (
                        <span className="px-2 py-0.5 bg-green-100 dark:bg-green-900/20 text-green-700 dark:text-green-400 rounded-full text-xs font-medium">有効</span>
                      ) : (
                        <span className="px-2 py-0.5 bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 rounded-full text-xs">未使用</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-gray-500 dark:text-gray-400">
                      {lk.activated_at ? new Date(lk.activated_at).toLocaleDateString('ja-JP') : '-'}
                    </td>
                    <td className="py-3 px-4 text-gray-500 dark:text-gray-400">
                      {new Date(lk.created_at).toLocaleDateString('ja-JP')}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
