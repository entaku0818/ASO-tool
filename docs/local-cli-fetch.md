# ローカルCLIからASO順位を取得する手順

ローカル（macOS）のシェルスクリプトやセッション開始フックから、ASO Compass のバックエンドAPI経由で
自分のアプリのキーワード順位を取得するための手順書。**DBは直接叩かず、公開APIのみを使う。**

実測日: 2026-08-02（すべて本番APIに対して実際にcurlを叩いて確認済み）

> **秘匿値についての注意**: このファイルにライセンスキー・JWT・メールアドレスの実値は書かない。
> 実値は macOS Keychain から取る（末尾「Keychain」節）。サンプルJSONの該当箇所はマスクしてある。

---

## 1. ベースURL

```
https://aso-api-671942133800.asia-northeast1.run.app
```

環境変数化しておくと扱いやすい: `ASO_API_BASE`

疎通確認は認証不要の `GET /health` が使える。

---

## 2. 認証フロー（ライセンスキー → アクティベート → JWT）

macOSアプリと同じ方式。ユーザー名/パスワードは不要で、**ライセンスキー＋メールアドレス**から
JWT を発行する。

### 2.1 アクティベート

```
POST /api/licenses/activate
Content-Type: application/json
（Authorization ヘッダは付けない）
```

**リクエストボディ**（フィールド名に注意。`license_key` ではなく **`key`**）

```json
{
  "key": "<ライセンスキー>",
  "email": "<メールアドレス>"
}
```

| フィールド | 型 | 必須 | 備考 |
|---|---|---|---|
| `key` | string | ✅ | 間違えて `license_key` にすると `400 {"error":"license key is required"}` |
| `email` | string | ✅ | ライセンスに紐づくメール |

**レスポンス `200 OK`**

```json
{
  "token": "eyJhbGciOi…<MASKED>",
  "user": {
    "id": "1c6afa63-…",
    "email": "<MASKED>",
    "name": "<MASKED>",
    "is_admin": false,
    "plan": "free",
    "created_at": "2026-06-20T01:06:13.035303Z"
  },
  "key": "ASOT-…<MASKED>"
}
```

| フィールド | 型 | 備考 |
|---|---|---|
| `token` | string | これが JWT。以降 `Authorization: Bearer <token>` で使う |
| `user.id` | string (UUID) | |
| `user.plan` | string | `free` / `pro` |
| `key` | string | 入力したライセンスキーがそのまま返る |

**エラー**: `400` = 必須項目欠落 / `404 license not found` = キー不正 /
`409` 相当 = 別メールで既にアクティベート済み（`ErrLicenseAlreadyUsed`）

### 2.2 JWTの有効期限

**24時間**（`backend/internal/service/auth.go:111` — `time.Now().Add(24 * time.Hour)`）。

`activate` は**同一メールなら冪等**（別メールが既にそのライセンスを使っている場合のみ拒否。
`backend/internal/service/license.go:67`）。したがって

> **毎回の実行で activate を呼び直して新しい JWT を取るのが最も単純で安全。**
> トークンをローカルにキャッシュする必要はない。

### 2.3 401 を受けたときの再アクティベート手順

トークンをキャッシュする実装にする場合は、macOSアプリと同じ挙動にする
（`macos/ASOCompass/ASO Compass/Services/APIClient.swift`）:

1. 任意のAPI呼び出しが `401` を返す
2. **保存済みのライセンスキー＋メールで `POST /api/licenses/activate` をやり直す**
   （このエンドポイントは Authorization ヘッダを付けないので、401リトライが再帰しない）
3. 返ってきた新 `token` を保存する
4. 元のリクエストの `Authorization` ヘッダだけ差し替えて **1回だけ**リトライする
5. activate 自体が `4xx`（ライセンス無効・失効）を返したら手動での再認証が必要。
   `5xx` やオフラインの場合は保存済みライセンスを消さずにそのまま失敗させる

---

## 3. アプリ一覧

```
GET /api/apps
Authorization: Bearer <token>
```

