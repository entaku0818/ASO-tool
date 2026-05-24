# ASO Tool — システム構成

## 全体像

```
┌─────────────────────────────────────────────────────────────┐
│                        ユーザー                              │
│           ブラウザ                   macOS アプリ            │
└──────────────┬───────────────────────────┬──────────────────┘
               │ HTTPS                     │ HTTPS
               ▼                           │
┌──────────────────────────┐               │
│  Cloud Run               │               │
│  aso-frontend (Next.js)  │               │
└──────────────┬───────────┘               │
               │ REST API                  │
               ▼                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Cloud Run                                                    │
│  aso-api (Go)                                                 │
│                                                               │
│  /api/auth/*      認証・JWT                                   │
│  /api/apps/*      アプリ管理・キーワード順位                   │
│  /api/billing/*   Stripe 課金                                 │
│  /api/public/*    認証なし公開 API                            │
└────────────────────────┬─────────────────────────────────────┘
                         │ pgx/v5
                         ▼
              ┌─────────────────────┐
              │  Cloud SQL          │
              │  PostgreSQL (aso-db)│
              └─────────────────────┘
```

## 外部サービス連携

```
aso-api ──→ 🍎 App Store Connect API    アナリティクスデータ取得
aso-api ──→ 🍎 Apple Search Ads API     キーワード人気スコア取得
aso-api ──→ 💳 Stripe                   決済・Webhook
aso-api ──→ 📧 Resend                   メール送信（購入完了等）

aso-frontend ──→ 🍎 App Store サジェスト    キーワード検索候補
aso-frontend ──→ 🤖 Google Play スクレイピング
aso-frontend ──→ ⚙️  db-manager (Cloud Run) Cloud SQL 起動
```

## 日次バッチ（毎朝 9:00 JST）

```
GitHub Actions ──→ Cloud Run Job (aso-batch)
                        │
                        ├── rankings        アプリのキーワード順位を更新
                        ├── tracked-keywords トラッキングキーワード順位を更新
                        ├── store-rankings  App Store / Google Play ランキング取得
                        └── keyword-cache   Apple Search Ads から人気スコアを取得・保存
                                            ※ ADMIN_ASA_* シークレット設定が必要
```

## CI/CD

```
git push → Backend CI  → ✅ → Deploy Backend  → Cloud Run (aso-api)
        → Frontend CI → ✅ → Deploy Frontend → Cloud Run (aso-frontend)
        → (手動)        →     Batch Scheduler → Cloud Run Job (aso-batch)
```

## 課金フロー

```
1. ユーザー → /api/billing/checkout
2. → Stripe Checkout ページへリダイレクト
3. 決済完了 → Stripe Webhook → /api/stripe/webhook
4. plan を pro に更新 + Resend でメール送信
```

---

## URL 一覧

| 名前 | URL |
|---|---|
| フロントエンド | https://aso-tool.entaku.app |
| バックエンド API | https://aso-api-671942133800.asia-northeast1.run.app |
| DB マネージャー | https://db-manager-te5er5txcq-an.a.run.app |
