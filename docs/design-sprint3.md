# スプリント3 UX 設計ドキュメント

作成日: 2026-04-04
担当: デザイナー
ステータス: エンジニアレビュー待ち

---

## 1. 成果シェアカード UI（最優先）

### 1-1. 追加箇所の特定

対象ファイル: `frontend/src/app/apps/[id]/page.tsx`
対象コンポーネント: `RisingKeywordsSection`（lines 132–178）

現状、急上昇キーワードテーブルの各行に「上昇幅」列はあるが、シェア導線がない。

```tsx
// 現状の上昇幅セル（line 168-172）
<td className="py-3 px-4 text-right">
  <span className="inline-flex items-center gap-1 text-green-600 font-bold">
    ▲ {kw.improvement}
  </span>
</td>
// ← ここにシェアボタンを追加
```

---

### 1-2. テーブル行へのシェアボタン追加

```
[RisingKeywordsSection テーブル 改善後]

┌────────────────────────────────────────────────────────────────┐
│ キーワード │ 国   │ 現在順位 │ 以前の順位 │ 上昇幅       │      │  ← thead
│──────────────────────────────────────────────────────────────│
│ パズルゲーム│ JP  │ 3位      │ 47位       │ ▲ 44  [ 🎉 ] │      │  ← シェアボタン追加
│ RPG攻略    │ JP  │ 12位     │ 23位       │ ▲ 11  [ 🎉 ] │      │
│ 無料ゲーム  │ US  │ 8位      │ 21位       │ ▲ 13  [ 🎉 ] │      │
└────────────────────────────────────────────────────────────────┘

シェアボタン [ 🎉 ] スタイル:
  bg-gradient-to-r from-orange-400 to-pink-500
  text-white text-xs font-semibold
  px-2.5 py-1 rounded-full
  hover: opacity-90 scale-105 transition-all
  ※ animate-bounce は初回マウント時 1 回のみ（useEffect + setTimeout で除去）
```

---

### 1-3. シェアモーダル

ボタンクリック後、シェアモーダルを表示:

```
┌─────────────────────────────────────────────────────┐
│  シェアする                                    [ ✕ ] │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │  [シェアカードプレビュー 400×210px]           │   │
│  │  （下記シェアカードビジュアル参照）            │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
│  シェアテキスト（編集可）:                          │
│  ┌─────────────────────────────────────────────┐   │
│  │ 「パズルゲーム」が #47 → #3 に上昇しました 🎉│   │  ← textarea (rows=3)
│  │ +44位アップ！ #ASO #AppStore               │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
│  ┌────────────────────┐  ┌───────────────────┐     │
│  │  🐦 X でシェア     │  │  📋 画像をコピー   │     │
│  └────────────────────┘  └───────────────────┘     │
│  ┌────────────────────┐                             │
│  │  🔗 URLをコピー    │                             │
│  └────────────────────┘                             │
│                                                     │
└─────────────────────────────────────────────────────┘

モーダルスタイル:
  bg-white rounded-2xl shadow-2xl p-6 max-w-lg w-full mx-4
  オーバーレイ: bg-black/50

ボタンスタイル:
  🐦 X:    bg-black text-white px-4 py-2 rounded-lg text-sm font-medium
  📋 コピー: bg-gray-100 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium
  🔗 URL:  bg-gray-100 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium
  コピー成功後: テキストが「✓ コピーしました」に 2 秒間変化
```

---

### 1-4. シェアカードビジュアル仕様（canvas 生成）