**レスポンス `200 OK`**（配列）

```json
[
  {
    "id": "b56b21d4-ccb7-4201-8718-6d5091b61f5e",
    "name": "シンプル録音 - 高音質ボイスレコーダー",
    "bundle_id": "com.entaku.VoiLog",
    "platform": "ios",
    "store_url": "https://apps.apple.com/jp/app/…/id6443528409?uo=4",
    "user_id": "1c6afa63-4de1-4bc2-a98c-df36335a60e1",
    "created_at": "2026-06-21T02:11:27.114Z",
    "updated_at": "2026-06-21T02:11:27.114Z"
  }
]
```

以降のエンドポイントで使う `{appID}` はこの `id`（UUID）。

---

## 4. 順位取得エンドポイント

### 4.1 `GET /api/apps/{appID}/keywords/ranks` — 最新順位＋変化量

```
GET /api/apps/{appID}/keywords/ranks
Authorization: Bearer <token>
```

パラメータなし。

**レスポンス `200 OK`**（配列。`keyword` 昇順）

```json
[
  {
    "keyword_id": "d3fbfd86-f748-4be0-ab6a-c402ab980810",
    "keyword": "録音",
    "country": "JP",
    "current_rank": 12,
    "previous_rank": 12,
    "change": 0
  }
]
```

| フィールド | 型 | 備考 |
|---|---|---|
| `keyword_id` | string (UUID) | |
| `keyword` | string | |
| `country` | string | 例 `JP` |
| `current_rank` | int \| **null** | **null = 圏外**（検索結果に自Appが出てこなかった）。詳細は 4.3 |
| `previous_rank` | int \| **null** | 履歴が1件しかなければ null |
| `change` | int \| **null** | `previous_rank - current_rank`。**正の値 = 順位が上がった** |

⚠️ **`previous_rank` は「前日」ではなく「1つ前のスクレイプ結果」**
（`backend/internal/repository/ranking.go:248` の `ROW_NUMBER() … WHERE rn = 2`）。
日次バッチ以外に手動スクレイプ（macOSアプリを開く、`POST /scrape/rankings`）が走ると
同じ日に複数行が入り、`change` が「数分前との差」になる。
**正確な前日比が必要なら 4.2 を使ってクライアント側で日付集計すること。**

### 4.2 `GET /api/apps/{appID}/rankings?days=N` — 日付付き履歴（前日比の計算用・推奨）

```
GET /api/apps/{appID}/rankings?days=8
Authorization: Bearer <token>
```

| パラメータ | 既定値 | 備考 |
|---|---|---|
| `days` | 30 | 何日分さかのぼるか |

**レスポンス `200 OK`**（配列。そのアプリの**全キーワード**の履歴が混ざって返る）

```json
[
  {
    "id": "e9901d9e-4f8d-4f39-b430-be1868529a96",
    "keyword_id": "d3fbfd86-f748-4be0-ab6a-c402ab980810",
    "rank": 12,
    "recorded_at": "2026-08-02T12:47:23.013875Z",
    "keyword": "録音",
    "country": "JP"
  }
]
```

| フィールド | 型 | 備考 |
|---|---|---|
| `id` | string (UUID) | ranking_history の行ID |
| `keyword_id` | string (UUID) | |
| `rank` | int \| **null** | |
| `recorded_at` | string (RFC3339, **UTC**) | JSTにするには +9h |
| `keyword` | string | |
| `country` | string | |

### 4.3 `rank: null`（圏外）の扱い ※必読

`rank` / `current_rank` が **`null` の場合は「圏外」**を意味する。エラーでもデータ欠損でもない。

`backend/internal/service/scraper.go:82` の `var rank *int` がそのまま保存されており、
App Store の検索結果に自分のアプリが出てこなかったキーワードは `rank = NULL` で
**行自体はちゃんと書き込まれる**。したがって:

- `null` を「取得失敗」として扱わないこと
- `null` を `0` に落とすと「1位より上」に見えてしまうので絶対にやらない
- 前日比は**両日とも数値のときだけ**計算し、片方でも `null` なら「—」等で逃がす

