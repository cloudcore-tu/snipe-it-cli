# snipe-it-cli

Snipe-IT（IT資産管理 OSS）を操作する Go 製 CLI ツール。コマンド名は `sit`。
インフラエンジニアの日常運用を第一に、Claude Code 等のコーディングエージェントからの利用も想定する。

## Snipe-IT について

IT資産（PC・サーバー・ライセンス等）を統合管理する OSS。PHP/Laravel 製。
全 API リソース・フィールド定義は公式 API リファレンスを参照すること。静的にここへ書かない（バージョンアップで陳腐化するため）。

```
公式ドキュメント: https://snipe-it.readme.io/docs/introduction
REST API リファレンス: https://snipe-it.readme.io/reference/api-overview
```

### REST API の基本構造

```
ベース URL: https://<hostname>/api/v1/
リスト:  GET    /api/v1/{resource}?limit=50&offset=0
詳細:   GET    /api/v1/{resource}/{id}
作成:   POST   /api/v1/{resource}
更新:   PATCH  /api/v1/{resource}/{id}
削除:   DELETE /api/v1/{resource}/{id}
```

**認証:**
```
Authorization: Bearer {api_token}
Accept: application/json
Content-Type: application/json
```
API トークンは Snipe-IT 管理画面 Settings > API Keys で生成する。

**ページネーション:**
```json
{
  "total": 250,
  "rows": [...]
}
```

クエリパラメータ: `limit`（件数）、`offset`（開始位置）

**Create/Update レスポンス:**
```json
{
  "status": "success",
  "messages": "Asset was created successfully.",
  "payload": { ...resource... }
}
```
クライアント層で `payload` を取り出して返す。

**レート制限:**
デフォルト 120 リクエスト/分。429 時は Retry-After を参照。

### 主要 API リソース一覧

| CLI コマンド | API パス | 説明 |
|---|---|---|
| `assets` | `/api/v1/hardware` | IT 資産（ハードウェア） |
| `users` | `/api/v1/users` | ユーザー |
| `licenses` | `/api/v1/licenses` | ソフトウェアライセンス |
| `accessories` | `/api/v1/accessories` | アクセサリー |
| `consumables` | `/api/v1/consumables` | 消耗品 |
| `components` | `/api/v1/components` | コンポーネント |
| `categories` | `/api/v1/categories` | カテゴリ |
| `companies` | `/api/v1/companies` | 会社 |
| `locations` | `/api/v1/locations` | ロケーション |
| `manufacturers` | `/api/v1/manufacturers` | メーカー |
| `models` | `/api/v1/models` | 機器モデル |
| `departments` | `/api/v1/departments` | 部門 |
| `statuslabels` | `/api/v1/statuslabels` | ステータスラベル |
| `suppliers` | `/api/v1/suppliers` | サプライヤー |
| `fieldsets` | `/api/v1/fieldsets` | カスタムフィールドセット |
| `maintenances` | `/api/v1/maintenances` | メンテナンス記録 |

## バージョニング方針

kubectl・gh に倣い、**snipe-it-cli 独自の semver** を採用する。Snipe-IT のバージョンには追従しない。

| 項目 | 方針 |
|------|------|
| snipe-it-cli 自身のバージョン | 独自 semver（v0.1.0〜） |
| 対応 Snipe-IT バージョン | README に明記（現在: v8.x 系） |
| バージョン表示 | `sit version` で CLI バージョン + 対応 Snipe-IT バージョンを表示 |

## グランドデザイン

### 想定ユーザー

- **人間のインフラエンジニア**（第一優先）: ターミナルから素早く確認・操作する
- **コーディングエージェント**（Claude Code 等）: Snipe-IT 操作を自動化する（`-o json` を明示指定）

kubectl / gh のような汎用 CLI として設計する。デフォルトは人間可読な出力とし、エージェントは `-o json` を明示指定する。

### コマンド体系

```
sit [global flags] {resource} {verb} [flags]

例:
sit assets list --status "Ready to Deploy"
sit assets get --id 123
sit assets create --data '{"name":"Laptop-001","asset_tag":"ASSET-001","model_id":1,"status_id":2}'
sit assets update --id 123 --data '{"status_id":3}'
sit assets delete --id 123 --yes
sit assets checkout --id 123 --data '{"checkout_to_type":"user","assigned_user":1}'
sit assets checkin --id 123
sit users list
sit licenses list
```

動詞は `list` / `get` / `create` / `update` / `delete` に統一する。  
資産固有の操作（`checkout` / `checkin`）は assets コマンドに追加する。

### 出力設計

