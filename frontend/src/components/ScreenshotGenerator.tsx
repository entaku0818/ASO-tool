'use client'

import { useState, useRef, useCallback } from 'react'
import JSZip from 'jszip'
import { generateScreenshots } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'
import { UpgradeModal } from '@/components/UpgradeModal'
import { useUpgradeModal } from '@/hooks/useUpgradeModal'

const LANGUAGES = [
  { code: 'ja', label: '日本語' },
  { code: 'en', label: 'English' },
  { code: 'zh-Hans', label: '中文(簡)' },
  { code: 'ko', label: '한국어' },
  { code: 'fr', label: 'Français' },
  { code: 'de', label: 'Deutsch' },
]

const DEVICE_PRESETS = [
  { id: 'iphone67', label: 'iPhone 6.7"', width: 393, height: 852, frameColor: '#1a1a1a', cornerRadius: 55 },
  { id: 'iphone65', label: 'iPhone 6.5"', width: 414, height: 896, frameColor: '#1a1a1a', cornerRadius: 55 },
  { id: 'ipad', label: 'iPad 12.9"', width: 1024, height: 1366, frameColor: '#1a1a1a', cornerRadius: 20 },
]

const BG_PRESETS = [
  { id: 'ocean',    label: 'Ocean',    from: '#667eea', to: '#764ba2', dir: 'tb' as const },
  { id: 'midnight', label: 'Midnight', from: '#0f0c29', to: '#302b63', dir: 'tb' as const },
  { id: 'sunset',   label: 'Sunset',   from: '#f093fb', to: '#f5576c', dir: 'tb' as const },
  { id: 'forest',   label: 'Forest',   from: '#134e5e', to: '#71b280', dir: 'tb' as const },
  { id: 'gold',     label: 'Gold',     from: '#f7971e', to: '#ffd200', dir: 'tb' as const },
  { id: 'rose',     label: 'Rose',     from: '#ee0979', to: '#ff6a00', dir: 'tb' as const },
  { id: 'sky',      label: 'Sky',      from: '#4facfe', to: '#00f2fe', dir: 'tb' as const },
  { id: 'dark',     label: 'Dark',     from: '#1a1a2e', to: '#16213e', dir: 'tb' as const },
]

type Captions = Record<string, string>
type BgType = 'solid' | 'gradient' | 'preset'
type GradDir = 'tb' | 'lr' | 'tlbr'
type ImageAlign = 'center' | 'bottom'

interface ScreenshotGeneratorProps {
  appName?: string
  appId?: string
}