実測（2026-08-02 時点）: 11アプリ / 79キーワード中、圏内27語・圏外52語。
新規登録したてのアプリは大半が圏外になるのが普通。

### 4.4 その他

| パス | パラメータ | 返すもの |
|---|---|---|
| `GET /api/apps/{appID}/keywords` | — | `[{id, app_id, keyword, country, created_at}]` |
| `GET /api/apps/{appID}/keywords/{keywordID}/rankings` | `days` or `limit` | `[{id, keyword_id, rank, recorded_at}]` |
| `GET /api/apps/{appID}/keywords/{keywordID}/rankings/latest` | — | 上記の単一オブジェクト |
| `GET /api/apps/{appID}/keywords/rising` | `days`(既定7) | `[{keyword_id, keyword, country, current_rank, previous_rank, improvement}]` |
| `GET /api/tracked-keywords` | — | `[{id, keyword, country, platform, user_id, created_at}]` |
| `GET /api/tracked-keywords/{id}/results` | — | SERP上位（自社以外も含む）`[{id, tracked_keyword_id, rank, app_name, bundle_id, developer, recorded_at}]` |

`GET /api/apps/{appID}/competitors/keyword-gap` は **Proプラン必須**（free だと
`403 {"error":"keyword gap analysis requires Pro plan"}`）。

---

## 5. 前日比はサーバ側？クライアント側？

**両方ある。用途で使い分ける。**

| | 前日比 | 判定 |
|---|---|---|
| `/keywords/ranks` の `change` | サーバが計算して返す | ただし「前日」ではなく**「1つ前のスクレイプ」との差**。手動スクレイプが混ざると壊れる |
| `/rankings?days=N` | 生の履歴のみ | **クライアント側で計算が必要**。正確な日次比較ができる |

朝イチのフックで**確実な前日比**を出すなら後者。推奨アルゴリズム:

1. `GET /api/apps/{appID}/rankings?days=8`
2. `recorded_at` を +9h して **JSTの日付**に丸める
3. `(keyword, country)` × JST日付 でグルーピングし、**同じ日に複数行あれば最後の1件**を採用
4. 直近日と、その1つ前の**別の日**を比較 → `前の日のrank - 直近のrank`（正 = 上昇）

---

## 6. 日次バッチのタイミング（重要）

`.github/workflows/batch-scheduler.yml` の cron は `0 0 * * *`（= **9:00 JST**）だが、
docker build/push → Cloud Run Job 実行を挟むため、**実際にDBへ書き込まれるのは 11:00 JST 前後**。

2026-07-20〜08-02 の実測 `recorded_at`（JST）:

```
07-20 11:32 / 07-21 11:10 / 07-22 11:09 / 07-23 11:15 / 07-24 11:11 / 07-25 11:11
07-26 11:21 / 07-27 11:29 / 07-28 11:04 / 07-29 11:07 / 07-30 10:57 / 07-31 11:22
08-01 11:23 / 08-02 11:14
```

> **朝9時にフックを回すと当日分はまだ入っていない**（最新＝前日分になる）。
> 当日分が欲しければ 11:30 JST 以降に実行すること。
> 「最新の1件」を当日と決め打ちせず、`recorded_at` の日付を必ず表示するのが安全。

---

## 7. Keychain（認証情報の置き場所）

平文でリポジトリに置かない。`credential-vault` スキル経由で macOS Keychain に保存する。

```bash
VAULT=~/.claude/skills/credential-vault/scripts/credential.sh
```

| 項目 | 値 |
|---|---|
| **service名** | `aso-tool-license` |
| **account名** | ASO Compass に登録したメールアドレス（＝`activate` の `email`） |
| **password（値）** | ライセンスキー（`ASOT-` で始まる文字列） |

取り出し方:

```bash
EMAIL=$($VAULT account aso-tool-license)   # → メールアドレス
KEY=$($VAULT get aso-tool-license)         # → ライセンスキー
```

