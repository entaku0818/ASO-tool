# スプリント4 UX 設計ドキュメント

作成日: 2026-04-05
担当: デザイナー
ステータス: エンジニアレビュー待ち

---

## 1. 成果シェア LP デザイン（最優先）

### 1-1. ページ概要・URL設計

| 項目 | 内容 |
|------|------|
| URL | `/share` （クエリパラメータ渡し） |
| 例 | `/share?kw=パズルゲーム&from=47&to=3&app=My%20Game&country=JP` |
| 認証 | 不要（public page） |
| Header | 公開ヘッダー（`publicPages` に `/share` を追加） |
| OGP | `og:image` に canvas 生成カード画像を設定（Next.js Metadata API） |

**クエリパラメータ一覧**:
```
kw      : キーワード名（必須）
from    : 以前の順位（必須）
to      : 現在の順位（必須）
app     : アプリ名（任意）
country : 国コード JP/US 等（任意、デフォルト JP）
```

---

### 1-2. モバイルファースト レイアウト全体図

```
[モバイル / max-w-lg mx-auto px-4]

┌─────────────────────────────────────┐  ← min-h-screen
│                                     │
│         ASO Tool                    │  ← ロゴ (上部, py-6)
│                                     │
│  ┌─────────────────────────────┐   │
│  │  [シェアカードビジュアル]    │   │  ← ①カードビジュアル
│  │   aspect-[1200/630]         │   │     w-full rounded-2xl shadow-xl
│  │   グラデーション背景         │   │
│  │   (canvas or CSS再現)       │   │
│  └─────────────────────────────┘   │
│                                     │
│  🎉 {キーワード}が                  │  ← ②ヘッドライン
│     {n}位上昇しました！             │     text-2xl font-black
│                                     │
│  {アプリ名} · {国}                  │  ← ③サブテキスト text-gray-500 text-sm
│                                     │
│  ─────────────────────────────     │  ← divider
│                                     │
│  ASO Tool で同じ成果を              │  ← ④バリュープロポジション
│  あなたのアプリにも。               │     text-lg text-gray-700
│                                     │
│  ✓ キーワード順位を自動追跡         │  ← ⑤3つの価値訴求（チェックリスト）
│  ✓ 急上昇キーワードを即キャッチ    │
│  ✓ App Store 最適化を継続改善      │
│                                     │
│  ┌─────────────────────────────┐   │
│  │   🚀  無料で始める           │   │  ← ⑥CTA（Primary）
│  └─────────────────────────────┘   │     bg-gradient-to-r from-blue-600
│                                     │     to-purple-600, w-full, py-4
│  ログイン済みの方は→ ダッシュボード │  ← ⑦セカンダリリンク
│                                     │
│  ─────────────────────────────     │
│                                     │
│  [SNSシェアボタン]                  │  ← ⑧シェア再拡散
│  [ 🐦 X でシェア ]  [ 📋 コピー ]   │
│                                     │
│  © ASO Tool  · プライバシーポリシー │  ← フッター text-xs text-gray-400
└─────────────────────────────────────┘


[デスクトップ / max-w-2xl mx-auto]
┌────────────────────────────────────────────────────────┐
│  ASO Tool                                              │
├─────────────────────────────┬──────────────────────────┤
│  [カードビジュアル]         │  🎉 {kw}が{n}位上昇！     │
│  aspect-[1200/630]         │                            │
│  left column               │  {アプリ名} · {国}         │
│  (flex-1)                  │  ──────────────────       │
│                            │  ASO Tool で同じ成果を     │
│                            │  あなたのアプリにも。       │
│                            │                            │
│                            │  ✓ キーワード順位を追跡    │
│                            │  ✓ 急上昇をキャッチ        │
│                            │  ✓ ASO を継続改善          │
│                            │                            │
│                            │  [🚀 無料で始める]          │
│                            │  ログイン済み → Dashboard  │
│                            │                            │
│                            │  [🐦 X] [📋 コピー]        │
└─────────────────────────────┴──────────────────────────┘
```

---

### 1-3. カードビジュアル仕様（CSS 再現版）

