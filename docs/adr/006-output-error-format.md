# ADR-006: 出力・エラーフォーマットの設計

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

netbox-cli（ADR-005）で確立した方針と同じ設計を snipe-it-cli に適用する。
想定ユーザーは人間のインフラエンジニアとコーディングエージェント（Claude Code 等）の両方。

## 決定

### 通常出力

| フォーマット | フラグ | 用途 |
| --- | --- | --- |
| table | デフォルト | 人間が一覧を確認する際 |
| JSON | `-o json` | エージェント・スクリプト向け |
| YAML | `-o yaml` | 設定ファイルへの取り込みや diff |
| custom-columns | `-o custom-columns=ID:.id,NAME:.name` | フィールドを選んで一覧表示 |
| jsonpath | `-o 'jsonpath={rows.#.id}'` | スクリプトでの値抽出（gjson 構文） |

デフォルトを table にする。エージェントは `-o json` を明示指定する。

### エラー出力

人間可読なテキスト形式で stderr に出力し、exit code で成否を表す。

```text
Error: API error: HTTP 404: No asset matches the given query.
```

| exit code | 意味 |
| --- | --- |
| 0 | 成功 |
| 1 | エラー（API エラー、認証失敗、引数の誤り等） |

### 対話的プロンプトの排除

`--yes` フラグで削除の確認を明示させる。エージェントが tty なし環境で実行することが多く、対話プロンプトはブロックの原因になるため。

## 理由

netbox-cli での運用で「デフォルト table・エージェントは `-o json` 明示」が正しいと確認済み。
kubectl 等の成熟した CLI ツールと同じ方針。

## 関連

- netbox-cli ADR-005（同一決定の元となった先行事例）
- [ADR-002](002-resource-def-pattern.md): ResourceDef による汎用 CRUD フレームワーク
