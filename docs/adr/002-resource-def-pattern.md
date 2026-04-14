# ADR-002: ResourceDef による汎用 CRUD フレームワーク

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

Snipe-IT の REST API は 16 以上のリソース（hardware, users, licenses, categories, locations 等）を提供し、すべて同一の CRUD パターン（list/get/create/update/delete）に従う。

各リソースに個別の cobra コマンド実装を書くと、コードの重複が膨大になり変更コストが高い。

netbox-cli の `ResourceDef` パターンが実績ある解として存在している。

## 決定

`cmd/internal/run/resource.go` に `ResourceDef` 構造体を定義し、リソースごとに `APIPath` を宣言するだけで list / get / create / update / delete のサブコマンドを自動生成する。

```go
type ResourceDef struct {
    Use     string
    Short   string
    DocsURL string
    APIPath string       // "hardware", "users", "licenses" 等
    ActionFns []ActionDef
}
```

netbox-cli との主な違い:

| 項目 | netbox-cli | snipe-it-cli |
|------|------|------|
| バックエンド | go-netbox（型付き SDK） | 直接 HTTP（汎用） |
| ResourceDef の関数フィールド | `ListFn`, `GetFn`, `CreateFn` 等 | なし（APIPath から自動） |
| バルク操作 | `BulkCreate/BulkDelete/BulkUpdate` | 未実装（必要時に追加） |

Snipe-IT API では型付き SDK なしに汎用 HTTP メソッドを呼ぶため、関数フィールドは不要。`APIPath` だけで全 CRUD を生成できる。

### カスタム操作

list/get/create/update/delete 以外の操作（checkout, checkin 等）は `ActionDef` で定義する。

```go
type ActionDef struct {
    Use   string
    Short string
    NeedsData bool  // --data フラグを受け付けるか
    Action    string // API アクションパス（"checkout", "checkin" 等）
}
```

`ActionDef` の実体は `POST /api/v1/{resource}/{id}/{action}` を呼ぶ。

## 採用した理由

- 各リソースのコード量が大幅に削減できる（1 リソースあたり ~3 行）
- フレームワーク側の改善が全リソースに即時反映される
- Snipe-IT の一貫した API パターンと相性が良い

## トレードオフ

- コマンド生成の詳細が `resource.go` に集中するため、個別リソースの挙動トレース時に間接参照が増える
- リソース固有のフラグ（assets の `--status`, `--assigned-to` 等）は `--filter key=value` で対応するため、型安全ではない

## 関連

- [ADR-001](001-no-go-sdk-direct-http.md): Go SDK を使わず直接 HTTP クライアントを実装する
- [ADR-003](003-configuration-design.md): 設定管理の設計