シェア LP ではサーバーサイドで canvas を生成できないため、**CSS で OGP カードを再現**する。

```
カードコンテナ:
  className: "relative w-full aspect-[1200/630] rounded-2xl overflow-hidden shadow-xl"
  style: { background: 'linear-gradient(135deg, #4F46E5, #7C3AED, #DB2777)' }

内部レイアウト (absolute inset-0 flex flex-col justify-between p-8 sm:p-12):

  [上段]
    🎉 (text-5xl sm:text-6xl)

  [中段: flex-1 flex flex-col justify-center gap-3]
    「{keyword}」
      text-white/70 text-lg sm:text-2xl font-medium

    [順位変化エリア: flex items-center gap-4 sm:gap-8]
      #{from}          →          #{to}
      (white/40,       (white/50   (white,
       text-4xl        text-2xl)   text-6xl sm:text-7xl,
       font-bold)                  font-black)

    +{improvement}位 ▲
      text-green-300 text-xl sm:text-2xl font-semibold

  [下段]
    App Store Optimization · powered by ASO Tool
      text-white/40 text-xs sm:text-sm text-right
```

---

### 1-4. CTA ボタン設計

```
Primary CTA:
  テキスト: "🚀  無料で始める"
  スタイル: w-full py-4 text-lg font-bold rounded-2xl
            bg-gradient-to-r from-blue-600 to-purple-600
            text-white shadow-lg hover:shadow-xl
            active:scale-95 transition-all
  リンク先: /login （または登録フォームへ）
  ※ ユーザーが認証済みの場合: "/" へリダイレクト

Secondary:
  テキスト: "ログイン済みの方はダッシュボードへ →"
  スタイル: text-sm text-blue-600 hover:underline text-center mt-2
  リンク先: /
```

---

### 1-5. エラー・エッジケース

| 状況 | 表示 |
|------|------|
| `kw` パラメータ欠如 | 「シェアリンクが無効です」メッセージ + ASO Tool トップへの CTA |
| `from` <= `to`（下落） | カードは表示するが「上昇」テキストは出さない。CTAは通常通り |
| 数字が異常値（0以下など） | 「−」表示にフォールバック |

---

### 1-6. エンジニア実装メモ

```tsx
// 新規ページ: frontend/src/app/share/page.tsx
// publicPages に '/share' を Header.tsx に追加

// URL パース
const searchParams = useSearchParams()
const keyword   = searchParams.get('kw') ?? ''
const fromRank  = Number(searchParams.get('from')) || 0
const toRank    = Number(searchParams.get('to'))   || 0
const appName   = searchParams.get('app') ?? ''
const country   = searchParams.get('country') ?? 'JP'
const improvement = fromRank - toRank  // 正の値 = 上昇

// シェア URL 生成（RisingKeywordsSection 側で作成）
const shareUrl = `${window.location.origin}/share?kw=${encodeURIComponent(kw.keyword)}&from=${kw.previous_rank}&to=${kw.current_rank}&app=${encodeURIComponent(appName)}&country=${kw.country.toUpperCase()}`

// X シェア
const twitterUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(tweetText)}&url=${encodeURIComponent(shareUrl)}`
```

---

## 2. キーワード制限告知バナー

### 2-1. Sprint 3 設計からの変更点

Sprint 3 では「8件以上で警告」と設計したが、**Sprint 4 要件では「7件到達時」に前倒し**。  
マーケターの意見を受けてトリガーを早める（余裕を持って Pro 訴求するため）。

| 件数 | Sprint 3 設計 | Sprint 4 更新後 |
|-----|-------------|----------------|
| 0〜6件 | カウンターのみ | カウンターのみ（変更なし） |
| **7件** | カウンターのみ | ⚠️ 警告バナー（新規） |
| 8〜9件 | 警告バナー | 警告バナー継続（変更なし） |
| 10件 | ブロックバナー | ブロックバナー（変更なし） |

---

### 2-2. 7件到達時 警告バナー（Dismissible）

```
[キーワードリスト 上部 or 入力フォーム直下]

