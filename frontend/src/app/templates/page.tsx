'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { getTemplates, Template, TemplateStyle } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'
import { UpgradeModal } from '@/components/UpgradeModal'
import { useUpgradeModal } from '@/hooks/useUpgradeModal'

// ── カテゴリ定義 ─────────────────────────────────────────────
const CATEGORIES = [
  { id: '',          label: 'すべて' },
  { id: 'game',      label: 'ゲーム' },
  { id: 'business',  label: 'ビジネス' },
  { id: 'education', label: '教育' },
]

const CATEGORY_STYLES: Record<string, { badge: string; accent: string }> = {
  game:      { badge: 'bg-purple-100 text-purple-700', accent: '#7C3AED' },
  business:  { badge: 'bg-blue-100 text-blue-700',     accent: '#1D4ED8' },
  education: { badge: 'bg-green-100 text-green-700',   accent: '#15803D' },
  lifestyle: { badge: 'bg-orange-100 text-orange-700', accent: '#C2410C' },
}

// ── グラデーションプレビューCSS生成 ──────────────────────────
function styleToCSS(style: TemplateStyle): React.CSSProperties {
  if (style.bg_gradient_from && style.bg_gradient_to) {
    const dir =
      style.bg_gradient_dir === 'lr'   ? 'to right' :
      style.bg_gradient_dir === 'tlbr' ? 'to bottom right' :
                                          'to bottom'
    return { background: `linear-gradient(${dir}, ${style.bg_gradient_from}, ${style.bg_gradient_to})` }
  }
  return { background: style.bg_color ?? '#4F46E5' }
}

// ── テンプレートカード ────────────────────────────────────────
function TemplateCard({
  template,
  isPro,
  onUse,
  onUnlock,
}: {
  template: Template
  isPro: boolean
  onUse: (t: Template) => void
  onUnlock: () => void
}) {
  const catStyle = CATEGORY_STYLES[template.category] ?? CATEGORY_STYLES.business
  const locked = template.is_pro && !isPro

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      {/* Thumbnail */}
      <div className="relative aspect-video overflow-hidden">
        <div className="w-full h-full" style={styleToCSS(template.style)} />
        {locked && (
          <div className="absolute inset-0 bg-black/40 backdrop-blur-sm flex flex-col items-center justify-center gap-1">
            <span className="text-2xl">🔒</span>
            <span className="text-white text-xs font-semibold">Pro プラン限定</span>
          </div>
        )}
      </div>

      {/* Body */}
      <div className="p-4">
        <span className={`inline-block text-xs font-semibold px-2 py-0.5 rounded-full mb-2 ${catStyle.badge}`}>
          {CATEGORIES.find(c => c.id === template.category)?.label ?? template.category}
        </span>
        <p className="font-semibold text-gray-900 text-sm mb-1">{template.name}</p>
        <p className="text-xs text-gray-500 mb-4">{template.description}</p>

        <div className="border-t border-gray-100 pt-3 flex gap-2">
          {locked ? (
            <button
              onClick={onUnlock}
              className="w-full py-2 text-sm font-semibold text-blue-600 border border-blue-200 rounded-lg hover:bg-blue-50 transition-colors"
            >
              Pro でアンロック
            </button>
          ) : (
            <>
              <button
                onClick={() => onUse(template)}
                className="flex-1 py-2 text-sm font-semibold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors"
              >
                このテンプレートを使う →
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

// ── ページ本体 ────────────────────────────────────────────────
export default function TemplatesPage() {
  const router = useRouter()
  const { user } = useAuth()
  const isPro = user?.is_pro ?? false
  const upgradeModal = useUpgradeModal()

  const [selectedCategory, setSelectedCategory] = useState('')
  const [templates, setTemplates] = useState<Template[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setIsLoading(true)
    setError(null)
    getTemplates(selectedCategory || undefined)
      .then(setTemplates)
      .catch(() => setError('テンプレートの取得に失敗しました'))
      .finally(() => setIsLoading(false))
  }, [selectedCategory])

  const handleUse = (template: Template) => {
    // Store style in sessionStorage — ScreenshotGenerator reads it on mount
    try {
      const payload: Record<string, string> = {}
      const s = template.style
      if (s.bg_gradient_from && s.bg_gradient_to) {
        payload.bg_gradient_from = s.bg_gradient_from
        payload.bg_gradient_to = s.bg_gradient_to
        if (s.bg_gradient_dir) payload.bg_gradient_dir = s.bg_gradient_dir
      } else if (s.bg_color) {
        payload.bg_color = s.bg_color
      }
      if (s.text_color) payload.text_color = s.text_color
      if (s.image_align) payload.image_align = s.image_align
      if (template.device) payload.device = template.device
      sessionStorage.setItem('pending_template', JSON.stringify(payload))
    } catch {
      // sessionStorage unavailable — fall through to navigation without preset
    }
    // Navigate to home; user selects an app → opens Screenshot tab where the style is applied
    router.push('/')
  }

  return (
    <div className="max-w-5xl mx-auto">
      {/* Page header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">テンプレートライブラリ</h1>
          <p className="text-sm text-gray-500 mt-1">
            App Store スクリーンショット用テンプレートを選んで始めよう
          </p>
        </div>
      </div>

      {/* Category filter */}
      <div className="flex gap-2 mb-6 overflow-x-auto pb-1">
        {CATEGORIES.map(cat => (
          <button
            key={cat.id}
            onClick={() => setSelectedCategory(cat.id)}
            className={`px-4 py-2 text-sm font-medium rounded-full whitespace-nowrap transition-colors ${
              selectedCategory === cat.id
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            {cat.label}
          </button>
        ))}
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-gray-500">読み込み中...</p>
        </div>
      ) : error ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-red-500">{error}</p>
        </div>
      ) : templates.length === 0 ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-gray-500">テンプレートがありません</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {templates.map(t => (
            <TemplateCard
              key={t.id}
              template={t}
              isPro={isPro}
              onUse={handleUse}
              onUnlock={() => upgradeModal.open('template_library')}
            />
          ))}
        </div>
      )}

      <UpgradeModal
        isOpen={upgradeModal.isOpen}
        onClose={upgradeModal.close}
        triggerFeature="テンプレートライブラリ（Pro限定）"
      />
    </div>
  )
}
