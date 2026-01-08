# ASO-tool

App Store Optimization ツール - キーワードランキング追跡、レビュー監視、競合分析

## 技術構成

| 層 | 技術 | デプロイ先 |
|---|---|---|
| Frontend | Next.js | Firebase Hosting |
| API | Go (chi) | Cloud Run |
| DB | PostgreSQL | Cloud SQL |
| 定期実行 | - | Cloud Scheduler |

## クイックスタート

```bash
# ローカル環境起動
docker-compose up -d

# Backend
cd backend && go run cmd/api/main.go

# Frontend
cd frontend && npm run dev
```

## ディレクトリ構成

```
.
├── backend/          # Go API サーバー (詳細: backend/README.md)
├── frontend/         # Next.js フロントエンド (詳細: frontend/README.md)
├── docs/             # デプロイ設定ガイド
│   ├── gcp-setup.md      # GCP (Cloud Run, Cloud SQL)
│   └── firebase-setup.md # Firebase Hosting
└── scripts/          # 開発用スクリプト
```

## ドキュメント

- [開発計画](./docs/development-plan.md) - フェーズ別開発ロードマップ
- [Backend README](./backend/README.md) - API 開発・起動方法
- [Frontend README](./frontend/README.md) - UI 開発・起動方法
- [GCP セットアップ](./docs/gcp-setup.md) - Cloud Run, Cloud SQL 設定
- [Firebase セットアップ](./docs/firebase-setup.md) - Firebase Hosting 設定
