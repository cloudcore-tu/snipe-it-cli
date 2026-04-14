# Changelog

このファイルはすべての変更履歴を記録します。
フォーマットは [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) に準拠し、
バージョン管理は [Semantic Versioning](https://semver.org/) に従います。

## [Unreleased]

### Added

- プロジェクト初期セットアップ
- Snipe-IT REST API（/api/v1）への直接 HTTP クライアント実装
- ResourceDef による汎用 CRUD フレームワーク（list/get/create/update/delete）
- 全 16 リソースのコマンド: assets, users, licenses, accessories, components, consumables, categories, companies, locations, manufacturers, models, departments, statuslabels, suppliers, fieldsets, maintenances
- assets/licenses/accessories/components の checkout/checkin アクション
- 出力フォーマット: table（デフォルト）、json、yaml、custom-columns、jsonpath
- 設定管理: config init/add/list コマンド、XDG 準拠、複数インスタンス対応
- 設定の優先順位: CLI フラグ > 環境変数（SNIPEIT_URL/TOKEN/TIMEOUT/OUTPUT）> 設定ファイル

### Changed

- コマンド名を `snip`（snipe の先頭4文字）に決定