```
カードサイズ: 1200 × 630px（OGP 標準）
プレビュー表示: 400 × 210px（scale 1/3）

┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  [背景: bg-gradient-to-br from-indigo-600 via-purple-700       │
│         to-pink-600 （既存 BG_PRESETS の配色を流用）]           │
│                                                                 │
│  🎉                                      ← 絵文字 80px, 左上    │
│                                                                 │
│  キーワード順位が上昇しました！            ← 38px, white, bold  │
│                                                                 │
│  ─────────────────────────────────────    ← rgba(255,255,255,0.3)│
│                                                                 │
│  「パズルゲーム」                          ← 24px, white/80%   │
│                                                                 │
│      #47          →         #3            ← 数字エリア          │
│   (54px, white/50%)      (72px, white,    │
│                           font-black)     │
│                                           │
│              +44位 ▲                      ← 28px, #86EFAC (green-300)│
│                                                                 │
│  ─────────────────────────────────────                          │
│                                                                 │
│  App Store Optimization · powered by ASO Tool  ← 16px white/40%│
│                                                 右下配置         │
└─────────────────────────────────────────────────────────────────┘

描画順序（canvas 2D）:
  1. グラデーション背景 fillRect
  2. セパレーター線 2本 (rgba stroke)
  3. 絵文字・ヘッドライン・キーワード名 fillText
  4. 前順位（薄い）・矢印・新順位（強調）fillText
  5. 上昇幅 fillText (green-300)
  6. フッターテキスト fillText
  7. canvas.toDataURL('image/png') → <img> プレビューに反映
```

---

### 1-5. シェアテキスト（X / Twitter 用テンプレート）

```
「{keyword}」が #{previousRank} → #{currentRank} に上昇しました 🎉
+{improvement}位アップ！App Store 最適化の成果が出ています 📈

#ASO #AppStoreOptimization #アプリ開発
```

---

### 1-6. エンジニア実装メモ

```tsx
// RisingKeyword 型（既存）
interface RisingKeyword {
  keyword_id: string
  keyword: string
  country: string
  current_rank: number
  previous_rank: number
  improvement: number
}

// 追加するステート（RisingKeywordsSection 内）
const [shareTarget, setShareTarget] = useState<RisingKeyword | null>(null)

// canvas 生成は ScreenshotGenerator の drawFrame パターンを流用
// navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]) で画像コピー
// X シェア URL: https://twitter.com/intent/tweet?text=encodeURIComponent(text)
```

---

## 2. キーワード制限 到達時のUI

### 2-1. 現状の問題

`apps/[id]/page.tsx` line 599-601:
```tsx
if (err instanceof Error && err.message.includes('402')) {
  upgradeModal.open('キーワード追加')
}
```

**問題点**: 上限到達まで何も表示されない。ユーザーが 10 件目を登録しようとして初めて UpgradeModal が唐突に出る。  
**理想**: 近づいたら事前警告 → 達したらインライン通知 → UpgradeModal への誘導。

---

### 2-2. 段階的 UI 設計

```
段階1: 0〜7件（通常）
┌──────────────────────────────────────────────────┐
│ キーワードを追加                     7 / 10       │  ← カウンター表示のみ（gray）
│ ┌─────────────────────┐ [+ 追加]                 │
│ │ キーワードを入力...   │                         │
│ └─────────────────────┘                         │
└──────────────────────────────────────────────────┘

段階2: 8〜9件（警告）
┌──────────────────────────────────────────────────┐
│ キーワードを追加               ⚠️ 残り1件         │  ← amber カウンター
│ ┌─────────────────────┐ [+ 追加]                 │
│ │ キーワードを入力...   │                         │
│ └─────────────────────┘                         │
│                                                  │
│ ┌─ インライン警告バナー ──────────────────────┐  │
│ │ ⚠️  キーワード上限まであと{n}件です。         │  │  ← amber-50 bg
│ │    Pro なら無制限に追加できます。             │  │
│ │    [ Pro にアップグレード →]                 │  │  ← amber border
│ └─────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘

段階3: 10件（上限到達）
┌──────────────────────────────────────────────────┐
│ キーワードを追加            🔒 10 / 10 （上限）   │  ← red カウンター
│ ┌─────────────────────┐ [+ 追加 disabled]        │  ← ボタン disabled
│ │ キーワードを入力...   │                         │
│ └─────────────────────┘                         │
│                                                  │
│ ┌─ 上限到達バナー ────────────────────────────┐  │
│ │ 🔒  キーワードの上限（10件）に達しました。   │  │  ← red-50 bg
│ │    Pro プランなら無制限に追加できます。       │  │
│ │                                              │  │
│ │  [  Pro を 7日間 無料で試す  →  ]           │  │  ← red border + blue CTA
│ └─────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘
```