┌──────────────────────────────────────────────────────────────┐
│  ⚠️  キーワードが 7 / 10 件に達しました                  [ ✕ ] │  ← dismissible
│                                                              │
│  無料プランの上限まであと 3 件。Pro なら無制限に追跡できます。 │
│                                                              │
│     [ Pro を 7日間 無料で試す →]                            │  ← inline CTA
└──────────────────────────────────────────────────────────────┘

スタイル:
  bg-amber-50 border border-amber-200 rounded-xl px-4 py-3
  flex items-start justify-between gap-3

⚠️アイコン: text-amber-500 flex-shrink-0 mt-0.5

テキスト:
  title: text-sm font-semibold text-amber-800
  body:  text-xs text-amber-700 mt-0.5

[ ✕ ] 閉じるボタン:
  text-amber-400 hover:text-amber-600
  クリック後: localStorage.setItem('kw_limit_warning_dismissed', '1')
  再表示条件: 10件到達で再表示（dismissible を上書き）

CTA リンク:
  text-xs font-semibold text-blue-600 hover:underline mt-1
  → upgradeModal.open('キーワード追加（上限7件）')
```

---

### 2-3. 10件到達時 ブロックバナー（Undismissible）

```
[キーワード入力フォーム 直下・常時固定表示]

┌──────────────────────────────────────────────────────────────┐
│  🔒  キーワード上限（10件）に達しました                       │  ← 閉じるボタンなし
│                                                              │
│  新しいキーワードを追加するには Pro プランが必要です。         │
│  キーワードの削除、または Pro へのアップグレードをどうぞ。    │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │       🚀  Pro を 7日間 無料で試す                       │  │  ← Primary CTA
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘

スタイル:
  bg-red-50 border-2 border-red-200 rounded-xl p-4

🔒アイコン: text-red-500

テキスト:
  title: text-sm font-semibold text-red-800
  body:  text-xs text-red-700 mt-1

CTA ボタン:
  w-full mt-3 py-2.5 text-sm font-bold
  bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl
  hover:opacity-90 transition-opacity
  → upgradeModal.open('キーワード追加（上限到達）')

入力フォームの状態:
  textarea: disabled opacity-50 cursor-not-allowed
  追加ボタン: disabled bg-gray-100 text-gray-400
```

---

### 2-4. カウンター表示の仕様（更新）

```
キーワード入力ラベル行:
  flex items-center justify-between mb-1

  左: <label>キーワードを追加</label>
  右: <KeywordCounter count={keywords.length} isPro={isPro} />

カウンター表示:
  0〜6件:  text-gray-400 text-xs "N / 10"
  7〜9件:  text-amber-500 text-xs font-medium "⚠️ N / 10"
  10件:    text-red-600 text-xs font-semibold "🔒 10 / 10"
  Pro:     text-gray-400 text-xs "N 件"
```

---

### 2-5. 表示優先度・重複排除ルール

```
同一画面で複数のバナーが出ないよう制御:

優先度（高→低）:
  1. 10件ブロックバナー（常時固定）
  2. 7〜9件警告バナー（dismiss 可能）
  3. カウンターのみ（常時）

実装:
  {isAtLimit
    ? <KeywordBlockBanner onUpgrade={...} />
    : isNearLimit && !isDismissed
    ? <KeywordWarningBanner remaining={...} onDismiss={...} onUpgrade={...} />
    : null
  }
```

---

## 3. AIキャプション一括生成 UI（Pro 専用）

### 3-1. 実装箇所

対象: `frontend/src/components/ScreenshotGenerator.tsx`
対象セクション: lines 534-558 「Language captions」ブロック

Sprint 2 設計（単一言語生成）の上位機能として追加。

---

### 3-2. UIレイアウト（全体）

```
[Language captions セクション]

┌──────────────────────────────────────────────────────┐
│ キャプション（言語ごと）           7 / 10 件          │
│                                                      │
│ [日本語] [English] [中文] [한국어] [Français] [Deutsch]│
│                                                      │
│ ┌──────────────────────────────────────────────┐    │
│ │ 日本語のキャプションを入力...                 │    │
│ └──────────────────────────────────────────────┘    │
│                                                      │
│  [ ✨ この言語で生成 ]   [ ✦ 全言語まとめて生成 Pro ]│  ← 2ボタン並列
└──────────────────────────────────────────────────────┘

