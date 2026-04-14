# snipe-it-cli

Snipe-IT（IT資産管理 OSS）を操作する Go 製 CLI ツール。コマンド名: `snip`。対応 Snipe-IT: v8.x

| OS | アーキテクチャ |
|----|--------------|
| macOS | amd64（Intel）/ arm64（Apple Silicon） |
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

API トークンの発行:

1. Snipe-IT にログインし、Settings > API Keys を開く
2. Create New Token をクリックしてトークンを作成する
3. 作成したトークンを以下の `YOUR_API_TOKEN` 部分に指定する

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

## 使い方

```text
snip [global flags] {resource} {verb} [flags]
```

### リソース一覧

| コマンド | 対象 |
|---------|------|
| `assets` | IT 資産（ハードウェア） |
| `users` | ユーザー |
| `licenses` | ソフトウェアライセンス |
| `accessories` | アクセサリー |
| `consumables` | 消耗品 |
| `components` | コンポーネント |
| `categories` | カテゴリ |
| `companies` | 会社 |
| `locations` | ロケーション |
| `manufacturers` | メーカー |
| `models` | 機器モデル |
| `departments` | 部門 |
| `statuslabels` | ステータスラベル |
| `suppliers` | サプライヤー |
| `fieldsets` | カスタムフィールドセット |
| `maintenances` | メンテナンス記録 |

### 操作一覧

| 操作 | フラグ |
|------|--------|
| `list` | `--filter key=value`（複数可）`--limit N`（デフォルト 50）`--offset N` |
| `get` | `--id ID` |
| `create` | `--data JSON` |
| `update` | `--id ID --data JSON`（PATCH） |
| `delete` | `--id ID --yes` |

### 資産固有操作

```bash
snip assets checkout --id 123 --data '{"checkout_to_type":"user","assigned_user":1}'
snip assets checkin  --id 123
```

### 例

```bash
# 資産
snip assets list
snip assets list --filter status_id=2 --limit 100
snip assets get --id 123
snip assets create --data '{"name":"Laptop-001","asset_tag":"ASSET-001","model_id":1,"status_id":2}'
snip assets update --id 123 --data '{"status_id":3}'
snip assets delete --id 123 --yes

# ユーザー
snip users list
snip users get --id 5

# ライセンス
snip licenses list
snip licenses create --data '{"name":"Office 365","category_id":3,"seats":10}'

# エージェント・スクリプト向け（JSON 出力）
snip assets list -o json | jq '.[].name'
snip assets list -o 'jsonpath={rows.#.name}'
```

### 出力フォーマット

`-o` または `--output` で指定。デフォルトは table（人間可読）。

```bash
snip assets list -o json
snip assets list -o yaml
snip assets list -o 'custom-columns=ID:.id,NAME:.name,STATUS:.status_label.name'
snip assets list -o 'jsonpath={rows.#.id}'
```

### グローバルフラグ

| フラグ | 説明 |
|--------|------|
| `--url URL` | Snipe-IT URL（環境変数 `SNIPEIT_URL`） |
| `--token TOKEN` | API トークン（環境変数 `SNIPEIT_TOKEN`） |
| `--profile NAME` | 使用するインスタンス名（環境変数 `SNIPE_PROFILE`） |
| `--timeout N` | リクエストタイムアウト秒数（環境変数 `SNIPEIT_TIMEOUT`） |
| `-o FORMAT` | 出力フォーマット（環境変数 `SNIPEIT_OUTPUT`） |
| `-v, --verbose` | API リクエスト概要ログを表示（INFO） |
| `--debug` | HTTP リクエスト/レスポンス詳細ログを表示（DEBUG）|

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
