# ADR-005: テストヘルパーを cmd/internal/testutil に配置する

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

netbox-cli（ADR-003）で確立した方針と同じ問題が snipe-it-cli にも生じる。
各モジュールのテスト（`cmd/assets/assets_test.go` 等）で `testServer`・`newTestBaseOptions`・`mustMarshalJSON` を共通化したいが、配置場所によって Go の `internal` パッケージ制約に引っかかる。

## 決定

`cmd/internal/testutil` に配置する。

## 理由

テストヘルパーは `cmd/internal/run.BaseOptions` と `internal/snipeit.Client` に依存する。

Go の `internal` パッケージルールでは、`internal` ディレクトリ配下のパッケージはその `internal` の **親ディレクトリ以下** からのみインポートできる。

```text
internal/testutil      ← cmd/ からはインポート不可（違反）
cmd/internal/testutil  ← cmd/ 以下ならどこからでもインポート可能 ✓
```

## トレードオフ

`internal/output` 等の `internal/` 配下パッケージは `cmd/internal/testutil` を使えない（逆方向の依存になるため）。`internal/` のテストは自前でヘルパーを定義する。

## 関連

- netbox-cli ADR-003（同一決定の元となった先行事例）
- [ADR-002](002-resource-def-pattern.md): ResourceDef による汎用 CRUD フレームワーク
