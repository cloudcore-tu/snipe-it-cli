# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- コマンド名を `snipe` から `sit`（Snipe-IT の頭字語）に変更

### Added
- プロジェクト初期セットアップ
- Snipe-IT REST API（/api/v1）への直接 HTTP クライアント実装
- ResourceDef による汎用 CRUD フレームワーク（list/get/create/update/delete）
- 全 16 リソースのコマンド: assets, users, licenses, accessories, components, consumables, categories, companies, locations, manufacturers, models, departments, statuslabels, suppliers, fieldsets, maintenances
- assets/licenses/accessories/components の checkout/checkin アクション
- 出力フォーマット: table（デフォルト）、json、yaml、custom-columns、jsonpath
- 設定管理: config init/add/list コマンド、XDG 準拠、複数インスタンス対応
- 設定の優先順位: CLI フラグ > 環境変数（SNIPEIT_URL/TOKEN/TIMEOUT/OUTPUT）> 設定ファイル
