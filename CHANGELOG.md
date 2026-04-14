# Changelog

このファイルはすべての変更履歴を記録します。
フォーマットは [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) に準拠し、
バージョン管理は [Semantic Versioning](https://semver.org/) に従います。

## [Unreleased]

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