ボタン配置: flex justify-end gap-2

[ ✨ この言語で生成 ] （Sprint 2 設計済み、変更なし）:
  bg-purple-50 text-purple-700 border border-purple-200 rounded-lg px-3 py-1.5 text-sm

[ ✦ 全言語まとめて生成 Pro ]:
  Pro の場合:
    bg-gradient-to-r from-purple-600 to-blue-600
    text-white rounded-lg px-3 py-1.5 text-sm font-semibold
    hover: opacity-90

  Free の場合:
    bg-gray-100 text-gray-400 border border-gray-200 rounded-lg px-3 py-1.5 text-sm
    末尾: <span className="text-blue-500 text-xs">Proで解除 →</span>
    クリック: upgradeModal.open('AIキャプション一括生成')
```

---

### 3-3. 全言語一括生成 進捗UI

「全言語まとめて生成」クリック後のローディング状態:

```
[生成中 — 6言語プログレス表示]

┌──────────────────────────────────────────────────────┐
│ キャプション（言語ごと）                              │
│                                                      │
│ ✦ 全言語を生成中...                                  │  ← section header (purple)
│                                                      │
│ ┌────────────────────────────────────────────────┐   │
│ │  日本語    ✅ 完了                              │   │
│ │  English   ✅ 完了                              │   │
│ │  中文      ⏳ 生成中...   (spinner)            │   │  ← 現在処理中の言語
│ │  한국어    ○  待機中                           │   │
│ │  Français  ○  待機中                           │   │
│ │  Deutsch   ○  待機中                           │   │
│ └────────────────────────────────────────────────┘   │
│                                                      │
│  ──────────────────────────────────────────────     │
│  ░░░░░░░░░░░░░░░░░░░░░░░░                           │  ← プログレスバー
│  3 / 6 言語完了                                     │  ← カウンター text-sm text-purple-600
│                                                      │
│                              [ キャンセル ]           │
└──────────────────────────────────────────────────────┘

各言語行のアイコン:
  完了:   ✅  text-green-500
  生成中: spinner (animate-spin w-4 h-4 text-purple-500)
  待機中: ○   text-gray-300

プログレスバー:
  bg-gray-200 rounded-full h-1.5 (外枠)
  bg-gradient-to-r from-purple-500 to-blue-500 rounded-full h-1.5
  width: `${(completedCount / 6) * 100}%`  transition-all duration-300

キャンセルボタン:
  text-sm text-gray-400 hover:text-gray-600
  クリック: AbortController で API リクエストをキャンセル
```

---

### 3-4. 全言語生成 完了状態

```
[生成完了]

┌──────────────────────────────────────────────────────┐
│ キャプション（言語ごと）                              │
│                                                      │
│ ✅  6言語のキャプションを生成しました                 │  ← success banner (green-50)
│     各言語タブで内容を確認・編集できます              │     border-green-200 text-green-700
│                                                      │
│ [日本語✓] [English✓] [中文✓] [한국어✓] [Français✓] [Deutsch✓]│  ← タブに ✓ バッジ
│                                                      │
│ ┌──────────────────────────────────────────────┐    │
│ │ 直感的な操作でゲームを極めろ                   │    │  ← 選択中の言語の内容
│ └──────────────────────────────────────────────┘    │
│                                                      │
│             [ ✦ 再生成 ]   [ ✨ この言語のみ再生成 ] │
└──────────────────────────────────────────────────────┘

言語タブのバッジ（生成済みの場合）:
  通常タブ: bg-blue-600 text-white (選択中) / bg-gray-100 text-gray-600
  生成済み: タブ末尾に tiny ✓ (text-green-500 text-[10px])

再生成ボタン:
  [ ✦ 再生成 ]:      全言語を再生成 → プログレスUIに戻る
  [ ✨ この言語のみ再生成 ]: 現在タブの言語のみ再生成
