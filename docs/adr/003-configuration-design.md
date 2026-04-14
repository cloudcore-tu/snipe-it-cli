# ADR-003: 設定管理の設計（XDG・複数インスタンス・kubectl 参考）

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

snipe-it-cli を使う場面として以下が想定される。

1. **インフラエンジニアの手動操作** — 開発・ステージング・本番など複数環境を日常的に切り替える
2. **コーディングエージェントの自動操作** — 環境変数や設定ファイルを通じて接続先を注入する

netbox-cli の設定設計（ADR-004）が実績ある解として存在しているため、同じ設計方針を採用する。

## 決定

### 設定ファイルの場所

XDG Base Directory Specification に準拠し `$XDG_CONFIG_HOME/snipe-it-cli/config.yaml`（未設定時は `~/.config/snipe-it-cli/config.yaml`）に置く。

`os.UserConfigDir()` は macOS で `~/Library/Application Support` を返すため使用しない。

### 複数インスタンス管理

kubectl の context パターンを参考に `FileConfig` 構造体を設計する。

```yaml
current: prod
instances:
  prod:
    url: https://snipeit.example.com
    token: prod-api-token
  staging:
    url: https://staging.example.com
    token: stg-api-token
timeout: 30
output: table
```

- `current` フィールドがデフォルトのインスタンス名を保持する
- プロファイル解決順: `--profile フラグ > SNIPE_PROFILE 環境変数 > current`

### 設定管理コマンド

| コマンド | 説明 |
| --- | --- |
| `snipeit config init [--name NAME] --url URL --token TOKEN` | 初期設定ファイルを生成 |
| `snipeit config add NAME --url URL --token TOKEN` | インスタンスを追加・更新 |
| `snipeit config list` | 登録済みインスタンスを一覧表示（* がアクティブ） |

### 設定値の優先順位

`CLI フラグ > 環境変数 > 設定ファイル（選択インスタンス） > デフォルト値`

### 環境変数

| 設定項目 | 環境変数 |
|---------|---------|
| Snipe-IT URL | `SNIPEIT_URL` |
| API トークン | `SNIPEIT_TOKEN` |
| タイムアウト（秒） | `SNIPEIT_TIMEOUT` |
| デフォルト出力形式 | `SNIPEIT_OUTPUT` |
| アクティブなインスタンス | `SNIPE_PROFILE` |

### YAML ライブラリ

`go.yaml.in/yaml/v3`（旧 `gopkg.in/yaml.v3` の後継）を使用。viper は不使用（nested map 構造との相性が悪い）。

## 代替案の却下理由

**`snipeit config use` コマンド（kubectl 方式）**: 設定ファイルをコマンドで書き換える方式は複数ターミナルセッションで危険。AWS CLI 方式（`SNIPE_PROFILE` 環境変数）の方が状態を持たない設計として優れるため却下。

**単一インスタンスのフラット設定**: 後から複数対応を追加すると config 構造の破壊的変更が発生するため却下。

**viper による設定管理**: ネストした `instances` マップを扱うと型アサーションが複雑になる。シンプルな YAML + `os.Getenv` で十分なため採用しない。