**値の取得元**: macOSアプリ ASO Compass が保存している UserDefaults。

```bash
defaults read com.entaku.ASOCompass user_email    # メールアドレス
defaults read com.entaku.ASOCompass license_key   # ライセンスキー
```

（未登録の場合の登録コマンド。値は標準入力で渡し、シェル履歴に残さない）

```bash
defaults read com.entaku.ASOCompass license_key \
  | $VAULT set aso-tool-license "$(defaults read com.entaku.ASOCompass user_email)"
```

---

## 8. 参考: 動作確認済みの最小スクリプト

```bash
#!/usr/bin/env bash
set -euo pipefail

BASE="${ASO_API_BASE:-https://aso-api-671942133800.asia-northeast1.run.app}"
VAULT="$HOME/.claude/skills/credential-vault/scripts/credential.sh"

EMAIL="$("$VAULT" account aso-tool-license)"
KEY="$("$VAULT" get aso-tool-license)"

JWT="$(curl -sS -X POST "$BASE/api/licenses/activate" \
  -H 'Content-Type: application/json' \
  -d "$(printf '{"key":"%s","email":"%s"}' "$KEY" "$EMAIL")" \
  | python3 -c 'import json,sys; print(json.load(sys.stdin)["token"])')"

curl -sS "$BASE/api/apps" -H "Authorization: Bearer $JWT" \
| python3 -c '
import json, sys, subprocess, datetime
base, jwt = sys.argv[1], sys.argv[2]
apps = json.load(sys.stdin)
if not apps:
    print("登録アプリなし"); sys.exit()

def api(path):
    out = subprocess.run(["curl","-sS",base+path,"-H","Authorization: Bearer "+jwt],
                         capture_output=True, text=True, check=True).stdout
    return json.loads(out)

def jst_date(iso):
    t = datetime.datetime.fromisoformat(iso.replace("Z","+00:00"))
    return (t + datetime.timedelta(hours=9)).date()

for app in apps:
    print("\n■ " + app["name"])
    rows = api("/api/apps/" + app["id"] + "/rankings?days=8")
    if not rows:
        print("  順位データなし"); continue
    by_kw = {}
    for r in sorted(rows, key=lambda x: x["recorded_at"]):
        by_kw.setdefault((r["keyword"], r["country"]), {})[jst_date(r["recorded_at"])] = r["rank"]
    def fmt(r):            # rank は null = 圏外
        return "圏外" if r is None else f"{r}位"

    for (kw, country), series in sorted(by_kw.items()):
        days = sorted(series)
        cur, curd = series[days[-1]], days[-1]
        line = f"  {kw} ({country}): {fmt(cur)}"
        if len(days) >= 2:
            prev, prevd = series[days[-2]], days[-2]
            if prev is not None and cur is not None:
                diff = prev - cur
                mark = "→ ±0" if diff == 0 else (f"↑ +{diff}" if diff > 0 else f"↓ {diff}")
            else:
                mark = "—"     # どちらかが圏外で数値比較できない
            line += f"  {mark}  (前回 {fmt(prev)} / {prevd})"
        else:
            line += "  (比較データなし)"
        print(line + f"  [{curd}]")
' "$BASE" "$JWT"
```

実行例:

```
■ シンプル録音 - AI文字起こし&ボイスメモ
  録音 (JP): 12位  → ±0  (前回 12位 / 2026-08-01)  [2026-08-02]
  長時間録音 (JP): 4位  (比較データなし)  [2026-08-02]
  議事録 (JP): 圏外  (比較データなし)  [2026-08-02]

■ ClipKit - コピー履歴・クリップボード
  クリップボード (JP): 3位  (比較データなし)  [2026-08-02]
  履歴 (JP): 11位  (比較データなし)  [2026-08-02]
```

⚠️ 2026-08-02 に11アプリ・79キーワードを登録したばかりなので、多くのキーワードは
履歴が1日分しかなく「比較データなし」になる。**翌朝のバッチ以降から前日比が出る。**
