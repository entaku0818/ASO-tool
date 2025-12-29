# ASO-tool

App Store Optimization ツール - キーワードランキング追跡、レビュー監視、競合分析

## 技術構成

| 層 | 技術 | デプロイ先 | 月額費用 |
|---|---|---|---|
| Frontend | Next.js | Firebase Hosting | $0 |
| API | Go (chi) | Cloud Run | $0 |
| DB | PostgreSQL | Cloud SQL | $7〜9 |
| 定期実行 | - | Cloud Scheduler | $0 |

**💰 月額合計: 約 $8〜12**

---

## 機能一覧

1. **キーワードランキング追跡** - 指定キーワードでの順位を毎日記録
2. **レビュー監視** - 新着レビューを毎時チェック・通知
3. **競合分析** - 競合アプリのメタデータ履歴を保存
4. **ダッシュボード** - 推移をグラフで可視化

---

## 開発計画

### Phase 1: Backend基盤 (API + DB)

```
backend/
├── cmd/api/main.go          # エントリーポイント
├── internal/
│   ├── handler/             # HTTPハンドラー
│   ├── service/             # ビジネスロジック
│   ├── repository/          # DB操作
│   └── scraper/             # データ取得
├── migrations/              # DBマイグレーション
├── Dockerfile
└── go.mod
```

**作業順序:**
1. プロジェクト構造作成 + Go mod init
2. PostgreSQL接続 + マイグレーション設定
3. アプリ情報取得API (`GET /api/apps/:id`)
4. キーワードランキング取得ロジック
5. Cloud Run用Dockerfile作成

### Phase 2: スクレイパー + 定期実行

1. App Store / Google Play スクレイパー実装
2. Cloud Scheduler連携用エンドポイント
3. ランキング履歴保存ロジック

### Phase 3: Frontend

```
frontend/
├── src/
│   ├── app/                 # Next.js App Router
│   ├── components/          # UIコンポーネント
│   └── lib/                 # API クライアント
├── next.config.js
└── package.json
```

**作業順序:**
1. Next.js プロジェクト作成
2. ダッシュボード画面
3. アプリ登録・管理画面
4. グラフ表示 (Chart.js or Recharts)

### Phase 4: インフラ + CI/CD

1. GCPプロジェクト設定
2. Cloud SQL インスタンス作成
3. Cloud Run デプロイ
4. Firebase Hosting 設定
5. GitHub Actions で自動デプロイ

---

## CI/CD 自動デプロイ構成

### 推奨: GitHub Actions

```
.github/workflows/
├── backend.yml    # Go → Cloud Run
└── frontend.yml   # Next.js → Firebase Hosting
```

#### Backend (Cloud Run) 自動デプロイ

```yaml
# .github/workflows/backend.yml
name: Deploy Backend

on:
  push:
    branches: [main]
    paths: ['backend/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}

      - uses: google-github-actions/setup-gcloud@v2

      - name: Build and Push to Artifact Registry
        run: |
          gcloud builds submit backend/ \
            --tag asia-northeast1-docker.pkg.dev/$PROJECT_ID/aso-tool/api

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy aso-api \
            --image asia-northeast1-docker.pkg.dev/$PROJECT_ID/aso-tool/api \
            --region asia-northeast1 \
            --allow-unauthenticated
```

#### Frontend (Firebase Hosting) 自動デプロイ

```yaml
# .github/workflows/frontend.yml
name: Deploy Frontend

on:
  push:
    branches: [main]
    paths: ['frontend/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install and Build
        working-directory: frontend
        run: |
          npm ci
          npm run build

      - uses: FirebaseExtended/action-hosting-deploy@v0
        with:
          repoToken: ${{ secrets.GITHUB_TOKEN }}
          firebaseServiceAccount: ${{ secrets.FIREBASE_SERVICE_ACCOUNT }}
          channelId: live
          projectId: your-project-id
```

### 必要なSecrets設定

| Secret名 | 用途 |
|---|---|
| `GCP_SA_KEY` | Cloud Run デプロイ用サービスアカウントJSON |
| `FIREBASE_SERVICE_ACCOUNT` | Firebase Hosting デプロイ用 |

---

## 懸念点・検討事項

### 1. データ取得方法について

**問題:** App Store / Google Play の検索ランキングは公式APIで取得できない

**選択肢:**

| 方法 | メリット | デメリット |
|---|---|---|
| **A. 自前スクレイピング** | 無料、自由度高い | 規約違反リスク、変更に弱い |
| **B. 外部API利用** (AppTweak, Sensor Tower等) | 安定、合法 | 有料 ($50〜/月) |
| **C. App Store Connect API のみ** | 公式、安全 | 自分のアプリのみ、ランキング取得不可 |

**おすすめ:** 最初は **C (公式API)** で自分のアプリのダウンロード数・レビューを取得。
ランキング追跡は後から検討（手動入力 or 外部API）

### 2. レビュー取得

- **App Store Connect API** → 自分のアプリのレビュー取得可能
- **Google Play Developer API** → 同様に可能
- 競合のレビューはRSS or スクレイピング必要

### 3. 本当にこの順番でいい？

現在の計画は **Backend → Scraper → Frontend → Infra** だが、
**もう一つの案:**

```
代替案: Infra → Backend → Frontend の順

1. 先にGCP/Firebaseセットアップ
2. Cloud SQL作成（これがないとBackend開発しにくい）
3. Backend開発（ローカルからCloud SQLに接続）
4. Frontend開発
```

**理由:** Cloud SQLがないとDB接続テストできない。
ローカルPostgreSQLでも開発できるが、本番環境との差異が出る可能性。

---

## 推奨する進め方

```
Step 1: ローカル開発環境
  └─ Docker Compose で PostgreSQL + Go API をローカル起動

Step 2: Backend MVP
  └─ CRUD API + 簡易スクレイパー

Step 3: GCPセットアップ + CI/CD
  └─ Cloud SQL作成、GitHub Actions設定

Step 4: Frontend
  └─ ダッシュボード実装

Step 5: 定期実行
  └─ Cloud Scheduler設定
```

これなら **ローカルで動くものを先に作って** から本番デプロイできる。

---

## 開発コマンド (予定)

```bash
# ローカル起動
docker-compose up -d

# Backend
cd backend && go run cmd/api/main.go

# Frontend
cd frontend && npm run dev

# マイグレーション
cd backend && go run cmd/migrate/main.go up
```

---

## 次のアクション

- [ ] `docker-compose.yml` 作成 (PostgreSQL + API)
- [ ] Go プロジェクト初期化
- [ ] DBスキーマ設計
- [ ] 最初のAPIエンドポイント作成
