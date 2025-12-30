# Firebase Hosting Setup Guide

フロントエンド（Next.js）を Firebase Hosting にデプロイするためのセットアップガイドです。

## 前提条件

- GCP プロジェクト作成済み（[GCP Setup Guide](./gcp-setup.md) 参照）
- Firebase CLI インストール済み
- Node.js 20 以上

## 1. Firebase CLI のインストール

```bash
npm install -g firebase-tools

# ログイン
firebase login
```

## 2. Firebase プロジェクトの初期化

```bash
# GCP プロジェクトを Firebase に追加
firebase projects:addfirebase aso-tool-prod

# プロジェクトディレクトリで Firebase を初期化
cd frontend
firebase init hosting
```

初期化時の選択肢：
- **What do you want to use as your public directory?** → `out` (静的エクスポート) または `.next` (Next.js)
- **Configure as a single-page app?** → `No` (Next.js のルーティングを使用)
- **Set up automatic builds and deploys with GitHub?** → `No` (GitHub Actions を使用)

## 3. firebase.json の設定

```json
{
  "hosting": {
    "public": "out",
    "ignore": [
      "firebase.json",
      "**/.*",
      "**/node_modules/**"
    ],
    "rewrites": [
      {
        "source": "/api/**",
        "run": {
          "serviceId": "aso-api",
          "region": "asia-northeast1"
        }
      },
      {
        "source": "**",
        "destination": "/index.html"
      }
    ],
    "headers": [
      {
        "source": "**/*.@(js|css)",
        "headers": [
          {
            "key": "Cache-Control",
            "value": "public, max-age=31536000, immutable"
          }
        ]
      }
    ]
  }
}
```

## 4. サービスアカウントの作成

```bash
PROJECT_ID=$(gcloud config get-value project)

# Firebase Hosting 用のサービスアカウントを作成
gcloud iam service-accounts create firebase-hosting \
  --display-name="Firebase Hosting Deployer"

# 必要な権限を付与
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:firebase-hosting@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/firebasehosting.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:firebase-hosting@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# サービスアカウントキーを生成
gcloud iam service-accounts keys create ~/firebase-sa-key.json \
  --iam-account=firebase-hosting@$PROJECT_ID.iam.gserviceaccount.com
```

## 5. GitHub Secrets の設定

GitHub リポジトリの Settings > Secrets and variables > Actions で以下を追加:

| Secret 名 | 値 | 説明 |
|-----------|---|------|
| `FIREBASE_SERVICE_ACCOUNT` | (JSONファイルの内容) | `~/firebase-sa-key.json` の内容 |
| `PRODUCTION_API_URL` | `https://aso-api-xxxxx.a.run.app` | Cloud Run の URL |

### Cloud Run URL の取得

```bash
gcloud run services describe aso-api \
  --region=asia-northeast1 \
  --format="value(status.url)"
```

## 6. Next.js の設定

### next.config.js

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',  // 静的エクスポートを有効化
  trailingSlash: true,
  images: {
    unoptimized: true,  // 静的エクスポート時は必要
  },
}

module.exports = nextConfig
```

### package.json にスクリプト追加

```json
{
  "scripts": {
    "build": "next build",
    "export": "next export",
    "type-check": "tsc --noEmit"
  }
}
```

## 7. 手動デプロイ（テスト用）

```bash
cd frontend

# ビルド
npm run build

# デプロイ
firebase deploy --only hosting
```

## Preview チャンネル

PR ごとにプレビュー環境を作成することも可能です：

```bash
firebase hosting:channel:deploy preview-pr-123
```

## カスタムドメインの設定

```bash
# ドメインを追加
firebase hosting:sites:update aso-tool --project=aso-tool-prod

# Firebase Console でドメイン設定
# https://console.firebase.google.com/project/aso-tool-prod/hosting
```

DNS に以下のレコードを追加：
- A レコード: Firebase が提供する IP アドレス
- TXT レコード: 所有権確認用

## トラブルシューティング

### デプロイが失敗する

```bash
# ローカルでビルドを確認
npm run build

# Firebase のログを確認
firebase hosting:channel:list
```

### API 呼び出しが失敗する

1. CORS 設定を確認（Cloud Run 側）
2. `NEXT_PUBLIC_API_URL` が正しく設定されているか確認
3. firebase.json の rewrites 設定を確認

### 静的ファイルが見つからない

1. `output: 'export'` が設定されているか確認
2. `out` ディレクトリが生成されているか確認
3. firebase.json の `public` が `out` になっているか確認