---

### 2-3. カウンター表示の仕様

```
表示位置: キーワード入力欄のラベル右端（flex justify-between）

0-7件:    text-gray-400, "N / 10"
8件:      text-amber-500, "⚠️ 残り2件"
9件:      text-amber-600, "⚠️ 残り1件"
10件:     text-red-600 font-semibold, "🔒 10 / 10（上限）"

Pro の場合: "N 件"（上限表示なし）
```

---

### 2-4. バナースタイル仕様

#### 警告バナー（8〜9件）
```
bg-amber-50 border border-amber-200 rounded-xl p-4
text-sm text-amber-800
CTAリンク: text-blue-600 underline font-medium → UpgradeModal open
```

#### 上限到達バナー（10件）
```
bg-red-50 border border-red-200 rounded-xl p-4
text-sm text-red-800
CTAボタン: w-full py-2 bg-blue-600 text-white rounded-lg text-sm font-semibold
           → UpgradeModal open（triggerFeature="キーワード追加（上限到達）"）
```

---

### 2-5. UpgradeModal との連携タイミング

| タイミング | 動作 |
|-----------|------|
| 8件目登録後 | インライン警告バナーを表示（UpgradeModal は開かない） |
| 9件目登録後 | インライン警告バナーを更新（残り1件） |
| 10件目登録試行 | API 402 → UpgradeModal open（既存動作を維持） |
| 上限到達バナーの CTA クリック | UpgradeModal open（新規） |

**原則**: UpgradeModal の多発を避ける。バナーの CTA で十分と判断した場合は Modal なしでも可。  
エンジニアと相談の上、10件時の挙動（自動 Modal vs バナー CTA 経由 Modal）を決める。

---

### 2-6. エンジニア実装メモ

```tsx
// keywords は useKeywords(appId) から取得済み（page.tsx line 580）
const FREE_KEYWORD_LIMIT = 10
const keywordCount = keywords.length
const isAtLimit = !isPro && keywordCount >= FREE_KEYWORD_LIMIT
const isNearLimit = !isPro && keywordCount >= FREE_KEYWORD_LIMIT - 2  // 8件以上

// 入力フォーム直下に条件付きレンダリング
{isAtLimit && <KeywordLimitBanner onUpgrade={() => upgradeModal.open('キーワード追加（上限到達）')} />}
{!isAtLimit && isNearLimit && <KeywordNearLimitBanner remaining={FREE_KEYWORD_LIMIT - keywordCount} onUpgrade={...} />}
```

---

## 3. BUG-12: テンプレートページの認証 UX

### 3-1. 現状の問題

`frontend/src/app/templates/page.tsx`:
- `useAuth()` は `{ user, isLoading }` を返す
- **`isLoading` が `true` の間**（JWT 検証中）、`isPro = user?.is_pro ?? false` が `false` になる
  → 全テンプレートがロック表示されてしまう
- **未認証ユーザー**がアクセスした場合、テンプレートが「ロック状態」で見えてしまう
  → ログインページへのリダイレクトがない

---

### 3-2. 認証フロー設計

```
[ページアクセス時の状態遷移]

     ┌──────────────────────────────────────────────────────────┐
     │  useAuth().isLoading === true                            │
     │                                                          │
     │  [ローディングスケルトン表示]                             │  ← 状態A
     └──────────────────────────────────────────────────────────┘
                         │ isLoading が false になる
            ┌────────────┴────────────┐
            ▼                         ▼
   user !== null                 user === null
   （認証済み）                   （未認証）
            │                         │
            ▼                         ▼
  [通常のテンプレート表示]    [リダイレクト or 未認証バナー]
                                       │
              ┌────────────────────────┤
              ▼                        ▼
     router.push('/login')    未認証バナー表示（ソフト）
     （ハード誘導）            ← こちらを推奨（BUG-12の方針次第）
```

---

### 3-3. ローディングスケルトン UI

認証確認中（`isLoading === true`）に表示するスケルトン:

