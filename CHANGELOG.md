# Changelog

このファイルはすべての変更履歴を記録します。
フォーマットは [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) に準拠し、
バージョン管理は [Semantic Versioning](https://semver.org/) に従います。

## [Unreleased]

## [1.1.0] - 2026-04-16

### Changed

- BaseOptions フィールドを unexported 化（client, printFlags, out）。Complete() 前のアクセスを型レベルで防止
- run package をファイル分割: validate.go, json.go, httpdispatch.go, filedownload.go を新設
- HTTP dispatch helper を本質的な命名に変更（Run*ByPath → *AndPrint 系）
- JSON helper を本質的な命名に変更（UnmarshalJSON → ParseJSONObject, ValidateJSON → CheckJSONSyntax 等）
- client.go 内部の命名改善（requestOptions → apiRequestSpec, doAPIRequest → sendAPIRequest）
- ファイル名を責務に合わせてリネーム（helpers.go → httpdispatch.go, binary.go → filedownload.go）
- 不足していた docstring を追記（extractPayload の失敗契約、sendAPIRequest のフロー等）
- command flow と HTTP 契約を整理し、境界と責務を全体で引き締めるリファクタリング
- path segment 連結を JoinPathSegments に統一し、path injection を防止
- release workflow の changelog 抽出を見出し込みに修正し、未検出時は fail するようにした
- release workflow の release notes からバージョン見出しを除外し、`### Added` / `### Changed` 以降だけ載せるようにした

## [1.0.2] - 2026-04-15

### Changed

- Homebrew tap 更新 workflow を GitHub App 依存から `HOMEBREW_TAP_TOKEN` secret ベースに変更
- 既存 release tag を後追い同期できる manual Homebrew formula sync workflow を追加

## [1.0.1] - 2026-04-15

### Added

- GitHub Release asset に shell completion（bash/zsh/fish）を追加
- deb/rpm パッケージに shell completion を同梱

### Changed

- GitHub Release asset の man page 名を `snip.en.1` / `snip.ja.1` に分離
- Homebrew tap 更新 workflow の formula 更新スクリプト呼び先を `snipe-it-cli` 用に修正

## [1.0.0] - 2026-04-15

### Added

- ローカル Snipe-IT を clean state から起動し、setup・token 作成・smoke・cleanup までを一撃で行う `scripts/snipeit-local-e2e.sh`
- Docker 上で `snip` 自体を実行する local E2E 検証フロー

### Changed

- ローカル検証手順を host 実行前提から dockerized `snip` smoke 前提に更新
- 0.x 系の試作段階を終え、CLI を `1.0.0` としてリリース可能な状態に整理

## [0.3.0] - 2026-04-14

### Added

- labels コマンド: list（JSON）+ get（PDF バイナリ → --output-file 保存）
- imports コマンド: CRUD + create（multipart CSV アップロード）+ process
- settings コマンド: get/update/login-attempts/backups/backup-download
- notes コマンド: list/create（--asset-id で資産に紐づくノート）
- fields コマンド拡張: associate/disassociate/reorder
- account コマンド拡張: eulas/tokens/token-create/token-delete
- クライアント追加: DeleteByPath/Upload（multipart）
- run ヘルパー追加: RunDeleteByPath/RunUpload/RunSaveBinary

## [0.2.0] - 2026-04-14

### Added

- Snipe-IT API 全エンドポイント対応（routes/api.php ベース）
- サブリソース参照コマンド群（history/checkedout/seats/bytag/byserial 等、30+ コマンド）
- 新リソース: fields（カスタムフィールド）、depreciations（償却設定）、groups（権限グループ）
- reports コマンド: activity/depreciation レポート
- account コマンド: requestable/requests/request/cancel-request
- assets: restore/bytag/byserial/assigned-* コマンド追加
- licenses: seats サブコマンドグループ（list/get/update）
- consumables: checkout アクション追加（routes/api.php にあり未実装だったため）
- manufacturers/models: restore アクション追加
- API クライアントに GetSub/GetByPath/PatchByPath/PostByPath メソッド追加
- BuildSubReadCmd/BuildPathReadCmd/RunGetByPath/RunPatchByPath/RunPostByPath ヘルパー追加

### Fixed

- delete コマンドの出力を PrintFlags 経由に統一（他動詞と形式不整合の解消）
- list コマンドの --offset 負値バリデーション追加

## [0.1.0] - 2026-04-14

### Added

- プロジェクト初期セットアップ
- Snipe-IT REST API（/api/v1）への直接 HTTP クライアント実装
- ResourceDef による汎用 CRUD フレームワーク（list/get/create/update/delete）
- 全 16 リソースのコマンド: assets, users, licenses, accessories, components, consumables, categories, companies, locations, manufacturers, models, departments, statuslabels, suppliers, fieldsets, maintenances
- assets/licenses/accessories/components の checkout/checkin アクション
- 出力フォーマット: table（デフォルト）、json、yaml、custom-columns、jsonpath
- 設定管理: config init/add/list コマンド、XDG 準拠、複数インスタンス対応
- 設定の優先順位: CLI フラグ > 環境変数（SNIPEIT_URL/TOKEN/TIMEOUT/OUTPUT）> 設定ファイル
- シェル補完: bash/zsh/fish/powershell（`snip completion <shell>`）
- man ページ: 英語（man/en/snip.1）・日本語（man/ja/snip.1）
- goreleaser によるクロスプラットフォームバイナリ生成（linux/darwin/windows、amd64/arm64）
- deb/rpm パッケージ生成（nfpm、man ページ含む）
- `--verbose/-v` フラグ: INFO ログを表示（API リクエスト概要）
- `--debug` フラグ: DEBUG ログを表示（HTTP リクエスト/レスポンス詳細）
- HTTP デバッグログで Authorization ヘッダーをマスク（`Bearer ***`）
- 設定ファイルのパーミッション 0600 未満時に警告（セキュリティ）
- 包括的テスト: 全コア機能に httptest.Server 統合テスト + race detector

### Changed

- コマンド名を `snip`（snipe の先頭4文字）に決定
- golangci-lint v2 対応（gofmt を formatters セクションへ移動）
