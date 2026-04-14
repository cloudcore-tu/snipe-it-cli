# ADR-001: Go SDK を使わず直接 HTTP クライアントを実装する

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

snipe-it-cli の Snipe-IT API クライアント実装方針を決定する必要がある。

Snipe-IT の REST API は `/api/v1/{resource}` の一貫したパターンを持ち、全リソースで同一の CRUD エンドポイントを提供する。

調査時点で利用可能な Go クライアントライブラリ:

| ライブラリ | Stars | 最終更新 | 備考 |
|---|---|---|---|
| `michellepellon/go-snipeit` | 0 | 2026-03-08 | 完成度不明 |
| `euracresearch/go-snipeit` | 2 | 2022-05-12 | 一部エンドポイントのみ対応 |

## 検討した選択肢

### A. 既存 Go ライブラリ（michellepellon/go-snipeit）を使用する

- Pros: 実装量を減らせる可能性がある
- Cons: Star 0、内部品質不明、プロジェクト基準（star 1,000 以上）を大幅に下回る、メンテ継続性が不安定

### B. 既存 Go ライブラリ（euracresearch/go-snipeit）を使用する

- Pros: 実績がある（Eurac Research 内部利用）
- Cons: Star 2、最終コミット 2022年（2年以上更新なし）、一部エンドポイントのみ対応

### C. 直接 HTTP クライアントを実装する

- Pros:
  - 外部依存ゼロで安定性が高い
  - Snipe-IT API の一貫したパターン（全リソースが同一 CRUD 構造）により、汎用メソッド数個で全リソースをカバーできる
  - ResourceDef に APIPath を持たせるだけで全リソースの CRUD を自動生成できる
  - テストで httptest.Server を使えばモック不要
- Cons: 初期実装コストがある（ただし汎用メソッドのみで少量）

## 決定

**C. 直接 HTTP クライアントを実装する**

Snipe-IT API の一貫したパターン（`GET /api/v1/{resource}`, `POST /api/v1/{resource}`, etc.）を活かし、以下の汎用メソッドを持つ `internal/snipeit.Client` を実装する:

```go
type Client struct { ... }

func (c *Client) List(ctx, path, params) ([]byte, error)
func (c *Client) GetByID(ctx, path, id) ([]byte, error)
func (c *Client) Create(ctx, path, data) ([]byte, error)   // payload を抽出して返す
func (c *Client) Update(ctx, path, id, data) ([]byte, error)  // PATCH、payload を抽出して返す
func (c *Client) Delete(ctx, path, id) error
func (c *Client) PostAction(ctx, path, id, action, data) ([]byte, error)  // checkout/checkin 等
```

`ResourceDef` は `APIPath` フィールドだけを持てばよく、関数フィールドが不要になる（netbox-cli の `ResourceDef` よりシンプル）。

## 結果

- 外部 Go ライブラリへの依存なし
- `cmd/internal/run/resource.go` が APIPath を使って汎用クライアントメソッドを呼ぶ
- 新リソースの追加は `ResourceDef{Use: "...", APIPath: "..."}` 数行で完了する
- httptest.Server を使ったテストが直接書ける

## 関連

- [ADR-002](002-resource-def-pattern.md): ResourceDef による汎用 CRUD フレームワーク
