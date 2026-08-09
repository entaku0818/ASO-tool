'use client'

import { useState, useEffect } from 'react'
import { createApp, createKeyword, fetchAppInfo, App, CreateAppRequest } from '@/lib/api'

// ── localStorage helpers ──────────────────────────────────────
const LS_COMPLETED  = 'onboarding_completed'
const LS_STEP       = 'onboarding_step'
const LS_SKIP_UNTIL = 'onboarding_skip_until'

export function isOnboardingCompleted(): boolean {
  try { return localStorage.getItem(LS_COMPLETED) === 'true' } catch { return false }
}

function markCompleted() {
  try { localStorage.setItem(LS_COMPLETED, 'true'); localStorage.removeItem(LS_STEP) } catch {}
}

function saveStep(step: number) {
  try { localStorage.setItem(LS_STEP, String(step)) } catch {}
}

function savedStep(): number {
  try { return Number(localStorage.getItem(LS_STEP) ?? '1') } catch { return 1 }
}

function skipFor24h() {
  try { localStorage.setItem(LS_SKIP_UNTIL, new Date(Date.now() + 86_400_000).toISOString()) } catch {}
}

function isSkipped(): boolean {
  try {
    const v = localStorage.getItem(LS_SKIP_UNTIL)
    return !!v && new Date(v) > new Date()
  } catch { return false }
}

// ── should show onboarding? ───────────────────────────────────
export function shouldShowOnboarding(userCreatedAt: string, appsCount: number): boolean {
  if (isOnboardingCompleted()) return false
  if (isSkipped()) return false
  const within72h = Date.now() - new Date(userCreatedAt).getTime() < 72 * 3_600_000
  return within72h && appsCount === 0
}

// ── keyword chip presets ──────────────────────────────────────
const KEYWORD_CHIPS = [
  'ゲーム', 'パズル', '無料', 'RPG', '子供向け',
  '英語学習', '家計簿', 'カレンダー', 'メモ', '写真編集',
]

// ── Progress bar ──────────────────────────────────────────────
function ProgressDots({ step, total }: { step: number; total: number }) {
  return (
    <div className="flex items-center gap-2 mb-6">
      {Array.from({ length: total }).map((_, i) => (
        <div
          key={i}
          className={`h-2 rounded-full transition-all ${
            i < step ? 'bg-blue-600 w-8' : 'bg-gray-200 dark:bg-gray-700 w-2'
          }`}
        />
      ))}
      <span className="ml-2 text-xs text-gray-400 dark:text-gray-500">STEP {step} / {total}</span>
    </div>
  )
}

// ── Props ─────────────────────────────────────────────────────
interface OnboardingFlowProps {
  userCreatedAt: string
  appsCount: number
  onComplete: (newApp: App) => void
  onDismiss: () => void
  onUpgradeClick?: () => void
}