export function ScreenshotGenerator({ appName, appId }: ScreenshotGeneratorProps) {
  const { user } = useAuth()
  const isPro = user?.plan === 'pro'
  const upgradeModal = useUpgradeModal()
  const [imageURL, setImageURL] = useState<string>('')
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [captions, setCaptions] = useState<Captions>({ ja: '', en: '', 'zh-Hans': '', ko: '', fr: '', de: '' })
  const [selectedLang, setSelectedLang] = useState('ja')
  const [bgColor, setBgColor] = useState('#4F46E5')
  const [bgType, setBgType] = useState<BgType>('solid')
  const [bgGradFrom, setBgGradFrom] = useState('#667eea')
  const [bgGradTo, setBgGradTo] = useState('#764ba2')
  const [bgGradDir, setBgGradDir] = useState<GradDir>('tb')
  const [bgPresetId, setBgPresetId] = useState('ocean')
  const [textColor, setTextColor] = useState('#FFFFFF')
  const [device, setDevice] = useState(DEVICE_PRESETS[0])
  const [imageAlign, setImageAlign] = useState<ImageAlign>('bottom')
  const [isGenerating, setIsGenerating] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setImageFile(file)
    const reader = new FileReader()
    reader.onload = (ev) => setImageURL(ev.target?.result as string)
    reader.readAsDataURL(file)
  }

  // Resolves the effective gradient colors (from/to) based on current bg settings
  const getEffectiveBg = useCallback((): { from: string; to: string; dir: GradDir } | null => {
    if (bgType === 'gradient') return { from: bgGradFrom, to: bgGradTo, dir: bgGradDir }
    if (bgType === 'preset') {
      const p = BG_PRESETS.find(p => p.id === bgPresetId) ?? BG_PRESETS[0]
      return { from: p.from, to: p.to, dir: p.dir }
    }
    return null
  }, [bgType, bgGradFrom, bgGradTo, bgGradDir, bgPresetId])

  const drawFrame = useCallback((canvas: HTMLCanvasElement, lang: string): Promise<void> => {
    return new Promise((resolve) => {
      const ctx = canvas.getContext('2d')
      if (!ctx) return resolve()

      const PADDING = 40
      // bottom-align uses a taller caption area so the device clearly "sits at the bottom"
      const CAPTION_H = imageAlign === 'bottom' ? 220 : 120
      const W = device.width + PADDING * 2
      const H = device.height + PADDING * 2 + CAPTION_H

      canvas.width = W
      canvas.height = H

      // Background
      const grad = getEffectiveBg()
      if (grad) {
        let canvasGrad: CanvasGradient
        if (grad.dir === 'lr') {
          canvasGrad = ctx.createLinearGradient(0, 0, W, 0)
        } else if (grad.dir === 'tlbr') {
          canvasGrad = ctx.createLinearGradient(0, 0, W, H)
        } else {
          canvasGrad = ctx.createLinearGradient(0, 0, 0, H)
        }
        canvasGrad.addColorStop(0, grad.from)
        canvasGrad.addColorStop(1, grad.to)
        ctx.fillStyle = canvasGrad
      } else {
        ctx.fillStyle = bgColor
      }
      ctx.fillRect(0, 0, W, H)

      // Caption text
      const caption = captions[lang] || ''
      if (caption) {
        ctx.fillStyle = textColor
        // 900 weight + W/13 size for strong visual impact (vs previous bold/W*18)
        ctx.font = `900 ${Math.round(W / 13)}px -apple-system, sans-serif`
        ctx.textAlign = 'center'
        ctx.textBaseline = 'middle'
        const lineH = Math.round(W / 14)
        const maxW = W - PADDING * 2
        const lines: string[] = []
        const isCJK = (ch: string) => /[\u3000-\u9fff\uac00-\ud7af\uff00-\uffef]/.test(ch)
        const segments = caption.split(' ')
        let line = ''
        for (const seg of segments) {
          if (Array.from(seg).some(isCJK)) {
            if (line) { lines.push(line); line = '' }
            for (const ch of Array.from(seg)) {
              const test = line + ch
              if (ctx.measureText(test).width > maxW) {
                if (line) lines.push(line)
                line = ch
              } else {
                line = test
              }
            }
          } else {
            const test = line ? `${line} ${seg}` : seg
            if (ctx.measureText(test).width > maxW) {
              if (line) lines.push(line)
              line = seg
            } else {
              line = test
            }
          }
        }
        if (line) lines.push(line)
        const totalH = lines.length * lineH
        const startY = CAPTION_H / 2 - totalH / 2 + lineH / 2
        lines.forEach((l, i) => ctx.fillText(l, W / 2, startY + i * lineH))
      }

      // Device frame — Y position based on alignment
      const fx = PADDING
      const fw = device.width
      const fh = device.height
      const fy = imageAlign === 'bottom'
        ? H - PADDING - fh  // bottom-align: device sits near the bottom
        : CAPTION_H         // center (current): device starts right after caption
      const r = device.cornerRadius

      ctx.fillStyle = device.frameColor
      ctx.beginPath()
      ctx.moveTo(fx + r, fy)
      ctx.lineTo(fx + fw - r, fy)
      ctx.arcTo(fx + fw, fy, fx + fw, fy + r, r)
      ctx.lineTo(fx + fw, fy + fh - r)
      ctx.arcTo(fx + fw, fy + fh, fx + fw - r, fy + fh, r)
      ctx.lineTo(fx + r, fy + fh)
      ctx.arcTo(fx, fy + fh, fx, fy + fh - r, r)
      ctx.lineTo(fx, fy + r)
      ctx.arcTo(fx, fy, fx + r, fy, r)
      ctx.closePath()
      ctx.fill()

      // Screen area (inset)
      const inset = 8
      const ir = Math.max(r - inset, 4)
      const sx = fx + inset
      const sy = fy + inset
      const sw = fw - inset * 2
      const sh = fh - inset * 2

      ctx.save()
      ctx.beginPath()
      ctx.moveTo(sx + ir, sy)
      ctx.lineTo(sx + sw - ir, sy)
      ctx.arcTo(sx + sw, sy, sx + sw, sy + ir, ir)
      ctx.lineTo(sx + sw, sy + sh - ir)
      ctx.arcTo(sx + sw, sy + sh, sx + sw - ir, sy + sh, ir)
      ctx.lineTo(sx + ir, sy + sh)
      ctx.arcTo(sx, sy + sh, sx, sy + sh - ir, ir)
      ctx.lineTo(sx, sy + ir)
      ctx.arcTo(sx, sy, sx + ir, sy, ir)
      ctx.closePath()
      ctx.clip()

      if (imageURL) {
        const img = new Image()
        img.onload = () => {
          const scale = Math.max(sw / img.width, sh / img.height)
          const dw = img.width * scale
          const dh = img.height * scale
          ctx.drawImage(img, sx + (sw - dw) / 2, sy + (sh - dh) / 2, dw, dh)
          ctx.restore()
          resolve()
        }
        img.onerror = () => {
          ctx.fillStyle = '#e5e7eb'
          ctx.fillRect(sx, sy, sw, sh)
          ctx.restore()
          resolve()
        }
        img.src = imageURL
      } else {
        ctx.fillStyle = '#e5e7eb'
        ctx.fillRect(sx, sy, sw, sh)
        ctx.restore()
        resolve()
      }
    })
  }, [bgColor, bgType, getEffectiveBg, textColor, captions, device, imageURL, imageAlign])

  const downloadSingle = async (lang: string) => {
    const canvas = canvasRef.current
    if (!canvas) return
    await drawFrame(canvas, lang)
    const a = document.createElement('a')
    a.download = `screenshot_${lang}_${device.id}.png`
    a.href = canvas.toDataURL('image/png')
    a.click()
  }

  const downloadAll = async () => {
    // Free plan: single language only
    if (!isPro) {
      await downloadSingle(selectedLang)
      return
    }

    setIsGenerating(true)

    try {
      if (appId && imageFile) {
        try {
          const activeCaptions = Object.fromEntries(
            Object.entries(captions).filter(([, v]) => v.trim() !== '')
          )
          const effectiveBg = getEffectiveBg()
          const result = await generateScreenshots(appId, {
            image: imageFile,
            device: device.id as 'iphone67' | 'iphone65' | 'ipad',
            bgColor: bgColor,
            bgGradientFrom: effectiveBg?.from,
            bgGradientTo: effectiveBg?.to,
            bgGradientDir: effectiveBg?.dir,
            textColor: textColor,
            captions: activeCaptions,
            imageAlign: imageAlign,
          })
          // Pro: bundle all languages into a ZIP
          const zip = new JSZip()
          for (const [lang, dataURL] of Object.entries(result.images)) {
            const base64 = dataURL.split(',')[1]
            zip.file(`screenshot_${lang}_${device.id}.png`, base64, { base64: true })
          }
          const blob = await zip.generateAsync({ type: 'blob' })
          const a = document.createElement('a')
          a.download = `screenshots_${device.id}.zip`
          a.href = URL.createObjectURL(blob)
          a.click()
          URL.revokeObjectURL(a.href)
          return
        } catch (err) {
          console.error('Server-side generation failed, falling back to canvas:', err)
        }
      }

      // Canvas fallback (Pro): ZIP from canvas renders
      const canvas = canvasRef.current
      if (!canvas) return
      const zip = new JSZip()
      for (const lang of LANGUAGES) {
        if (!captions[lang.code] && lang.code !== 'ja') continue
        await drawFrame(canvas, lang.code)
        const base64 = canvas.toDataURL('image/png').split(',')[1]
        zip.file(`screenshot_${lang.code}_${device.id}.png`, base64, { base64: true })
      }
      const blob = await zip.generateAsync({ type: 'blob' })
      const a = document.createElement('a')
      a.download = `screenshots_${device.id}.zip`
      a.href = URL.createObjectURL(blob)
      a.click()
      URL.revokeObjectURL(a.href)
    } finally {
      setIsGenerating(false)
    }
  }

  const previewLangs = LANGUAGES.filter(l => captions[l.code] || l.code === selectedLang)

  // Build preview background style
  const previewBg = (() => {
    const g = getEffectiveBg()
    if (!g) return { backgroundColor: bgColor }
    const dir = g.dir === 'lr' ? 'to right' : g.dir === 'tlbr' ? 'to bottom right' : 'to bottom'
    return { background: `linear-gradient(${dir}, ${g.from}, ${g.to})` }
  })()

  return (
    <div className="space-y-4">
      {/* Device + alignment */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">デバイス</label>
          <select
            value={device.id}
            onChange={(e) => setDevice(DEVICE_PRESETS.find(d => d.id === e.target.value)!)}
            className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {DEVICE_PRESETS.map(d => (
              <option key={d.id} value={d.id}>{d.label}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">画像位置</label>
          <div className="flex rounded-lg border overflow-hidden">
            {(['center', 'bottom'] as ImageAlign[]).map((align) => (
              <button
                key={align}
                onClick={() => setImageAlign(align)}
                className={`flex-1 py-2 text-sm transition-colors ${
                  imageAlign === align
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-600 hover:bg-gray-50'
                }`}
              >
                {align === 'center' ? '中央' : '下寄せ'}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Background */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">背景</label>

        {/* Type selector */}
        <div className="flex rounded-lg border overflow-hidden mb-3">
          {([
            { id: 'solid',    label: '単色' },
            { id: 'gradient', label: 'グラデーション' },
            { id: 'preset',   label: 'プリセット' },
          ] as { id: BgType; label: string }[]).map(({ id, label }) => (
            <button
              key={id}
              onClick={() => setBgType(id)}
              className={`flex-1 py-2 text-sm transition-colors ${
                bgType === id
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-600 hover:bg-gray-50'
              }`}
            >
              {label}
            </button>
          ))}
        </div>

        {/* Solid */}
        {bgType === 'solid' && (
          <div className="flex gap-3 items-center">
            <div className="flex items-center gap-2">
              <input type="color" value={bgColor} onChange={e => setBgColor(e.target.value)}
                className="h-9 w-14 border rounded cursor-pointer" />
              <span className="text-xs text-gray-500">背景色</span>
            </div>
            <div className="flex items-center gap-2">
              <input type="color" value={textColor} onChange={e => setTextColor(e.target.value)}
                className="h-9 w-14 border rounded cursor-pointer" />
              <span className="text-xs text-gray-500">テキスト色</span>
            </div>
          </div>
        )}

        {/* Gradient */}
        {bgType === 'gradient' && (
          <div className="space-y-3">
            <div className="flex gap-3 items-center">
              <div className="flex items-center gap-2">
                <input type="color" value={bgGradFrom} onChange={e => setBgGradFrom(e.target.value)}
                  className="h-9 w-14 border rounded cursor-pointer" />
                <span className="text-xs text-gray-500">開始色</span>
              </div>
              <span className="text-gray-400">→</span>
              <div className="flex items-center gap-2">
                <input type="color" value={bgGradTo} onChange={e => setBgGradTo(e.target.value)}
                  className="h-9 w-14 border rounded cursor-pointer" />
                <span className="text-xs text-gray-500">終了色</span>
              </div>
              <div className="flex items-center gap-2">
                <input type="color" value={textColor} onChange={e => setTextColor(e.target.value)}
                  className="h-9 w-14 border rounded cursor-pointer" />
                <span className="text-xs text-gray-500">テキスト色</span>
              </div>
            </div>
            <div className="flex gap-2">
              {([
                { id: 'tb', label: '↓ 縦' },
                { id: 'lr', label: '→ 横' },
                { id: 'tlbr', label: '↘ 斜め' },
              ] as { id: GradDir; label: string }[]).map(({ id, label }) => (
                <button
                  key={id}
                  onClick={() => setBgGradDir(id)}
                  className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
                    bgGradDir === id
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Preset */}
        {bgType === 'preset' && (
          <div className="space-y-3">
            <div className="grid grid-cols-8 gap-2">
              {BG_PRESETS.map(p => (
                <button
                  key={p.id}
                  onClick={() => setBgPresetId(p.id)}
                  title={p.label}
                  className={`relative h-10 rounded-lg overflow-hidden transition-all ${
                    bgPresetId === p.id ? 'ring-2 ring-blue-600 ring-offset-1 scale-105' : 'hover:scale-105'
                  }`}
                  style={{ background: `linear-gradient(to bottom, ${p.from}, ${p.to})` }}
                >
                  {bgPresetId === p.id && (
                    <span className="absolute inset-0 flex items-center justify-center text-white text-xs">✓</span>
                  )}
                </button>
              ))}
            </div>
            <div className="flex items-center gap-2">
              <input type="color" value={textColor} onChange={e => setTextColor(e.target.value)}
                className="h-9 w-14 border rounded cursor-pointer" />
              <span className="text-xs text-gray-500">テキスト色</span>
            </div>
          </div>
        )}
      </div>

      {/* Image upload */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">スクリーンショット画像</label>
        <div
          onClick={() => fileInputRef.current?.click()}
          className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center cursor-pointer hover:border-blue-400 transition-colors"
        >
          {imageURL ? (
            <img src={imageURL} alt="preview" className="max-h-32 mx-auto object-contain" />
          ) : (
            <div className="text-gray-400">
              <p className="text-sm">クリックして画像を選択</p>
              <p className="text-xs mt-1">PNG / JPG</p>
            </div>
          )}
        </div>
        <input ref={fileInputRef} type="file" accept="image/*" className="hidden" onChange={handleFileChange} />
      </div>

      {/* Language captions */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">キャプション（言語ごと）</label>
        <div className="flex flex-wrap gap-1 mb-3">
          {LANGUAGES.map(l => (
            <button
              key={l.code}
              onClick={() => setSelectedLang(l.code)}
              className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
                selectedLang === l.code
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {l.label}
            </button>
          ))}
        </div>
        <textarea
          value={captions[selectedLang] || ''}
          onChange={(e) => setCaptions(prev => ({ ...prev, [selectedLang]: e.target.value }))}
          placeholder={`${LANGUAGES.find(l => l.code === selectedLang)?.label} のキャプションを入力`}
          rows={2}
          className="w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Preview */}
      <div>
        <div className="flex items-center justify-between mb-2">
          <label className="block text-sm font-medium text-gray-700">プレビュー</label>
          <span className="text-xs text-gray-500">※ 実際のダウンロード画像は高解像度です</span>
        </div>
        <div
          className="rounded-xl p-4 flex flex-col items-center"
          style={{ ...previewBg, minHeight: 200 }}
        >
          {captions[selectedLang] && (
            <p className="text-center mb-3 text-base" style={{ color: textColor, fontWeight: 900 }}>
              {captions[selectedLang]}
            </p>
          )}
          <div
            className={`relative overflow-hidden ${imageAlign === 'bottom' ? 'mt-auto' : ''}`}
            style={{
              width: Math.min(device.width * 0.3, 180),
              height: Math.min(device.height * 0.3, 380),
              backgroundColor: device.frameColor,
              borderRadius: device.cornerRadius * 0.3,
              padding: 3,
            }}
          >
            <div className="w-full h-full overflow-hidden" style={{ borderRadius: (device.cornerRadius - 4) * 0.3 }}>
              {imageURL ? (
                <img src={imageURL} alt="screenshot" className="w-full h-full object-cover" />
              ) : (
                <div className="w-full h-full bg-gray-200 flex items-center justify-center">
                  <span className="text-gray-400 text-xs">画像なし</span>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Download buttons */}
      <div className="flex flex-wrap gap-2 items-center">
        {/* Individual language DL — available to all plans */}
        {previewLangs.map(l => (
          <button
            key={l.code}
            onClick={() => downloadSingle(l.code)}
            className="px-3 py-1.5 bg-white border border-gray-300 rounded-lg text-sm hover:bg-gray-50 flex items-center gap-1"
          >
            ↓ {l.label}
          </button>
        ))}

        {/* ZIP bulk DL — Pro only */}
        {isPro ? (
          <button
            onClick={downloadAll}
            disabled={isGenerating}
            className="ml-auto px-4 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50 flex items-center gap-1"
          >
            {isGenerating ? '生成中...' : '↓ ZIP一括DL'}
          </button>
        ) : (
          <button
            onClick={() => upgradeModal.open('多言語スクリーンショット一括生成')}
            className="ml-auto flex items-center gap-1.5 px-4 py-1.5 bg-gray-100 border border-gray-200 text-gray-400 rounded-lg text-sm hover:bg-gray-200 transition-colors"
          >
            🔒 ZIP一括DL
            <span className="text-blue-500 text-xs">Proで解除 →</span>
          </button>
        )}
      </div>

      {/* Hidden canvas for rendering */}
      <canvas ref={canvasRef} className="hidden" />

      <UpgradeModal
        isOpen={upgradeModal.isOpen}
        onClose={upgradeModal.close}
        triggerFeature={upgradeModal.triggerFeature}
      />
    </div>
  )
}