| 用途 | 出力形式 |
|------|---------|
| デフォルト | table（人間向け） |
| `--output json` | JSON（エージェント・スクリプト向け） |
| `--output yaml` | YAML |
| `--output custom-columns=ID:.id,NAME:.name` | 列指定 |
| `--output 'jsonpath={rows.#.id}'` | 値抽出（gjson 構文） |

正常時の結果は stdout に出力する。エラーは stderr にプレーンテキストで出力し、exit code で成否を表す。対話的な確認プロンプトは出さない（`--yes` フラグで制御）。

### バージョン表示

```bash
$ sit version
snipe-it-cli v0.1.0
Snipe-IT API: v1 (compatible with Snipe-IT v8.x)
```

## 基本姿勢

- 推測禁止。不明な仕様・用語・技術判断は確認する
- 実装前に調査する。まず公式ドキュメントを確認し、不明な場合は GitHub Issues を調査する
- 変更対象だけでなく関連モジュール・呼び出し元も確認してから着手する

## 設計方針

- **DDD**: ドメイン層（Snipe-IT リソース操作）はインフラ層（HTTP クライアント・設定）に依存しない
- **TDD**: テストを先に書く。Red → Green → Refactor
- 依存はインターフェースを介する。具体実装に直接依存しない（テスト時にモック差し替え可能）
- 設定値のハードコード禁止。URL・認証情報・タイムアウト等はすべて環境変数・設定ファイルから注入

### kubectl から採用する実装パターン

**Options パターン（全コマンド共通）**

各コマンドは `XXXOptions` 構造体と `Complete → Validate → Run` の3段階で実装する。

```go
type ListAssetsOptions struct {
    Status string
    Output string
    // ...
}

func (o *ListAssetsOptions) Complete(cmd *cobra.Command) error { /* フラグ → 構造体へ移し替え、クライアント初期化 */ }
func (o *ListAssetsOptions) Validate() error                  { /* 引数の矛盾チェック */ }
func (o *ListAssetsOptions) Run(ctx context.Context) error    { /* メイン処理 */ }
```

**PrintFlags（出力フォーマット）**

`--output` フラグを `PrintFlags` 構造体として共通化し、各コマンドに渡す。

```go
type PrintFlags struct {
    OutputFormat string // table | json | yaml（デフォルト: table）
}
```

**ResourceDef（汎用 CRUD フレームワーク）**

`cmd/internal/run/resource.go` に定義。`ResourceDef` を宣言するだけで list/get/create/update/delete コマンドが自動生成される。

Snipe-IT API はリソースごとに一貫したパターンを持つため、型付き SDK なしに APIPath から汎用 HTTP メソッドを呼ぶ設計にしている（ADR-001 参照）。

```go
type ResourceDef struct {
    Use     string
    Short   string
    DocsURL string
    APIPath string      // "hardware", "users" 等
    ActionFns []ActionDef // checkout/checkin 等の追加アクション
}
```

**エラー構造化**

API エラーをユーザー向けメッセージに変換し、stderr にプレーンテキストで出力する。

```
exit 0  : 成功
exit 1  : エラー（API エラー、認証失敗、引数の誤り等）
```

```
$ sit assets get --id 99999
Error: API error: 404 Not Found
detail: No asset matches the given query.
```

## ツール・ビルド

- **言語**: Go 1.26.2（mise 管理）
- **CLI フレームワーク**: cobra
- **設定管理**: `go.yaml.in/yaml/v3` + `os.Getenv`（viper 不使用）
- **Snipe-IT クライアント**: 直接 HTTP（公式 Go SDK なし。既存ライブラリは star 数基準未達のため採用外。ADR-001 参照）
- **JSON パス**: `github.com/tidwall/gjson`（jsonpath/custom-columns 出力用）
- **テスト**: go test + testify
- **ロギング**: log/slog（構造化ログ）
- **lint**: golangci-lint

### コマンド

```bash
mise install           # Go バージョンをインストール
go build ./...         # ビルド
go test ./...          # テスト実行
golangci-lint run      # lint
```

## 設定管理

設定ファイル: `$XDG_CONFIG_HOME/snipe-it-cli/config.yaml`（未設定時は `~/.config/snipe-it-cli/config.yaml`）

複数インスタンスを管理できる（ADR-003 参照）。

```yaml
current: prod
instances:
  prod:
    url: https://snipeit.example.com
    token: your-api-token
  staging:
    url: https://staging.example.com
    token: stg-api-token
timeout: 30   # オプション
output: table  # オプション
```

