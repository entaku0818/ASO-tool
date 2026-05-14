'use client'

import { useState, useEffect } from 'react'
import {
  AppMetadataVersion,
  UpsertMetadataRequest,
  getMetadata,
  upsertMetadata,
  deleteMetadata,
} from '@/lib/api'

const LOCALES = [
  { value: 'ja', label: '日本語' },
  { value: 'en-US', label: 'English (US)' },
  { value: 'zh-Hans', label: '中文(简体)' },
  { value: 'ko', label: '한국어' },
]

const CHAR_LIMITS = {
  title: 30,
  subtitle: 30,
  promotional_text: 170,
  keywords: 100,
  description: 4000,
}

function CharCounter({ value, max }: { value: string; max: number }) {
  const len = value.length
  const over = len > max
  return (
    <span className={`text-xs ${over ? 'text-red-500 font-bold' : 'text-gray-400'}`}>
      {len}/{max}
    </span>
  )
}

type FormState = {
  title: string
  subtitle: string
  description: string
  keywords: string
  promotional_text: string
}

const emptyForm: FormState = {
  title: '',
  subtitle: '',
  description: '',
  keywords: '',
  promotional_text: '',
}

type MetadataSectionProps = {
  appId: string
}

export function MetadataSection({ appId }: MetadataSectionProps) {
  const [versions, setVersions] = useState<AppMetadataVersion[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [selectedLocale, setSelectedLocale] = useState('ja')
  const [versionTag, setVersionTag] = useState('draft')
  const [form, setForm] = useState<FormState>(emptyForm)
  const [isSaving, setIsSaving] = useState(false)
  const [savedId, setSavedId] = useState<string | null>(null)

  const fetchVersions = async () => {
    try {
      const data = await getMetadata(appId)
      setVersions(data)
    } catch (e) {
      console.error('Failed to fetch metadata:', e)
    }
  }

  useEffect(() => {
    setIsLoading(true)
    fetchVersions().finally(() => setIsLoading(false))
  }, [appId])

  // Populate form when locale/version changes
  useEffect(() => {
    const existing = versions.find(v => v.locale === selectedLocale && v.version_tag === versionTag)
    if (existing) {
      setForm({
        title: existing.title ?? '',
        subtitle: existing.subtitle ?? '',
        description: existing.description ?? '',
        keywords: existing.keywords ?? '',
        promotional_text: existing.promotional_text ?? '',
      })
      setSavedId(existing.id)
    } else {
      setForm(emptyForm)
      setSavedId(null)
    }
  }, [selectedLocale, versionTag, versions])

  const handleSave = async () => {
    setIsSaving(true)
    try {
      const req: UpsertMetadataRequest = {
        locale: selectedLocale,
        version_tag: versionTag,
        title: form.title || undefined,
        subtitle: form.subtitle || undefined,
        description: form.description || undefined,
        keywords: form.keywords || undefined,
        promotional_text: form.promotional_text || undefined,
      }
      await upsertMetadata(appId, req)
      await fetchVersions()
    } catch (e) {
      alert('保存に失敗しました')
    } finally {
      setIsSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('このメタデータを削除しますか？')) return
    try {
      await deleteMetadata(appId, id)
      await fetchVersions()
    } catch (e) {
      alert('削除に失敗しました')
    }
  }

  const field = (key: keyof FormState, label: string, multiline = false) => {
    const max = CHAR_LIMITS[key]
    return (
      <div>
        <div className="flex items-center justify-between mb-1">
          <label className="text-sm font-medium text-gray-700">{label}</label>
          <CharCounter value={form[key]} max={max} />
        </div>
        {multiline ? (
          <textarea
            value={form[key]}
            onChange={(e) => setForm({ ...form, [key]: e.target.value })}
            rows={key === 'description' ? 6 : 3}
            className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
          />
        ) : (
          <input
            type="text"
            value={form[key]}
            onChange={(e) => setForm({ ...form, [key]: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        )}
      </div>
    )
  }

  if (isLoading) {
    return <div className="text-gray-500">読み込み中...</div>
  }

  return (
    <div className="space-y-4">
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h3 className="text-lg font-semibold">メタデータ管理</h3>
          <p className="text-sm text-gray-500">App Storeのタイトル・説明文・キーワードを管理</p>
        </div>

        <div className="p-4 space-y-4">
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="text-sm font-medium text-gray-700 block mb-1">言語</label>
              <select
                value={selectedLocale}
                onChange={(e) => setSelectedLocale(e.target.value)}
                className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {LOCALES.map(l => (
                  <option key={l.value} value={l.value}>{l.label}</option>
                ))}
              </select>
            </div>
            <div className="flex-1">
              <label className="text-sm font-medium text-gray-700 block mb-1">バージョンタグ</label>
              <input
                type="text"
                value={versionTag}
                onChange={(e) => setVersionTag(e.target.value)}
                placeholder="draft, v1.0, ..."
                className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          {field('title', 'タイトル（30文字以内）')}
          {field('subtitle', 'サブタイトル（30文字以内）')}
          {field('keywords', 'キーワード（100文字以内、カンマ区切り）')}
          {field('promotional_text', 'プロモーションテキスト（170文字以内）', true)}
          {field('description', '説明文（4000文字以内）', true)}

          <div className="flex justify-end gap-2 pt-2">
            {savedId && (
              <button
                onClick={() => handleDelete(savedId)}
                className="px-4 py-2 text-sm text-red-600 border border-red-300 rounded-lg hover:bg-red-50"
              >
                削除
              </button>
            )}
            <button
              onClick={handleSave}
              disabled={isSaving}
              className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {isSaving ? '保存中...' : savedId ? '更新' : '保存'}
            </button>
          </div>
        </div>
      </div>

      {versions.length > 0 && (
        <div className="bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h4 className="font-medium">保存済みバージョン ({versions.length}件)</h4>
          </div>
          <table className="w-full text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="py-2 px-4 text-left font-medium text-gray-600">言語</th>
                <th className="py-2 px-4 text-left font-medium text-gray-600">タグ</th>
                <th className="py-2 px-4 text-left font-medium text-gray-600">タイトル</th>
                <th className="py-2 px-4 text-left font-medium text-gray-600">更新日時</th>
                <th className="py-2 px-4 text-left font-medium text-gray-600"></th>
              </tr>
            </thead>
            <tbody>
              {versions.map((v) => (
                <tr
                  key={v.id}
                  className={`border-t hover:bg-gray-50 cursor-pointer ${
                    v.locale === selectedLocale && v.version_tag === versionTag ? 'bg-blue-50' : ''
                  }`}
                  onClick={() => {
                    setSelectedLocale(v.locale)
                    setVersionTag(v.version_tag)
                  }}
                >
                  <td className="py-2 px-4">{LOCALES.find(l => l.value === v.locale)?.label ?? v.locale}</td>
                  <td className="py-2 px-4 font-mono text-xs">{v.version_tag}</td>
                  <td className="py-2 px-4 text-gray-700">{v.title ?? '-'}</td>
                  <td className="py-2 px-4 text-gray-500">
                    {new Date(v.updated_at).toLocaleDateString('ja-JP')}
                  </td>
                  <td className="py-2 px-4">
                    <button
                      onClick={(e) => { e.stopPropagation(); handleDelete(v.id) }}
                      className="text-red-500 hover:text-red-700 text-xs"
                    >
                      削除
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
