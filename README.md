# snipe-it-cli

Snipe-IT（IT資産管理 OSS）を操作する Go 製 CLI ツール。コマンド名: `snip`。対応 Snipe-IT: v8.x

| OS | アーキテクチャ |
|----|--------------|
| macOS | amd64（Intel）/ arm64（Apple Silicon）|
| Linux | amd64 / arm64 |
| Windows | amd64 / arm64 |

## インストール

### Homebrew（macOS / Linux）

```bash
brew tap cloudcore-tu/tap https://github.com/cloudcore-tu/homebrew-tap.git
brew install snipe-it-cli
```

### deb / rpm（Linux）

[GitHub Releases](https://github.com/cloudcore-tu/snipe-it-cli/releases) からダウンロード。

```bash
# Debian / Ubuntu
sudo dpkg -i snipe-it-cli_<VERSION>_linux_amd64.deb

# RHEL / Fedora
sudo rpm -i snipe-it-cli_<VERSION>_linux_amd64.rpm
```

man ページは `/usr/share/man/` に自動配置される。

### バイナリ（直接ダウンロード）

```bash
VERSION=$(curl -s https://api.github.com/repos/cloudcore-tu/snipe-it-cli/releases/latest | grep tag_name | cut -d'"' -f4 | tr -d v)
curl -L "https://github.com/cloudcore-tu/snipe-it-cli/releases/download/v${VERSION}/snip_${VERSION}_linux_amd64.tar.gz" | tar xz
mv snip /usr/local/bin/
```

## セットアップ

1. Snipe-IT → Settings > API Keys でトークンを発行する
2. 設定ファイルを初期化する:

```bash
snip config init --url https://snip.example.com --token YOUR_API_TOKEN
```

環境変数での指定（設定ファイルより優先）:

```bash
export SNIPEIT_URL=https://snip.example.com
export SNIPEIT_TOKEN=YOUR_API_TOKEN
```

複数インスタンス:

```bash
snip config add staging --url https://staging.example.com --token STAGING_TOKEN
snip config list
export SNIPE_PROFILE=staging          # セッション単位で切り替え
snip --profile prod assets list       # コマンド単位で切り替え
```

## コマンド体系

```text
snip [global flags] {resource} {verb} [flags]
```

Snipe-IT REST API v1 の全エンドポイントに対応。約 150 のサブコマンドを 5 カテゴリに分類。

---

### 1. CRUD リソース（19 種）

標準動詞 `list` / `get` / `create` / `update` / `delete` をすべてのリソースで提供。

| リソース | 対象 | API パス |
|---------|------|---------|
| `assets` | IT 資産（ハードウェア） | `/api/v1/hardware` |
| `users` | ユーザー | `/api/v1/users` |
| `licenses` | ソフトウェアライセンス | `/api/v1/licenses` |
| `accessories` | アクセサリー | `/api/v1/accessories` |
| `consumables` | 消耗品 | `/api/v1/consumables` |
| `components` | コンポーネント | `/api/v1/components` |
| `categories` | カテゴリ | `/api/v1/categories` |
| `companies` | 会社 | `/api/v1/companies` |
| `locations` | ロケーション | `/api/v1/locations` |
| `manufacturers` | メーカー | `/api/v1/manufacturers` |
| `models` | 機器モデル | `/api/v1/models` |
| `departments` | 部門 | `/api/v1/departments` |
| `statuslabels` | ステータスラベル | `/api/v1/statuslabels` |
| `suppliers` | サプライヤー | `/api/v1/suppliers` |
| `fieldsets` | カスタムフィールドセット | `/api/v1/fieldsets` |
| `maintenances` | メンテナンス記録 | `/api/v1/maintenances` |
| `fields` | カスタムフィールド | `/api/v1/fields` |
| `depreciations` | 償却設定 | `/api/v1/depreciations` |
| `groups` | 権限グループ | `/api/v1/groups` |

**動詞フラグ:**

| 動詞 | フラグ |
|------|--------|
| `list` | `--filter key=value`（複数可） `--limit N`（デフォルト 50） `--offset N` |
| `get` | `--id N` |
| `create` | `--data JSON` |
| `update` | `--id N --data JSON`（PATCH） |
| `delete` | `--id N --yes` |

```bash
snip assets list --filter status_id=2 --limit 100
snip assets get --id 123
snip assets create --data '{"name":"Laptop-001","asset_tag":"ASSET-001","model_id":1,"status_id":2}'
snip assets update --id 123 --data '{"status_id":3}'
snip assets delete --id 123 --yes
```

---

### 2. アクション

リソース固有の操作。`POST /api/v1/{resource}/{id}/{action}` を呼ぶ。

| リソース | アクション |
|---------|-----------|
| `assets` | checkout, checkin, audit, restore |
| `licenses` | checkout, checkin |
| `accessories` | checkout, checkin |
| `consumables` | checkout |
| `components` | checkout, checkin |
| `manufacturers` | restore |
| `models` | restore |

```bash
snip assets checkout --id 123 --data '{"checkout_to_type":"user","assigned_user":1}'
snip assets checkin  --id 123
snip assets audit    --id 123
snip assets restore  --id 123
snip licenses checkout --id 10 --data '{"assigned_to":5}'
snip accessories checkin --id 7 --data '{"note":"返却済み"}'
```

---

### 3. サブリソース参照

リソースに紐づく関連データを取得する。`--id N` で親 ID を指定。

| コマンド | 取得内容 |
|---------|---------|
| `assets history --id N` | 資産の操作履歴 |
| `assets licenses --id N` | 資産に割り当てられたライセンス |
| `assets assigned-assets --id N` | 資産に割り当てられた資産 |
| `assets assigned-accessories --id N` | 資産に割り当てられたアクセサリー |
| `assets assigned-components --id N` | 資産に割り当てられたコンポーネント |
| `assets bytag --tag TAG` | 資産タグで資産を検索 |
| `assets byserial --serial SERIAL` | シリアル番号で資産を検索 |
| `users assets --id N` | ユーザーに割り当てられた資産 |
| `users licenses --id N` | ユーザーに割り当てられたライセンス |
| `users accessories --id N` | ユーザーに割り当てられたアクセサリー |
| `users consumables --id N` | ユーザーに割り当てられた消耗品 |
| `licenses history --id N` | ライセンスの操作履歴 |
| `licenses seats list --id N` | ライセンスシート一覧 |
| `licenses seats get --id N --seat-id M` | ライセンスシート詳細 |
| `licenses seats update --id N --seat-id M --data JSON` | ライセンスシート更新 |
| `accessories history --id N` | アクセサリーの操作履歴 |
| `accessories checkedout --id N` | アクセサリーの貸し出し一覧 |
| `consumables history --id N` | 消耗品の操作履歴 |
| `consumables users --id N` | 消耗品を受け取ったユーザー |
| `components history --id N` | コンポーネントの操作履歴 |
| `components assets --id N` | コンポーネントに割り当てられた資産 |
| `locations users --id N` | ロケーションのユーザー |
| `locations assets --id N` | ロケーションの資産 |
| `locations assigned-assets --id N` | ロケーションに割り当てられた資産 |
| `locations assigned-accessories --id N` | ロケーションに割り当てられたアクセサリー |
| `locations history --id N` | ロケーションの操作履歴 |
| `statuslabels assetlist --id N` | ステータスラベルの資産一覧 |
| `statuslabels counts-by-label` | ラベルごとの資産数 |
| `statuslabels counts-by-type` | タイプごとの資産数 |
| `fieldsets fields --id N` | フィールドセットに属するフィールド |
| `maintenances history --id N` | メンテナンス記録の履歴 |
| `models history --id N` | モデルの操作履歴 |

---

### 4. レポート・アカウント

| コマンド | 説明 |
|---------|------|
| `reports activity` | 操作履歴レポート |
| `reports depreciation` | 償却レポート |
| `account requestable` | リクエスト可能な資産一覧 |
| `account requests` | 自分のリクエスト一覧 |
| `account request --id N` | 資産をリクエスト |
| `account cancel-request --id N` | リクエストをキャンセル |
| `account eulas` | EULA 一覧 |
| `account tokens` | 個人アクセストークン一覧 |
| `account token-create --data JSON` | 個人アクセストークン作成 |
| `account token-delete --token-id N` | 個人アクセストークン削除 |

---

### 5. 管理操作

#### ラベル

```bash
snip labels list                                              # ラベルテンプレート一覧
snip labels get --name barcode --output-file label.pdf       # PDF ダウンロード
```

#### インポート

```bash
snip imports list
snip imports create --file assets.csv --type hardware        # CSV アップロード（multipart）
snip imports process --id 3                                  # インポート実行
snip imports delete --id 3 --yes
```

#### 設定

```bash
snip settings get                                            # サーバー設定取得
snip settings update --data '{"per_page":50}'               # サーバー設定更新
snip settings login-attempts                                 # ログイン試行履歴
snip settings backups                                        # バックアップ一覧
snip settings backup-download --output-file backup.zip       # 最新バックアップ DL
snip settings backup-download --name 2026-04-14.zip --output-file backup.zip
```

#### ノート

```bash
snip notes list   --asset-id 123
snip notes create --asset-id 123 --data '{"note":"バッテリー交換済み"}'
```

#### カスタムフィールド追加操作

```bash
snip fields associate    --id 5 --fieldset-id 2
snip fields disassociate --id 5 --fieldset-id 2
snip fields reorder      --fieldset-id 2 --data '[3,1,2]'
```

---

## 出力フォーマット

`-o` / `--output` で指定。デフォルトは table（人間可読）。

```bash
snip assets list -o json
snip assets list -o yaml
snip assets list -o 'custom-columns=ID:.id,NAME:.name,STATUS:.status_label.name'
snip assets list -o 'jsonpath={rows.#.id}'
```

エージェント・スクリプト向けは `-o json` を明示指定する:

```bash
snip assets list -o json | jq '.rows[].name'
snip assets list -o 'jsonpath={rows.#.name}'
```

## グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--url URL` | Snipe-IT URL（環境変数 `SNIPEIT_URL`） |
| `--token TOKEN` | API トークン（環境変数 `SNIPEIT_TOKEN`） |
| `--profile NAME` | 使用するインスタンス名（環境変数 `SNIPE_PROFILE`） |
| `--timeout N` | リクエストタイムアウト秒数（環境変数 `SNIPEIT_TIMEOUT`） |
| `-o FORMAT` | 出力フォーマット（環境変数 `SNIPEIT_OUTPUT`） |
| `-v, --verbose` | API リクエスト概要ログを表示（INFO） |
| `--debug` | HTTP リクエスト/レスポンス詳細ログを表示（DEBUG） |

優先順位: `CLI フラグ > 環境変数 > 設定ファイル > デフォルト`

## シェル補完

```bash
# bash
snip completion bash > /etc/bash_completion.d/snip

# zsh（~/.zshrc に追加）
source <(snip completion zsh)

# fish
snip completion fish > ~/.config/fish/completions/snip.fish
```

## ドキュメント

- [docs/adr/](docs/adr/) — 設計判断の記録
- [docs/api-coverage.md](docs/api-coverage.md) — API カバレッジ一覧
- [CHANGELOG.md](CHANGELOG.md) — 変更履歴

```bash
snip --help
snip assets --help
snip assets list --help
man snip
```

## トラブルシューティング

| 症状 | 対処 |
|------|------|
| `Error: API error: 401 Unauthorized` | `--token` または `SNIPEIT_TOKEN` を確認 |
| `Error: API error: 404 Not Found` | `--id` の値を確認 |
| タイムアウト | `--timeout 60` で延長 |
| 詳細を確認したい | `--verbose` または `--debug` を追加 |
| 設定ファイルのパーミッション警告 | `chmod 0600 ~/.config/snipe-it-cli/config.yaml` |