```
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  テンプレートライブラリ                                        │
│  App Store スクリーンショット用テンプレートを選んで始めよう     │
│                                                              │
│  ┌── フィルタバー（実際のボタンを表示）─────────────────────┐ │
│  │  [すべて]  [ゲーム]  [ビジネス]  [教育]                  │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │  ← スケルトン
│  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │
│  │ ░░░░░       │  │ ░░░░░       │  │ ░░░░░       │
│  │ ░░░░░░░░░   │  │ ░░░░░░░░░   │  │ ░░░░░░░░░   │
│  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │  │ ░░░░░░░░░░░ │
│  │ ...         │  │ ...         │  │ ...         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                              │
└──────────────────────────────────────────────────────────────┘

スケルトンスタイル:
  bg-gray-200 animate-pulse rounded-xl
  サムネイル部: aspect-video rounded-t-xl
  バッジ部: h-4 w-12 rounded-full
  タイトル部: h-4 w-32 rounded
  説明文部: h-3 w-full rounded
  ボタン部: h-9 w-full rounded-lg
```

---

### 3-4. 未認証状態の表示（ソフトリダイレクト推奨）

```
┌──────────────────────────────────────────────────────────────┐
│  テンプレートライブラリ                                        │
│                                                              │
│  ┌─ 未認証バナー ──────────────────────────────────────────┐ │
│  │                                                          │ │
│  │  🔒  テンプレートの利用にはログインが必要です             │ │
│  │                                                          │ │
│  │  アカウントを作成してテンプレートを無料で使いましょう     │ │
│  │                                                          │ │
│  │  [  ログイン / 新規登録  →  ]                           │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                              │
│  [テンプレートプレビュー（ぼかし状態）がグリッドで見える]      │  ← 訴求用にコンテンツは見せる
│  ░░░░░░░░░░   ░░░░░░░░░░   ░░░░░░░░░░                       │
│  blur-md + pointer-events-none                               │
│                                                              │
└──────────────────────────────────────────────────────────────┘

バナースタイル:
  bg-blue-50 border border-blue-200 rounded-xl p-6 mb-6 text-center
  テキスト: text-gray-700
  CTA: bg-blue-600 text-white px-6 py-2 rounded-lg font-semibold
       → router.push('/login')
```

**ソフトリダイレクトを推奨する理由**: テンプレートのビジュアルをぼかして見せることで、ログイン後の価値を予感させられる。ハードリダイレクト（即 /login）は UX が荒い。

---

### 3-5. エンジニア実装メモ

```tsx
// templates/page.tsx に追加
const { user, isLoading: authLoading } = useAuth()
const isPro = user?.is_pro ?? false

// 認証ロード中
if (authLoading) return <TemplateSkeletonGrid />  // スケルトンコンポーネント

// 未認証
if (!user) return (
  <div>
    <UnauthBanner />
    <div className="blur-md pointer-events-none">
      <TemplateGrid templates={dummyTemplates} />  // or実データをそのまま表示
    </div>
  </div>
)

// 認証済み → 既存の表示ロジック（変更不要）
```

**注意**: `authLoading` の確認を `isLoading`（テンプレートAPI）と混同しないこと。  
`authLoading` が先に解決してから `getTemplates()` を呼ぶ順序にする。

---

## エンジニア連携サマリー（スプリント3）

| # | 機能 | 実装難易度 | バックエンド変更 | 優先度 |
|---|------|----------|---------------|--------|
| 1 | シェアカード（canvas生成 + モーダル） | M | 不要（フロントのみ） | 最優先 |
| 2 | シェアボタン（RisingKeywords行追加） | XS | 不要 | 最優先 |
| 3 | キーワードカウンター表示 | XS | 不要 | High |
| 4 | キーワード上限インライン警告バナー | S | 不要 | High |
| 5 | テンプレートページ スケルトン + 未認証バナー | S | 不要 | Low（BUG-12対応） |

**全タスクがフロントエンドのみで完結**。バックエンド変更不要なため、エンジニアはフロント作業に集中できる。

---

*このドキュメントはデザイナーが作成。実装前にエンジニア・PDM のレビューを受けること。*
