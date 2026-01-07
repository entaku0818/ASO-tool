# Frontend

Next.js + Tailwind CSS で構築したダッシュボード UI

## ディレクトリ構成

```
frontend/
├── src/
│   ├── app/          # Next.js App Router
│   ├── components/   # UI コンポーネント
│   ├── contexts/     # React Context
│   ├── hooks/        # カスタムフック
│   └── lib/          # API クライアント
├── public/           # 静的ファイル
└── __tests__/        # テスト
```

## 開発コマンド

```bash
# 開発サーバー起動
npm run dev

# ビルド
npm run build

# テスト実行
npm test

# 型チェック
npm run type-check
```

## 環境変数

| 変数名 | 説明 |
|--------|------|
| `NEXT_PUBLIC_API_URL` | Backend API の URL |

## 技術スタック

- **Next.js 14** - React フレームワーク
- **Tailwind CSS** - スタイリング
- **Recharts** - グラフ表示
- **Jest** - テスト
