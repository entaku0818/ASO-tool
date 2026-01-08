'use client'

import { useState, useEffect } from 'react'
import {
  Screenshot,
  DeviceType,
  getScreenshots,
  createScreenshot,
  updateScreenshot,
  deleteScreenshot,
  reorderScreenshots,
} from '@/lib/api'

type ScreenshotSectionProps = {
  appId: string
  platform: 'ios' | 'android'
}

const DEVICE_TYPES: { value: DeviceType; label: string }[] = [
  { value: 'iphone', label: 'iPhone' },
  { value: 'ipad', label: 'iPad' },
  { value: 'android', label: 'Android' },
  { value: 'tablet', label: 'Tablet' },
]

const LOCALES = [
  { value: 'ja-JP', label: '日本語' },
  { value: 'en-US', label: 'English (US)' },
  { value: 'en-GB', label: 'English (UK)' },
  { value: 'zh-CN', label: '中国語 (簡体)' },
  { value: 'zh-TW', label: '中国語 (繁体)' },
  { value: 'ko-KR', label: '韓国語' },
]

export function ScreenshotSection({ appId, platform }: ScreenshotSectionProps) {
  const [screenshots, setScreenshots] = useState<Screenshot[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [showAddForm, setShowAddForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [isAdding, setIsAdding] = useState(false)
  const [selectedDeviceFilter, setSelectedDeviceFilter] = useState<DeviceType | ''>('')
  const [newScreenshot, setNewScreenshot] = useState({
    url: '',
    device_type: platform === 'ios' ? 'iphone' : 'android' as DeviceType,
    locale: 'ja-JP',
    description: '',
  })
  const [editForm, setEditForm] = useState({
    url: '',
    device_type: 'iphone' as DeviceType,
    locale: 'ja-JP',
    description: '',
  })

  const fetchScreenshots = async () => {
    try {
      setIsLoading(true)
      const deviceType = selectedDeviceFilter || undefined
      const data = await getScreenshots(appId, deviceType)
      setScreenshots(data)
    } catch (e) {
      console.error('Failed to fetch screenshots:', e)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchScreenshots()
  }, [appId, selectedDeviceFilter])

  const handleAdd = async () => {
    if (!newScreenshot.url.trim()) return
    setIsAdding(true)
    try {
      await createScreenshot(appId, {
        url: newScreenshot.url.trim(),
        device_type: newScreenshot.device_type,
        locale: newScreenshot.locale,
        description: newScreenshot.description.trim() || undefined,
      })
      setNewScreenshot({
        url: '',
        device_type: platform === 'ios' ? 'iphone' : 'android',
        locale: 'ja-JP',
        description: '',
      })
      setShowAddForm(false)
      fetchScreenshots()
    } catch (e) {
      alert('スクリーンショットの追加に失敗しました')
    } finally {
      setIsAdding(false)
    }
  }

  const handleEdit = (screenshot: Screenshot) => {
    setEditingId(screenshot.id)
    setEditForm({
      url: screenshot.url,
      device_type: screenshot.device_type,
      locale: screenshot.locale,
      description: screenshot.description || '',
    })
  }

  const handleUpdate = async (id: string) => {
    try {
      await updateScreenshot(appId, id, {
        url: editForm.url,
        device_type: editForm.device_type,
        locale: editForm.locale,
        description: editForm.description || undefined,
      })
      setEditingId(null)
      fetchScreenshots()
    } catch (e) {
      alert('更新に失敗しました')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('このスクリーンショットを削除しますか？')) return
    try {
      await deleteScreenshot(appId, id)
      fetchScreenshots()
    } catch (e) {
      alert('削除に失敗しました')
    }
  }

  const handleMoveUp = async (index: number) => {
    if (index === 0) return
    const newOrder = [...screenshots]
    const temp = newOrder[index]
    newOrder[index] = newOrder[index - 1]
    newOrder[index - 1] = temp
    try {
      await reorderScreenshots(appId, newOrder.map(s => s.id))
      fetchScreenshots()
    } catch (e) {
      alert('並び替えに失敗しました')
    }
  }

  const handleMoveDown = async (index: number) => {
    if (index === screenshots.length - 1) return
    const newOrder = [...screenshots]
    const temp = newOrder[index]
    newOrder[index] = newOrder[index + 1]
    newOrder[index + 1] = temp
    try {
      await reorderScreenshots(appId, newOrder.map(s => s.id))
      fetchScreenshots()
    } catch (e) {
      alert('並び替えに失敗しました')
    }
  }

  const getDeviceLabel = (deviceType: DeviceType) => {
    return DEVICE_TYPES.find(d => d.value === deviceType)?.label || deviceType
  }

  const getLocaleLabel = (locale: string) => {
    return LOCALES.find(l => l.value === locale)?.label || locale
  }

  if (isLoading) {
    return <div className="text-gray-500">読み込み中...</div>
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-4 border-b flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">スクリーンショット ({screenshots.length}件)</h3>
          <p className="text-sm text-gray-500">アプリのスクリーンショットを管理</p>
        </div>
        <div className="flex gap-2 items-center">
          <select
            value={selectedDeviceFilter}
            onChange={(e) => setSelectedDeviceFilter(e.target.value as DeviceType | '')}
            className="px-3 py-1.5 text-sm border rounded"
          >
            <option value="">すべてのデバイス</option>
            {DEVICE_TYPES.map((d) => (
              <option key={d.value} value={d.value}>{d.label}</option>
            ))}
          </select>
          <button
            onClick={() => setShowAddForm(!showAddForm)}
            className="px-3 py-1.5 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            {showAddForm ? 'キャンセル' : '追加'}
          </button>
        </div>
      </div>

      {showAddForm && (
        <div className="p-4 bg-gray-50 border-b">
          <div className="space-y-3">
            <input
              type="url"
              value={newScreenshot.url}
              onChange={(e) => setNewScreenshot({ ...newScreenshot, url: e.target.value })}
              placeholder="スクリーンショットURL"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="flex gap-3">
              <select
                value={newScreenshot.device_type}
                onChange={(e) => setNewScreenshot({ ...newScreenshot, device_type: e.target.value as DeviceType })}
                className="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {DEVICE_TYPES.map((d) => (
                  <option key={d.value} value={d.value}>{d.label}</option>
                ))}
              </select>
              <select
                value={newScreenshot.locale}
                onChange={(e) => setNewScreenshot({ ...newScreenshot, locale: e.target.value })}
                className="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {LOCALES.map((l) => (
                  <option key={l.value} value={l.value}>{l.label}</option>
                ))}
              </select>
            </div>
            <input
              type="text"
              value={newScreenshot.description}
              onChange={(e) => setNewScreenshot({ ...newScreenshot, description: e.target.value })}
              placeholder="説明 (任意)"
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={handleAdd}
              disabled={isAdding || !newScreenshot.url.trim()}
              className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {isAdding ? '追加中...' : 'スクリーンショットを追加'}
            </button>
          </div>
        </div>
      )}

      {screenshots.length === 0 ? (
        <p className="p-4 text-gray-500">スクリーンショットが登録されていません</p>
      ) : (
        <div className="p-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {screenshots.map((screenshot, index) => (
              <div key={screenshot.id} className="border rounded-lg overflow-hidden bg-white shadow-sm">
                {editingId === screenshot.id ? (
                  <div className="p-3 space-y-2">
                    <input
                      type="url"
                      value={editForm.url}
                      onChange={(e) => setEditForm({ ...editForm, url: e.target.value })}
                      className="w-full px-2 py-1 text-sm border rounded"
                    />
                    <select
                      value={editForm.device_type}
                      onChange={(e) => setEditForm({ ...editForm, device_type: e.target.value as DeviceType })}
                      className="w-full px-2 py-1 text-sm border rounded"
                    >
                      {DEVICE_TYPES.map((d) => (
                        <option key={d.value} value={d.value}>{d.label}</option>
                      ))}
                    </select>
                    <select
                      value={editForm.locale}
                      onChange={(e) => setEditForm({ ...editForm, locale: e.target.value })}
                      className="w-full px-2 py-1 text-sm border rounded"
                    >
                      {LOCALES.map((l) => (
                        <option key={l.value} value={l.value}>{l.label}</option>
                      ))}
                    </select>
                    <input
                      type="text"
                      value={editForm.description}
                      onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                      placeholder="説明"
                      className="w-full px-2 py-1 text-sm border rounded"
                    />
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleUpdate(screenshot.id)}
                        className="flex-1 px-2 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
                      >
                        保存
                      </button>
                      <button
                        onClick={() => setEditingId(null)}
                        className="flex-1 px-2 py-1 text-sm bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
                      >
                        キャンセル
                      </button>
                    </div>
                  </div>
                ) : (
                  <>
                    <div className="aspect-[9/16] bg-gray-100 relative">
                      <img
                        src={screenshot.url}
                        alt={screenshot.description || 'Screenshot'}
                        className="w-full h-full object-contain"
                        onError={(e) => {
                          const target = e.target as HTMLImageElement
                          target.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y="50" x="50" text-anchor="middle" dominant-baseline="middle" font-size="14" fill="%23999">No Image</text></svg>'
                        }}
                      />
                      <div className="absolute top-2 left-2 flex gap-1">
                        <span className="px-2 py-0.5 text-xs bg-blue-600 text-white rounded">
                          {getDeviceLabel(screenshot.device_type)}
                        </span>
                        <span className="px-2 py-0.5 text-xs bg-gray-600 text-white rounded">
                          {getLocaleLabel(screenshot.locale)}
                        </span>
                      </div>
                      <div className="absolute top-2 right-2 flex gap-1">
                        <span className="px-2 py-0.5 text-xs bg-gray-800 text-white rounded">
                          #{index + 1}
                        </span>
                      </div>
                    </div>
                    <div className="p-3">
                      {screenshot.description && (
                        <p className="text-sm text-gray-600 mb-2 line-clamp-2">{screenshot.description}</p>
                      )}
                      <div className="flex items-center justify-between">
                        <div className="flex gap-1">
                          <button
                            onClick={() => handleMoveUp(index)}
                            disabled={index === 0}
                            className="p-1 text-gray-500 hover:text-gray-700 disabled:opacity-30"
                            title="上に移動"
                          >
                            ↑
                          </button>
                          <button
                            onClick={() => handleMoveDown(index)}
                            disabled={index === screenshots.length - 1}
                            className="p-1 text-gray-500 hover:text-gray-700 disabled:opacity-30"
                            title="下に移動"
                          >
                            ↓
                          </button>
                        </div>
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleEdit(screenshot)}
                            className="text-blue-500 hover:text-blue-700 text-sm"
                          >
                            編集
                          </button>
                          <button
                            onClick={() => handleDelete(screenshot.id)}
                            className="text-red-500 hover:text-red-700 text-sm"
                          >
                            削除
                          </button>
                        </div>
                      </div>
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
