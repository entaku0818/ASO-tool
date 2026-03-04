'use client'

import { useState, useRef, useCallback } from 'react'

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

type Captions = Record<string, string>

interface ScreenshotGeneratorProps {
  appName?: string
}

export function ScreenshotGenerator({ appName }: ScreenshotGeneratorProps) {
  const [imageURL, setImageURL] = useState<string>('')
  const [captions, setCaptions] = useState<Captions>({ ja: '', en: '', 'zh-Hans': '', ko: '', fr: '', de: '' })
  const [selectedLang, setSelectedLang] = useState('ja')
  const [bgColor, setBgColor] = useState('#4F46E5')
  const [textColor, setTextColor] = useState('#FFFFFF')
  const [device, setDevice] = useState(DEVICE_PRESETS[0])
  const [isGenerating, setIsGenerating] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const canvasRef = useRef<HTMLCanvasElement>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = (ev) => setImageURL(ev.target?.result as string)
    reader.readAsDataURL(file)
  }

  const drawFrame = useCallback((canvas: HTMLCanvasElement, lang: string): Promise<void> => {
    return new Promise((resolve) => {
      const ctx = canvas.getContext('2d')
      if (!ctx) return resolve()

      const PADDING = 40
      const CAPTION_H = 120
      const W = device.width + PADDING * 2
      const H = device.height + PADDING * 2 + CAPTION_H

      canvas.width = W
      canvas.height = H

      // Background
      ctx.fillStyle = bgColor
      ctx.fillRect(0, 0, W, H)

      // Caption text
      const caption = captions[lang] || ''
      if (caption) {
        ctx.fillStyle = textColor
        ctx.font = `bold ${Math.round(W / 18)}px -apple-system, sans-serif`
        ctx.textAlign = 'center'
        ctx.textBaseline = 'middle'
        // Word wrap
        const words = caption.split(' ')
        const lineH = Math.round(W / 14)
        const maxW = W - PADDING * 2
        let line = ''
        const lines: string[] = []
        for (const word of words) {
          const test = line ? `${line} ${word}` : word
          if (ctx.measureText(test).width > maxW) {
            if (line) lines.push(line)
            line = word
          } else {
            line = test
          }
        }
        if (line) lines.push(line)
        const totalH = lines.length * lineH
        const startY = CAPTION_H / 2 - totalH / 2 + lineH / 2
        lines.forEach((l, i) => ctx.fillText(l, W / 2, startY + i * lineH))
      }

      // Device frame
      const fx = PADDING
      const fy = CAPTION_H
      const fw = device.width
      const fh = device.height
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
          // Cover fit
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
  }, [bgColor, textColor, captions, device, imageURL])

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
    setIsGenerating(true)
    const canvas = canvasRef.current
    if (!canvas) { setIsGenerating(false); return }
    for (const lang of LANGUAGES) {
      if (!captions[lang.code] && lang.code !== 'ja') continue
      await drawFrame(canvas, lang.code)
      const a = document.createElement('a')
      a.download = `screenshot_${lang.code}_${device.id}.png`
      a.href = canvas.toDataURL('image/png')
      a.click()
      await new Promise(r => setTimeout(r, 300))
    }
    setIsGenerating(false)
  }

  const previewLangs = LANGUAGES.filter(l => captions[l.code] || l.code === selectedLang)

  return (
    <div className="space-y-4">
      {/* Device + image settings */}
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
          <label className="block text-sm font-medium text-gray-700 mb-1">背景色</label>
          <div className="flex gap-2">
            <input type="color" value={bgColor} onChange={e => setBgColor(e.target.value)}
              className="h-10 w-16 border rounded cursor-pointer" />
            <input type="color" value={textColor} onChange={e => setTextColor(e.target.value)}
              className="h-10 w-16 border rounded cursor-pointer" title="テキスト色" />
            <span className="text-xs text-gray-500 self-center">背景 / テキスト</span>
          </div>
        </div>
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
          style={{ backgroundColor: bgColor, minHeight: 200 }}
        >
          {captions[selectedLang] && (
            <p className="font-bold text-center mb-3 text-sm" style={{ color: textColor }}>
              {captions[selectedLang]}
            </p>
          )}
          <div
            className="relative overflow-hidden"
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
      <div className="flex flex-wrap gap-2">
        {previewLangs.map(l => (
          <button
            key={l.code}
            onClick={() => downloadSingle(l.code)}
            className="px-3 py-1.5 bg-white border border-gray-300 rounded-lg text-sm hover:bg-gray-50 flex items-center gap-1"
          >
            ↓ {l.label}
          </button>
        ))}
        <button
          onClick={downloadAll}
          disabled={isGenerating}
          className="px-4 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50 ml-auto"
        >
          {isGenerating ? '生成中...' : '全言語を一括ダウンロード'}
        </button>
      </div>

      {/* Hidden canvas for rendering */}
      <canvas ref={canvasRef} className="hidden" />
    </div>
  )
}
