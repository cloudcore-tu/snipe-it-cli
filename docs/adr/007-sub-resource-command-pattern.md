# ADR-007: サブリソースコマンドの実装パターン

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

Snipe-IT API には CRUD の他にサブリソース参照エンドポイントが多数存在する。

- `GET /api/v1/hardware/{id}/history`
- `GET /api/v1/licenses/{id}/seats`
- `GET /api/v1/statuslabels/assets/name`（静的パス）

など 30+ パターン。これらをどう CLI コマンドに落とし込むかを決定する必要があった。

## 検討した選択肢

### A. ResourceDef に SubQueryDef フィールドを追加する

```go
type ResourceDef struct {
    // ...
    SubQueries []SubQueryDef
}

type SubQueryDef struct {
    Use     string
    SubPath string
}
```

- Pros: ResourceDef の宣言だけで全コマンドが生成される
- Cons: ResourceDef が肥大化する。静的パス（`statuslabels/assets/name`）や特殊フラグ（`bytag --tag`, `byserial --serial`）をフィールドで表現しようとすると
  さらに複雑になる。Codex コードレビューで却下推奨。

### B. 各パッケージに専用ヘルパー関数を直接書く

ResourceDef は CRUD のみ。サブリソースは各パッケージで `BuildSubReadCmd` / `BuildPathReadCmd` を呼ぶ。

```go
// cmd/assets/assets.go
cmd.AddCommand(run.BuildSubReadCmd("history", "資産の操作履歴", "hardware", "history"))
cmd.AddCommand(run.BuildSubReadCmd("licenses", "資産に割り当てられたライセンス", "hardware", "licenses"))
```

- Pros:
  - ResourceDef がシンプルなまま（CRUD のみ）
  - 特殊なサブコマンド（bytag, byserial, seats など）を各パッケージで自由に実装できる
  - テスト・トレースが容易（明示的な関数呼び出し）
- Cons: 各パッケージで `AddCommand` 呼び出しが増える（ただし 1〜3 行/コマンド）

## 決定

**B を採用。**

`cmd/internal/run/resource.go` に以下のヘルパーを追加し、各パッケージで呼び出す:

```go
// BuildSubReadCmd: GET /api/v1/{parentAPIPath}/{id}/{subPath}
func BuildSubReadCmd(use, short, parentAPIPath, subPath string) *cobra.Command

// BuildPathReadCmd: GET /api/v1/{apiPath}（静的パス。id 不要）
func BuildPathReadCmd(use, short, apiPath string) *cobra.Command

// RunGetByPath: 任意パスへの GET
func RunGetByPath(ctx context.Context, o *BaseOptions, urlPath string) error

// RunPostByPath / RunPatchByPath / RunDeleteByPath / RunSaveBinary / RunUpload
// 特殊な操作パターン（バイナリ出力、ファイルアップロード等）に対応
```

ResourceDef は `APIPath` + `ActionFns` のみ保持し、サブリソース参照はスコープ外とする。

## 採用した理由

- Codex レビューが「ResourceDef に SubQueryDef を追加するとフレームワークが肥大化し、
  例外ケース（静的パス、特殊フラグ）の扱いが複雑になる」と指摘
- `BuildSubReadCmd` は 1 行で `cobra.Command` を返す軽量ヘルパーであり、
  ResourceDef の複雑化なしに再利用性を確保できる
- 特殊ケース（`bytag --tag`, `seats list/get/update`, `notes --asset-id` 等）を
  各パッケージで自由に実装できる柔軟性が重要

## トレードオフ

- 各リソースパッケージのコード量が増える（`assets.go` 等に `AddCommand` 行が増加）
- ただし 1 コマンド = 1 行で読みやすく、ResourceDef の黒魔術的な自動生成より追跡しやすい

## 実装パターン例

```go
// 通常サブリソース
cmd.AddCommand(run.BuildSubReadCmd("history", "操作履歴", "hardware", "history"))

// 静的パス（id 不要）
cmd.AddCommand(run.BuildPathReadCmd("counts-by-label", "ラベルごとの資産数", "statuslabels/assets/name"))

// 特殊フラグ（--tag / --serial）
cmd.AddCommand(buildByTagCmd())   // 各パッケージで独自実装

// ネスト CRUD（licenses seats）
cmd.AddCommand(buildSeatsCmd())   // list/get/update を含むグループ
```

## 関連

- [ADR-002](002-resource-def-pattern.md): ResourceDef による汎用 CRUD フレームワーク
- [ADR-001](001-no-go-sdk-direct-http.md): 直接 HTTP クライアント実装
