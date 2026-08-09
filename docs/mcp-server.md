# ASO Compass MCP サーバ

Claude Code / Claude Desktop などの MCP クライアントから、ASO Compass のデータを直接読めるようにする
stdio ベースの MCP サーバ。

**設計方針**: バックエンドが LLM を呼んで「AIアドバイス」を生成するのをやめ、代わりに**データを
MCP で公開して、分析は AI クライアント側にやらせる**。分析の内容・粒度・切り口をその場で指示できるし、
バックエンドに ANTHROPIC_API_KEY を置く必要もなくなる。

- 実装: `backend/cmd/mcp/`
- トランスポート: stdio（JSON-RPC 2.0）
- SDK: `github.com/modelcontextprotocol/go-sdk` v1.0.0
- データ取得はすべて既存の REST API 経由。DBには直接触れない
  （API の詳細は [local-cli-fetch.md](local-cli-fetch.md)）

---

## 1. ビルド

```bash
cd backend
go build -o ~/bin/aso-mcp ./cmd/mcp
```

Go 1.24 でビルドできる（SDK は v1.0.0 に固定してある。v1.1.0 以降は Go 1.25 が必要で、
CI と Dockerfile の Go 1.24 と噛み合わないため上げていない）。

## 2. 設定

環境変数だけ。設定ファイルは無い。

| 変数 | 必須 | 内容 |
|---|---|---|
| `ASO_LICENSE_KEY` | ✅ | ライセンスキー（macOSアプリと同じもの） |
| `ASO_EMAIL` | ✅ | ライセンスに紐づくメールアドレス |
| `ASO_API_BASE` | — | 既定は本番。ローカル検証時のみ指定 |

値は Keychain から取る（service `aso-tool-license`。詳細は
[local-cli-fetch.md](local-cli-fetch.md) の Keychain 節）:

```bash
VAULT=~/.claude/skills/credential-vault/scripts/credential.sh
export ASO_LICENSE_KEY="$($VAULT get aso-tool-license)"
export ASO_EMAIL="$($VAULT account aso-tool-license)"
```

**認証の挙動**: 起動時にライセンスキーで `POST /api/licenses/activate` を叩いて JWT を取得する。
JWT は24時間で失効するが、`401` を受けたら自動で再アクティベートして**1回だけ**リトライする
（macOSアプリと同じ方式）。トークンをディスクにキャッシュはしない。

## 3. Claude Code への登録

```bash
claude mcp add aso-tool \
  --env ASO_LICENSE_KEY="$($VAULT get aso-tool-license)" \
  --env ASO_EMAIL="$($VAULT account aso-tool-license)" \
  -- ~/bin/aso-mcp
```

> ⚠️ ライセンスキーが `~/.claude.json` に平文で入る。気になる場合は
> `aso-mcp` を直接登録せず、Keychain から読んで `exec` するラッパースクリプトを噛ませる。

## 4. 提供ツール

| ツール | 引数 | 返すもの |
|---|---|---|
| `list_apps` | なし | 登録アプリ一覧。以降のツールに渡す `app_id` はここから |
| `list_keyword_ranks` | `app_id` | 全キーワードの最新順位・前日順位・変化量 |
| `get_rank_history` | `app_id`, `days`(既定30) | 日付付きの順位履歴。推移分析用 |
| `list_rising_keywords` | `app_id`, `days`(既定7) | 上昇したキーワードのみ、上昇幅の降順 |
| `search_app_store` | `keyword`, `country`(既定jp), `limit`(既定50) | App Store 検索結果。**未登録キーワードでも調べられる** |
| `get_keyword_gap` | `app_id` | 競合が取れていて自Appが取れていない語（Pro必須） |

### 読み取り専用

すべて `GET` のみ。アプリやキーワードの登録・削除・スクレイプ実行といった**状態を変える操作は
意図的に公開していない**。AI クライアントが誤ってデータを書き換える事故を防ぐため。
登録作業が必要なら REST API を直接叩く。

### 圏外の扱い

`current_rank` / `rank` が `null` は**圏外**（検索結果に自Appが出てこなかった）であって
取得失敗ではない。`0` に丸めると「1位より上」に見えるので絶対にやらないこと。
ツールの description にもこの旨を書いてあるので、モデルは基本的に正しく解釈する。

## 5. 使い方の例

Claude Code で自然文で聞けばよい。

```
シンプル録音の順位で、直近1週間で下がったキーワードを教えて
```
→ `list_apps` → `get_rank_history` を呼んで、モデルが日付集計して答える。

```
ClipKit が「クリップボード」で3位だけど、上位2つは何のアプリ？勝つには何が要る？
```
→ `search_app_store` で SERP を取って、モデルが競合分析する。

```
全アプリの圏内キーワードを数えて、どのアプリに注力すべきか整理して
```
→ `list_apps` → 各アプリに `list_keyword_ranks` → モデルが横断集計。

## 6. 動作確認

```bash
# ツール一覧が返るか
ASO_LICENSE_KEY=... ASO_EMAIL=... ~/bin/aso-mcp <<'EOF'
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"0"}}}
EOF
```

`{"result":{"serverInfo":{"name":"aso-tool","version":"0.1.0"},...}}` が返れば起動している。

**stdout は JSON-RPC 専用**。ログや診断は必ず stderr に出すこと（`log.SetOutput(os.Stderr)`）。
stdout に1行でも余計な出力を混ぜるとプロトコルが壊れる。

## 7. テスト

`backend/cmd/mcp/client_test.go` が認証まわりを固めている:

- 初回リクエスト前に activate すること / activate に `Authorization` を付けないこと
- リクエストボディのキーが `key`（`license_key` ではない）であること
- 401 で再アクティベートして新トークンでリトライすること
- リトライも401なら**それ以上リトライしない**こと
- 403 などのステータスコードと理由が呼び出し元に伝わること

```bash
cd backend && go test ./cmd/mcp/...
```

---

## 廃止した機能

`GET /api/apps/{appID}/aso-advice`（ASOアドバイス）とその macOS UI は**削除済み**。
バックエンドから Claude API を呼んで固定フォーマットの助言を生成する機能だったが、
MCP 経由で AI クライアントに分析させるほうが柔軟なため置き換えた。

- 削除: `backend/internal/service/aso_advice.go`, `backend/internal/handler/aso_advice.go`,
  `macos/ASOCompass/ASO Compass/Views/ASOAdviceView.swift`、および各所の導線
- `ANTHROPIC_API_KEY` は**キャプション生成**（`POST /api/apps/{appID}/captions/generate`）が
  まだ使っているので、環境変数自体は残す必要がある
