# GCP Setup Guide

ASO Compass を GCP にデプロイするためのセットアップガイドです。

## 前提条件

- Google Cloud アカウント
- `gcloud` CLI インストール済み
- GitHub リポジトリへのアクセス権限

## 1. GCP プロジェクトの作成

```bash
# プロジェクト作成
gcloud projects create aso-tool-prod --name="ASO Tool Production"

# プロジェクトを選択
gcloud config set project aso-tool-prod

# 課金アカウントをリンク (課金アカウントIDは gcloud billing accounts list で確認)
gcloud billing projects link aso-tool-prod --billing-account=BILLING_ACCOUNT_ID
```

## 2. 必要な API の有効化

```bash
# 必要な API を有効化
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  sqladmin.googleapis.com \
  secretmanager.googleapis.com \
  cloudbuild.googleapis.com
```

## 3. Artifact Registry リポジトリの作成

```bash
# Docker イメージ用のリポジトリを作成
gcloud artifacts repositories create aso-tool \
  --repository-format=docker \
  --location=asia-northeast1 \
  --description="ASO Compass Docker images"
```

## 4. Cloud SQL (PostgreSQL) のセットアップ

```bash
# Cloud SQL インスタンスを作成 (db-f1-micro は最小構成)
gcloud sql instances create aso-db \
  --database-version=POSTGRES_16 \
  --tier=db-f1-micro \
  --region=asia-northeast1 \
  --storage-type=SSD \
  --storage-size=10GB

# データベースを作成
gcloud sql databases create aso_tool --instance=aso-db

# ユーザーを作成
gcloud sql users create aso \
  --instance=aso-db \
  --password=YOUR_SECURE_PASSWORD
```

### Cloud SQL 接続情報の取得

```bash
# インスタンスの接続名を取得
gcloud sql instances describe aso-db --format="value(connectionName)"
# 出力例: aso-tool-prod:asia-northeast1:aso-db
```

## 5. サービスアカウントの作成

```bash
# サービスアカウントを作成
gcloud iam service-accounts create github-actions \
  --display-name="GitHub Actions Deployer"

# プロジェクトIDを変数に設定
PROJECT_ID=$(gcloud config get-value project)

# 必要な権限を付与
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

# サービスアカウントキーを生成
gcloud iam service-accounts keys create ~/gcp-sa-key.json \
  --iam-account=github-actions@$PROJECT_ID.iam.gserviceaccount.com
```

## 6. GitHub Secrets の設定

GitHub リポジトリの Settings > Secrets and variables > Actions で以下を設定:

| Secret 名 | 値 | 説明 |
|-----------|---|------|
| `GCP_PROJECT_ID` | `aso-tool-prod` | GCP プロジェクト ID |
| `GCP_SA_KEY` | (JSONファイルの内容) | `~/gcp-sa-key.json` の内容をコピー |
| `DATABASE_URL` | (下記参照) | Cloud SQL 接続文字列 |

### DATABASE_URL のフォーマット

Cloud Run から Cloud SQL に接続する場合、Unix ソケット経由で接続します：

```
postgres://aso:YOUR_PASSWORD@/aso_tool?host=/cloudsql/PROJECT_ID:asia-northeast1:aso-db
```

例：
```
postgres://aso:secure_password@/aso_tool?host=/cloudsql/aso-tool-prod:asia-northeast1:aso-db
```

## 7. Cloud Run の Cloud SQL 接続設定

初回デプロイ後、Cloud Run サービスに Cloud SQL 接続を追加：

```bash
# Cloud SQL 接続を追加
gcloud run services update aso-api \
  --region=asia-northeast1 \
  --add-cloudsql-instances=aso-tool-prod:asia-northeast1:aso-db
```

または、GitHub Actions ワークフローの deploy コマンドに追加：

```yaml
gcloud run deploy aso-api \
  --add-cloudsql-instances=$PROJECT_ID:asia-northeast1:aso-db \
  # ... 他のオプション
```

## 8. マイグレーションの実行

Cloud SQL にマイグレーションを実行するには、Cloud SQL Auth Proxy を使用：

```bash
# Cloud SQL Auth Proxy をダウンロード
curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.8.0/cloud-sql-proxy.darwin.arm64
chmod +x cloud-sql-proxy

# プロキシを起動 (別ターミナルで)
./cloud-sql-proxy aso-tool-prod:asia-northeast1:aso-db --port=5433

# マイグレーションを実行
DATABASE_URL="postgres://aso:YOUR_PASSWORD@localhost:5433/aso_tool" \
  go run ./cmd/migrate up
```

## セキュリティのベストプラクティス

1. **サービスアカウントキーの管理**
   - キーファイルは生成後すぐに削除
   - GitHub Secrets に保存したら `rm ~/gcp-sa-key.json`

2. **最小権限の原則**
   - サービスアカウントには必要最小限の権限のみ付与

3. **Cloud SQL**
   - パブリック IP を無効化し、プライベート IP のみ使用を推奨
   - 本番環境では強力なパスワードを使用

## コスト見積もり（概算）

| サービス | 構成 | 月額費用（概算） |
|---------|------|-----------------|
| Cloud Run | 最小構成、低トラフィック | $0 〜 $10 |
| Cloud SQL | db-f1-micro, 10GB | 約 $10 |
| Artifact Registry | 1GB 以下 | 約 $0.1 |
| **合計** | | **約 $10 〜 $20** |

※ 無料枠を超えない範囲では $0 になる可能性もあります

## トラブルシューティング

### デプロイが失敗する

```bash
# ログを確認
gcloud run services logs read aso-api --region=asia-northeast1
```

### データベース接続エラー

1. Cloud SQL インスタンスが起動しているか確認
2. Cloud Run サービスに Cloud SQL 接続が追加されているか確認
3. DATABASE_URL のフォーマットが正しいか確認

### 権限エラー

```bash
# サービスアカウントの権限を確認
gcloud projects get-iam-policy $PROJECT_ID \
  --filter="bindings.members:github-actions@$PROJECT_ID.iam.gserviceaccount.com"
```
