# ASO Compass — Redesign

macOS ネイティブアプリ + Web 版の UI リデザイン提案（2案比較）。
Astro（tryastro.app）の美学を基調にした、A. Aurora（ダーク）/ B. Daylight（ライト）の2方向。

## 開く

`index.html` をブラウザで直接開くだけで動きます。ビルド不要。

## 構成

```
index.html              ルート。design_canvas に4アートボードを配置
src/
  data.js               モックデータ（アプリ・キーワード・順位履歴・ギャップ・メタデータ）
  themes.js             AURORA_THEME / DAYLIGHT_THEME のトークン定義
  aso-shared.jsx        共通パーツ（チャート/キーワード一覧/Segmented/EmptyState など）
  aso-app.jsx           macOS ネイティブ風アプリ本体（AsoApp）
  aso-web-app.jsx       Web ブラウザ風アプリ本体（AsoWebApp）
  starters/
    design-canvas.jsx   2案を並べるためのパン&ズーム可能なキャンバス
```

## 4つのアートボード

1. **macOS · Aurora** — 1280×800、ダーク + ミント／ラベンダーのオーロラ
2. **macOS · Daylight** — 1280×800、ウォームオフホワイト + コーラル
3. **Web · Aurora** — 1440×900、同テーマでブラウザクローム + トップナビ
4. **Web · Daylight** — 1440×900、同上

## 動作

全アートボードでインタラクティブ:
- アプリ切替（サイドバー / Web はドロップダウン）
- 3タブ切替（キーワード / 競合ギャップ / メタデータ）
- キーワード選択 → ランキング推移チャート更新（7/14/30日切替）
- ギャップ表のソート＆「すべて／圏外のみ」フィルタ
- メタデータ：言語＋タグ編集／文字数カウンタ（超過で赤）／差分検知
- 「キーワード追加」モーダル

## 編集のヒント

- **配色／質感** を変えるなら `src/themes.js`
- **モックデータ** を増やしたいなら `src/data.js`
- **コンポーネントの中身** を直すなら `src/aso-app.jsx` / `src/aso-web-app.jsx`
- **共通部品** は `src/aso-shared.jsx`

## 既存リポジトリとの関係

参照元: [entaku0818/ASOCompass](https://github.com/entaku0818/ASOCompass)
- macOS 版は `macos/ASOCompass/ASO Compass/Views/*.swift` の構造をベースに再設計
- Web 版は `frontend/src/app/apps/[id]/page.tsx` の構造を踏襲