// ── Main component ────────────────────────────────────────────
export function OnboardingFlow({ userCreatedAt, appsCount, onComplete, onDismiss, forceShow, onUpgradeClick }: OnboardingFlowProps & { forceShow?: boolean }) {
  const [step, setStep] = useState<1 | 2 | 3 | 'done'>(() => {
    const s = savedStep()
    return (s === 2 || s === 3 ? s : 1) as 1 | 2
  })
  const [visible, setVisible] = useState(forceShow || shouldShowOnboarding(userCreatedAt, appsCount))

  // Step 2 state
  const [appName, setAppName]     = useState('')
  const [bundleId, setBundleId]   = useState('')
  const [platform, setPlatform]   = useState<'ios' | 'android'>('ios')
  const [storeUrl, setStoreUrl]   = useState('')
  const [urlInput, setUrlInput]   = useState('')
  const [isFetchingInfo, setIsFetchingInfo] = useState(false)
  const [fetchInfoError, setFetchInfoError] = useState('')
  const [isCreatingApp, setIsCreatingApp] = useState(false)
  const [appError, setAppError]   = useState('')
  const [createdApp, setCreatedApp] = useState<App | null>(null)

  // Auto-extract bundle ID from store URL
  const handleUrlInput = async (url: string) => {
    setUrlInput(url)
    setFetchInfoError('')
    if (!url.trim()) return

    // Extract App Store numeric ID: /id1234567890
    const iosMatch = url.match(/\/id(\d+)/)
    // Extract Google Play package: ?id=com.example.app
    const androidMatch = url.match(/[?&]id=([a-zA-Z0-9_.]+)/)

    const detectedPlatform = iosMatch ? 'ios' : androidMatch ? 'android' : null
    const extractedId = iosMatch?.[1] ?? androidMatch?.[1] ?? null

    if (!extractedId || !detectedPlatform) return

    setIsFetchingInfo(true)
    try {
      const info = await fetchAppInfo(extractedId, detectedPlatform)
      setAppName(info.name || appName)
      setBundleId(info.bundle_id || extractedId)
      setStoreUrl(url)
      setPlatform(detectedPlatform as 'ios' | 'android')
    } catch {
      setFetchInfoError('アプリ情報の取得に失敗しました。手動で入力してください。')
    } finally {
      setIsFetchingInfo(false)
    }
  }

  // Step 3 state
  const [selectedChips, setSelectedChips] = useState<string[]>([])
  const [customKeyword, setCustomKeyword] = useState('')
  const [pendingKeywords, setPendingKeywords] = useState<string[]>([])
  const [isAddingKeywords, setIsAddingKeywords] = useState(false)

  useEffect(() => {
    if (step !== 1) saveStep(typeof step === 'number' ? step : 3)
  }, [step])

  if (!visible) return null

  // ── Step 1: Welcome ────────────────────────────────────────
  const handleStart = () => setStep(2)

  const handleSkip = () => {
    skipFor24h()
    setVisible(false)
    onDismiss()
  }

  // ── Step 2: App registration ───────────────────────────────
  const handleCreateApp = async () => {
    if (!appName.trim() || !bundleId.trim()) return
    setIsCreatingApp(true)
    setAppError('')
    try {
      const req: CreateAppRequest = {
        name: appName.trim(),
        bundle_id: bundleId.trim(),
        platform,
        store_url: storeUrl.trim() || undefined,
      }
      const app = await createApp(req)
      setCreatedApp(app)
      setStep(3)
    } catch (e) {
      const is402 = e instanceof Error && e.message.includes('402')
      setAppError(is402
        ? 'Freeプランでは1アプリまで登録できます。アップグレードをご検討ください。'
        : 'アプリの登録に失敗しました。入力内容を確認してください。'
      )
    } finally {
      setIsCreatingApp(false)
    }
  }

  // ── Step 3: Keywords ───────────────────────────────────────
  const allKeywords = Array.from(new Set([...selectedChips, ...pendingKeywords]))

  const toggleChip = (kw: string) => {
    setSelectedChips(prev =>
      prev.includes(kw) ? prev.filter(k => k !== kw) : [...prev, kw]
    )
  }

  const addCustom = () => {
    const kw = customKeyword.trim()
    if (!kw || allKeywords.includes(kw)) return
    setPendingKeywords(prev => [...prev, kw])
    setCustomKeyword('')
  }

  const removeKeyword = (kw: string) => {
    setSelectedChips(prev => prev.filter(k => k !== kw))
    setPendingKeywords(prev => prev.filter(k => k !== kw))
  }

  const handleFinish = async () => {
    if (!createdApp || allKeywords.length < 3) return
    setIsAddingKeywords(true)
    try {
      await Promise.all(allKeywords.map(kw => createKeyword(createdApp.id, kw)))
      markCompleted()
      setStep('done')
    } catch {
      // partial failure is acceptable — mark complete anyway
      markCompleted()
      setStep('done')
    } finally {
      setIsAddingKeywords(false)
    }
  }

  // ── Overlay wrapper ────────────────────────────────────────
  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-[#12161e] rounded-2xl shadow-2xl dark:shadow-black/20 w-full max-w-md overflow-hidden">
        {/* Gradient bar */}
        <div className="h-1.5 bg-gradient-to-r from-blue-600 to-purple-600" />

        <div className="p-8">
          {/* ── Step 1: Welcome ──────────────────────────── */}
          {step === 1 && (
            <div className="text-center">
              <div className="text-5xl mb-4">🎯</div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">ASO Compass へようこそ！</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
                あなたのアプリのキーワード順位を自動で追跡・分析します。
                3ステップでセットアップが完了します。
              </p>
              <div className="flex items-center justify-center gap-2 mb-8 text-xs text-gray-400 dark:text-gray-500">
                <span className="flex items-center gap-1">
                  <span className="w-5 h-5 rounded-full bg-blue-600 text-white flex items-center justify-center text-xs">1</span>
                  アプリ登録
                </span>
                <span className="text-gray-300 dark:text-gray-600">→</span>
                <span className="flex items-center gap-1">
                  <span className="w-5 h-5 rounded-full bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 flex items-center justify-center text-xs">2</span>
                  キーワード追加
                </span>
                <span className="text-gray-300 dark:text-gray-600">→</span>
                <span className="flex items-center gap-1">
                  <span className="w-5 h-5 rounded-full bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 flex items-center justify-center text-xs">3</span>
                  完了
                </span>
              </div>
              <button
                onClick={handleStart}
                className="w-full py-3 bg-blue-600 dark:bg-blue-500 text-white rounded-xl font-semibold hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors mb-3"
              >
                アプリを登録して始める →
              </button>
              <button
                onClick={handleSkip}
                className="w-full py-2 text-sm text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 transition-colors"
              >
                スキップして後で設定する
              </button>
            </div>
          )}

          {/* ── Step 2: App registration ──────────────────── */}
          {step === 2 && (
            <div>
              <ProgressDots step={1} total={3} />
              <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-1">アプリを登録する</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">ストアのURLを貼り付けると自動で情報を取得します</p>

              {/* URL auto-fill */}
              <div className="mb-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl p-3">
                <label className="block text-xs font-medium text-blue-700 dark:text-blue-400 mb-1">
                  ストアURL（App Store / Google Play）
                </label>
                <div className="flex gap-2">
                  <input
                    type="url"
                    value={urlInput}
                    onChange={e => handleUrlInput(e.target.value)}
                    placeholder="https://apps.apple.com/jp/app/.../id..."
                    className="flex-1 px-3 py-2 text-sm border border-blue-200 dark:border-blue-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
                  />
                  {isFetchingInfo && <div className="w-5 h-5 border-2 border-blue-500 border-t-transparent rounded-full animate-spin self-center flex-shrink-0" />}
                </div>
                {fetchInfoError && <p className="text-xs text-red-500 mt-1">{fetchInfoError}</p>}
                {bundleId && !fetchInfoError && <p className="text-xs text-blue-600 dark:text-blue-400 mt-1">✓ アプリ情報を取得しました</p>}
              </div>

              <div className="space-y-3">
                <div>
                  <label className="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">アプリ名 *</label>
                  <input
                    type="text"
                    value={appName}
                    onChange={e => setAppName(e.target.value)}
                    placeholder="My App"
                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">Bundle ID *</label>
                  <input
                    type="text"
                    value={bundleId}
                    onChange={e => setBundleId(e.target.value)}
                    placeholder="com.example.myapp"
                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">プラットフォーム</label>
                  <select
                    value={platform}
                    onChange={e => setPlatform(e.target.value as 'ios' | 'android')}
                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                  >
                    <option value="ios">iOS</option>
                    <option value="android">Android</option>
                  </select>
                </div>
              </div>

              {appError && <p className="mt-3 text-sm text-red-500">{appError}</p>}

              <button
                onClick={handleCreateApp}
                disabled={isCreatingApp || !appName.trim() || !bundleId.trim()}
                className="w-full mt-5 py-3 bg-blue-600 dark:bg-blue-500 text-white rounded-xl font-semibold hover:bg-blue-700 dark:hover:bg-blue-600 disabled:opacity-50 transition-colors"
              >
                {isCreatingApp ? '登録中...' : '次へ →'}
              </button>
            </div>
          )}

          {/* ── Step 3: Keywords ──────────────────────────── */}
          {step === 3 && (
            <div>
              <ProgressDots step={2} total={3} />
              <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-1">キーワードを追加する</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
                よく使われるキーワードから選ぶか、直接入力してください（3つ以上推奨）
              </p>

              {/* Chips */}
              <div className="flex flex-wrap gap-2 mb-4">
                {KEYWORD_CHIPS.map(kw => (
                  <button
                    key={kw}
                    onClick={() => toggleChip(kw)}
                    className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
                      selectedChips.includes(kw)
                        ? 'bg-blue-600 dark:bg-blue-500 text-white'
                        : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
                    }`}
                  >
                    {kw}
                  </button>
                ))}
              </div>

              {/* Custom input */}
              <div className="flex gap-2 mb-4">
                <input
                  type="text"
                  value={customKeyword}
                  onChange={e => setCustomKeyword(e.target.value)}
                  onKeyDown={e => { if (e.key === 'Enter') { e.preventDefault(); addCustom() } }}
                  placeholder="キーワードを入力..."
                  className="flex-1 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500"
                />
                <button
                  onClick={addCustom}
                  disabled={!customKeyword.trim()}
                  className="px-4 py-2 text-sm bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 disabled:opacity-40 transition-colors"
                >
                  追加
                </button>
              </div>

              {/* Selected */}
              {allKeywords.length > 0 && (
                <div className="mb-4">
                  <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                    選択中: <span className={allKeywords.length >= 3 ? 'text-blue-600 dark:text-blue-400 font-semibold' : 'text-orange-500'}>
                      {allKeywords.length}/3 以上
                    </span>
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {allKeywords.map(kw => (
                      <span
                        key={kw}
                        className="flex items-center gap-1 px-3 py-1 bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 rounded-full text-sm"
                      >
                        {kw}
                        <button
                          onClick={() => removeKeyword(kw)}
                          className="ml-0.5 text-blue-400 dark:text-blue-500 hover:text-blue-700 dark:hover:text-blue-300 leading-none"
                          aria-label={`${kw}を削除`}
                        >
                          ×
                        </button>
                      </span>
                    ))}
                  </div>
                </div>
              )}

              <button
                onClick={handleFinish}
                disabled={allKeywords.length < 3 || isAddingKeywords}
                className="w-full py-3 bg-blue-600 dark:bg-blue-500 text-white rounded-xl font-semibold hover:bg-blue-700 dark:hover:bg-blue-600 disabled:opacity-50 transition-colors"
              >
                {isAddingKeywords ? '登録中...' : '完了して順位を確認 →'}
              </button>
              {allKeywords.length < 3 && (
                <p className="mt-2 text-center text-xs text-gray-400 dark:text-gray-500">
                  あと {3 - allKeywords.length} 件追加してください
                </p>
              )}
            </div>
          )}

          {/* ── Completion ───────────────────────────────── */}
          {step === 'done' && (
            <div className="text-center">
              <div className="text-5xl mb-4">✅</div>
              <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">セットアップ完了！</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-5">
                キーワード順位の追跡を開始しました。<br />
                最初のデータは明日の朝に反映されます。
              </p>

              {onUpgradeClick && (
                <div className="bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 border border-blue-100 dark:border-blue-800 rounded-xl p-4 mb-4 text-left">
                  <p className="text-xs font-bold text-blue-700 dark:text-blue-300 uppercase tracking-wide mb-2">✦ Pro でできること</p>
                  <ul className="space-y-1.5 mb-4">
                    {[
                      '競合アプリのキーワードを逆引き',
                      'アプリ・キーワード数が無制限',
                      'CSVエクスポートで分析を効率化',
                    ].map((f) => (
                      <li key={f} className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                        <span className="text-green-500 flex-shrink-0">✓</span>
                        {f}
                      </li>
                    ))}
                  </ul>
                  <button
                    onClick={onUpgradeClick}
                    className="w-full py-2.5 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg text-sm font-semibold hover:from-blue-700 hover:to-purple-700 transition-all"
                  >
                    Proにアップグレード — ¥1,650/月〜
                  </button>
                </div>
              )}

              <button
                onClick={() => { setVisible(false); if (createdApp) onComplete(createdApp) }}
                className={`w-full py-3 rounded-xl font-semibold transition-colors ${onUpgradeClick ? 'border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 text-sm' : 'bg-blue-600 dark:bg-blue-500 text-white hover:bg-blue-700 dark:hover:bg-blue-600'}`}
              >
                ダッシュボードへ →
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