```

---

### 3-5. エラー状態（一括生成）

```
[一部の言語で生成失敗]

┌──────────────────────────────────────────────────────┐
│ ⚠️  一部の言語で生成に失敗しました                    │  ← amber banner
│                                                      │
│ ┌────────────────────────────────────────────────┐   │
│ │  日本語    ✅ 完了                              │   │
│ │  English   ✅ 完了                              │   │
│ │  中文      ✗  失敗（再試行）                   │   │  ← retry リンク
│ │  한국어    ✅ 完了                              │   │
│ │  Français  ✗  失敗（再試行）                   │   │
│ │  Deutsch   ✅ 完了                              │   │
│ └────────────────────────────────────────────────┘   │
│                                                      │
│  [ 失敗した言語のみ再試行 ]  [ 完了した言語で続ける ] │
└──────────────────────────────────────────────────────┘

失敗行:
  ✗ アイコン: text-red-400
  「再試行」: text-xs text-blue-600 hover:underline (個別リトライ)
```

---

### 3-6. エンジニア実装メモ（一括生成）

```tsx
// 追加ステート
type LangStatus = 'idle' | 'loading' | 'done' | 'error'
const [bulkStatus, setBulkStatus] = useState<Record<string, LangStatus>>({})
const [isBulkGenerating, setIsBulkGenerating] = useState(false)

// 処理フロー（逐次 or 並列の選択はバックエンドと相談）
const handleBulkGenerate = async () => {
  setIsBulkGenerating(true)
  setBulkStatus(Object.fromEntries(LANGUAGES.map(l => [l.code, 'idle'])))

  for (const lang of LANGUAGES) {
    setBulkStatus(prev => ({ ...prev, [lang.code]: 'loading' }))
    try {
      const result = await generateCaption(appId, imageFile, lang.code)
      setCaptions(prev => ({ ...prev, [lang.code]: result.caption }))
      setBulkStatus(prev => ({ ...prev, [lang.code]: 'done' }))
    } catch {
      setBulkStatus(prev => ({ ...prev, [lang.code]: 'error' }))
    }
  }

  setIsBulkGenerating(false)
}

// Pro チェック
const handleBulkClick = () => {
  if (!isPro) {
    upgradeModal.open('AIキャプション一括生成')
    return
  }
  if (!imageFile) {
    // inline warning 表示（Sprint 2 設計と同様）
    return
  }
  handleBulkGenerate()
}
```

---

## エンジニア連携サマリー（スプリント4）

| # | 機能 | 実装難易度 | バックエンド変更 | 優先度 |
|---|------|----------|---------------|--------|
| 1 | シェア LP `/share` ページ新設 | S | 不要（フロントのみ） | 最優先 |
| 2 | Header に `/share` を public page 追加 | XS | 不要 | 最優先（LP とセット） |
| 3 | キーワード 7件 警告バナー（dismissible） | S | 不要 | High |
| 4 | キーワード 10件 ブロックバナー | XS | 不要（Sprint 3 設計の微修正） | High |
| 5 | AIキャプション 全言語一括生成ボタン | S | AI API エンドポイント必要（Sprint 2 と同API） | Medium |
| 6 | 全言語生成 進捗プログレスUI | S | 不要（フロントのみ） | Medium |

**注意**: #5 の AI API は Sprint 2 で設計済みの単一言語 API をループ呼び出しすることで実現可能。バックエンドに「全言語一括 API」を作らずフロントで逐次処理する方が実装コストが低い。

---

## デザイン整合メモ（スプリント横断）

| テーマ | Sprint 2 | Sprint 3 | Sprint 4（今回） |
|--------|---------|---------|----------------|
| シェアカード | canvas 生成設計 | テーブル行ボタン設計 | LP ページとして具現化 |
| キーワード上限 | — | 8件〜警告設計 | **7件〜**に前倒し修正 |
| AI一括生成 | 単一言語ボタン | — | 全言語一括 + 進捗UI |

---

*このドキュメントはデザイナーが作成。実装前にエンジニア・PDM のレビューを受けること。*
