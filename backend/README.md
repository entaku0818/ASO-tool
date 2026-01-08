# Backend

Go (chi) で構築した REST API サーバー

## ディレクトリ構成

```
backend/
├── cmd/
│   ├── api/          # API サーバー
│   └── batch/        # バッチジョブ
├── internal/
│   ├── handler/      # HTTP ハンドラー
│   ├── service/      # ビジネスロジック
│   ├── repository/   # DB 操作
│   ├── scraper/      # データ取得
│   ├── model/        # データモデル
│   └── middleware/   # ミドルウェア
└── migrations/       # DB マイグレーション
```

## 開発コマンド

```bash
# API サーバー起動
make run

# テスト実行
make test

# ビルド
make build

# Docker イメージ作成
make docker-build
```

## 環境変数

| 変数名 | 説明 |
|--------|------|
| `DATABASE_URL` | PostgreSQL 接続文字列 |
| `PORT` | サーバーポート (default: 8080) |

## API エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/apps` | アプリ一覧取得 |
| POST | `/api/apps` | アプリ登録 |
| GET | `/api/apps/:id` | アプリ詳細取得 |
| GET | `/api/apps/:id/keywords` | キーワード一覧 |
| GET | `/api/apps/:id/rankings` | ランキング履歴 |
| GET | `/api/apps/:id/reviews` | レビュー一覧 |