| 設定項目 | 環境変数 | 設定ファイルキー |
|---------|---------|----------------|
| Snipe-IT URL | `SNIPEIT_URL` | `instances.<name>.url` |
| API トークン | `SNIPEIT_TOKEN` | `instances.<name>.token` |
| タイムアウト（秒） | `SNIPEIT_TIMEOUT` | `timeout` |
| デフォルト出力形式 | `SNIPEIT_OUTPUT` | `output` |
| アクティブなインスタンス（セッション） | `SNIPE_PROFILE` | — |
| デフォルトインスタンス | — | `current` |

設定値の優先順位: `CLI フラグ > 環境変数 > 設定ファイル（選択インスタンス） > デフォルト値`

### 設定管理コマンド

```bash
sit config init [--name NAME] --url URL --token TOKEN  # 初期設定ファイル生成
sit config add NAME --url URL --token TOKEN            # インスタンス追加・更新
sit config list                                        # 一覧表示（* がアクティブ）
export SNIPE_PROFILE=NAME                                # セッション単位で切り替え
sit --profile NAME ...                                 # コマンド単位で切り替え
```

## ロギング方針

- `log/slog` で構造化ログを出力する
- ログレベル: ERROR / WARN / INFO / DEBUG（`--debug` フラグで有効化）
- API トークン等の秘密情報はログに出力しない

## Git 運用

Contextual Commits 形式でコミットする。

```text
type(scope): 簡潔な説明

## なぜ
この変更が必要な背景・理由

## 何を
具体的な変更内容

## 影響
他モジュールへの影響、注意点（なければ省略）

Refs: #issue-number（あれば）
```

type: feat, fix, docs, refactor, test, chore, ci

## ライブラリ選定基準

- GitHub star 1,000 以上 / 最終更新 2 ヶ月以内 / 開発元が信頼できる
- 基準を満たさない場合は導入前に確認を求める

## ドキュメント更新

あらゆるコード変更・設計変更時に、関連ドキュメントを必ず同時に更新する。ドキュメントの更新漏れはバグと同じ扱い。

### CLAUDE.md（最優先）

以下に該当する変更は CLAUDE.md を更新する。

- ツール・ライブラリの追加・削除・バージョン変更
- ビルド・テスト・lint のコマンドや設定の変更
- ディレクトリ構成・設計方針・グランドデザインの変更
- 設定項目の追加・変更

### その他のドキュメント

- 新しい設計判断は `docs/adr/` に ADR を作成する
- 機能追加・変更・削除時は `CHANGELOG.md` を更新する（Keep a Changelog 形式）

## スキル運用

| スキル | 用途 |
|--------|------|
| `contextual-commit` | Contextual Commits 形式でのコミット |
| `changelog-update` | CHANGELOG.md の更新 |
| `adr-writer` | Architecture Decision Record の作成 |
| `doc-sync` | コード変更に対するドキュメント整合性チェック |

## ディレクトリ構成

```text
.claude/
  CLAUDE.md             # プロジェクト固有ルール
  skills/               # 繰り返しワークフロー定義

cmd/
  root.go               # ルートコマンド（cobra）、グローバルフラグ定義（--profile 含む）
  version.go            # sit version（CLI バージョン + 対応 Snipe-IT バージョン）
  config/               # sit config サブコマンド群
    config.go           # config コマンドのエントリポイント
    init.go             # sit config init
    add.go              # sit config add NAME
    list.go             # sit config list
  internal/run/
    run.go              # BaseOptions / ParseFilters / FormatAPIError 等の共通ユーティリティ
    resource.go         # ResourceDef / ActionDef（汎用 CRUD フレームワーク）
  assets/               # sit assets（/api/v1/hardware）
    assets.go           # ResourceDef + checkout/checkin 追加コマンド
  users/                # sit users
    users.go
  licenses/             # sit licenses
    licenses.go
  {resource}/           # その他リソースは ResourceDef のみ
    {resource}.go

internal/
  snipeit/              # Snipe-IT API クライアント（直接 HTTP）
    client.go           # NewClient、List/GetByID/Create/Update/Delete 汎用メソッド
  config/               # 設定管理（go.yaml.in/yaml/v3 + os.Getenv）
    config.go           # FileConfig / Instance / Config 型定義
  output/               # 出力フォーマット（table / JSON / YAML / custom-columns / jsonpath）
    output.go

docs/
  adr/                  # Architecture Decision Records
  api-coverage.md       # API カバレッジ一覧

CHANGELOG.md            # 変更履歴（Keep a Changelog 形式）
mise.toml               # ツールバージョン管理
```
