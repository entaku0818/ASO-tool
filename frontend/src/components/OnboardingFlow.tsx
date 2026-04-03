'use client'

import { useState, useEffect } from 'react'
import { createApp, createKeyword, App, CreateAppRequest } from '@/lib/api'

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
            i < step ? 'bg-blue-600 w-8' : 'bg-gray-200 w-2'
          }`}
        />
      ))}
      <span className="ml-2 text-xs text-gray-400">STEP {step} / {total}</span>
    </div>
  )
}

// ── Props ─────────────────────────────────────────────────────
interface OnboardingFlowProps {
  userCreatedAt: string
  appsCount: number
  onComplete: (newApp: App) => void
  onDismiss: () => void
}

// ── Main component ────────────────────────────────────────────
export function OnboardingFlow({ userCreatedAt, appsCount, onComplete, onDismiss }: OnboardingFlowProps) {
  const [step, setStep] = useState<1 | 2 | 3 | 'done'>(() => {
    const s = savedStep()
    return (s === 2 || s === 3 ? s : 1) as 1 | 2
  })
  const [visible, setVisible] = useState(shouldShowOnboarding(userCreatedAt, appsCount))

  // Step 2 state
  const [appName, setAppName]     = useState('')
  const [bundleId, setBundleId]   = useState('')
  const [platform, setPlatform]   = useState<'ios' | 'android'>('ios')
  const [storeUrl, setStoreUrl]   = useState('')
  const [isCreatingApp, setIsCreatingApp] = useState(false)
  const [appError, setAppError]   = useState('')
  const [createdApp, setCreatedApp] = useState<App | null>(null)

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
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden">
        {/* Gradient bar */}
        <div className="h-1.5 bg-gradient-to-r from-blue-600 to-purple-600" />

        <div className="p-8">
          {/* ── Step 1: Welcome ──────────────────────────── */}
          {step === 1 && (
            <div className="text-center">
              <div className="text-5xl mb-4">🎯</div>
              <h2 className="text-xl font-bold text-gray-900 mb-2">ASO Tool へようこそ！</h2>
              <p className="text-sm text-gray-500 mb-6">
                あなたのアプリのキーワード順位を自動で追跡・分析します。
                3ステップでセットアップが完了します。
              </p>
              <div className="flex items-center justify-center gap-2 mb-8 text-xs text-gray-400">
                <span className="flex items-center gap-1">
                  <span className="w-5 h-5 rounded-full bg-blue-600 text-white flex items-center justify-center text-xs">1</span>
                  アプリ登録
                </span>
                <span className="text-gray-300">→</span>
                <span className="flex items-center gap-1">
                  <span className="w-5 h-5 rounded-full bg-gray-200 text-gray-500 flex items-center justify-center text-xs">2</span>
                  キーワード追加
                </span>
                <span className="text-gray-300">→</span>
                <span className="flex items-center gap-1">
                  <span className="w-5 h-5 rounded-full bg-gray-200 text-gray-500 flex items-center justify-center text-xs">3</span>
                  完了
                </span>
              </div>
              <button
                onClick={handleStart}
                className="w-full py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 transition-colors mb-3"
              >
                アプリを登録して始める →
              </button>
              <button
                onClick={handleSkip}
                className="w-full py-2 text-sm text-gray-400 hover:text-gray-600 transition-colors"
              >
                スキップして後で設定する
              </button>
            </div>
          )}

          {/* ── Step 2: App registration ──────────────────── */}
          {step === 2 && (
            <div>
              <ProgressDots step={1} total={3} />
              <h2 className="text-lg font-bold text-gray-900 mb-1">アプリを登録する</h2>
              <p className="text-sm text-gray-500 mb-5">アプリの情報を入力してください</p>

              <div className="space-y-3">
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">アプリ名 *</label>
                  <input
                    type="text"
                    value={appName}
                    onChange={e => setAppName(e.target.value)}
                    placeholder="My App"
                    className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">Bundle ID *</label>
                  <input
                    type="text"
                    value={bundleId}
                    onChange={e => setBundleId(e.target.value)}
                    placeholder="com.example.myapp"
                    className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">プラットフォーム</label>
                  <select
                    value={platform}
                    onChange={e => setPlatform(e.target.value as 'ios' | 'android')}
                    className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="ios">iOS</option>
                    <option value="android">Android</option>
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-600 mb-1">ストアURL（任意）</label>
                  <input
                    type="url"
                    value={storeUrl}
                    onChange={e => setStoreUrl(e.target.value)}
                    placeholder="https://apps.apple.com/jp/app/..."
                    className="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              {appError && <p className="mt-3 text-sm text-red-500">{appError}</p>}

              <button
                onClick={handleCreateApp}
                disabled={isCreatingApp || !appName.trim() || !bundleId.trim()}
                className="w-full mt-5 py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                {isCreatingApp ? '登録中...' : '次へ →'}
              </button>
            </div>
          )}

          {/* ── Step 3: Keywords ──────────────────────────── */}
          {step === 3 && (
            <div>
              <ProgressDots step={2} total={3} />
              <h2 className="text-lg font-bold text-gray-900 mb-1">キーワードを追加する</h2>
              <p className="text-sm text-gray-500 mb-4">
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
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
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
                  className="flex-1 px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button
                  onClick={addCustom}
                  disabled={!customKeyword.trim()}
                  className="px-4 py-2 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 disabled:opacity-40 transition-colors"
                >
                  追加
                </button>
              </div>

              {/* Selected */}
              {allKeywords.length > 0 && (
                <div className="mb-4">
                  <p className="text-xs text-gray-500 mb-2">
                    選択中: <span className={allKeywords.length >= 3 ? 'text-blue-600 font-semibold' : 'text-orange-500'}>
                      {allKeywords.length}/3 以上
                    </span>
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {allKeywords.map(kw => (
                      <span
                        key={kw}
                        className="flex items-center gap-1 px-3 py-1 bg-blue-50 text-blue-700 rounded-full text-sm"
                      >
                        {kw}
                        <button
                          onClick={() => removeKeyword(kw)}
                          className="ml-0.5 text-blue-400 hover:text-blue-700 leading-none"
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
                className="w-full py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                {isAddingKeywords ? '登録中...' : '完了して順位を確認 →'}
              </button>
              {allKeywords.length < 3 && (
                <p className="mt-2 text-center text-xs text-gray-400">
                  あと {3 - allKeywords.length} 件追加してください
                </p>
              )}
            </div>
          )}

          {/* ── Completion ───────────────────────────────── */}
          {step === 'done' && (
            <div className="text-center">
              <div className="text-5xl mb-4">✅</div>
              <h2 className="text-xl font-bold text-gray-900 mb-2">セットアップ完了！</h2>
              <p className="text-sm text-gray-500 mb-6">
                キーワード順位の追跡を開始しました。<br />
                最初のデータは明日の朝に反映されます。
              </p>

              <div className="bg-blue-50 rounded-xl p-4 mb-6 text-left">
                <p className="text-xs font-semibold text-blue-700 mb-1">── Pro ならもっとできます ──</p>
                <p className="text-sm text-gray-600">
                  競合アプリが使っているキーワードを逆引きして、上位表示のチャンスを見つけましょう。
                </p>
              </div>

              <button
                onClick={() => { setVisible(false); if (createdApp) onComplete(createdApp) }}
                className="w-full py-3 bg-blue-600 text-white rounded-xl font-semibold hover:bg-blue-700 transition-colors"
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
