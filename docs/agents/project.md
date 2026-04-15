# 共通エージェントガイド

このファイルは Claude Code と Codex の共通エージェントガイドの正本。
`.claude/CLAUDE.md` や `AGENTS.md` のような tool 固有ファイルは薄く保ち、ここを参照する。

## エージェント運用

このファイルをエージェント向けプロジェクト文脈の正本とする。
Codex 用の `AGENTS.md` はこのファイルへのシンボリックリンクとして扱う。

### 応答圧縮ルール

トークン消費抑制のため、通常の応答・進捗報告・レビュー結果は「原始人」寄りの簡潔な文体を使う。技術的中身は削らず、装飾だけ削る。

- 敬語・挨拶・前置き・クッション言葉を省く
- ぼかしを避ける。言い切れることだけ言う
- 自明な補足や背景説明を足さない。聞かれたことにだけ答える
- 箇条書き優先。1項目1意味。重複説明しない
- 助詞や接続を少し崩してもよいが、意味が落ちるなら崩さない
- コード、コマンド、識別子、エラー文、API仕様は絶対に圧縮しない
- 高リスク事項（破壊的操作、セキュリティ、認証、データ損失、重大な設計判断）は簡潔さより明確さを優先する

## プロジェクト概要

- プロジェクト名: `snipe-it-cli`
- バイナリ名: `snip`
- 目的: terminal、script、coding agent から Snipe-IT を操作する Go 製 CLI
- 主な利用者: まず人間の運用者、その次に agent

## プロダクトの形

- デフォルト出力は人間向けの table。
- agent や script が機械可読出力を必要とする場合は明示的に `-o json` を付ける。
- コマンド体系は `snip [global flags] {resource} {verb} [flags]`。
- 基本 CRUD resource は `cmd/internal/run.ResourceDef` で実装する。

## コード設計ルール

- resource package は宣言的に保つ。再利用ロジックは各 `cmd/<resource>` ではなく `cmd/internal/run` に置く。
- inline `RunE` でも `Complete -> Validate -> Run` の形を優先する。
- API call 前に ID と必須文字列を検証する。
- JSON input は中央で検証し、`json.Marshal` や共通 helper で済む場面では手書き JSON を避ける。
- JSON response 出力は `run.BaseOptions` を通して一貫させる。
- `ResourceDef` は one-off 機能で膨らませず、明示 helper で済むならその形を取る。

## 設計ガードレール

- カプセル化を優先する。可変状態や命令的手順は、小さな helper や options struct の内側に閉じ込める。
- 関心を分離する。command wiring、validation、API transport、rendering、永続化、自動化を同じ層に混ぜない。
- 契約による設計を徹底する。入力、出力、不変条件、失敗条件を明示し、境界で前提条件を検証する。
- 副作用を隔離する。ファイル書き込み、環境変数依存、network call、process 全体への変更は監査しやすい狭い境界に閉じ込める。
- テストにも同じ原則を適用する。重複した setup、assertion、fixture、副作用は helper に集約し、ケース固有の意図だけを各 test に残す。
- 変更が不要な場合は no-op で終える。隠れた mutation を起こさない。

## 現在の共通抽象

- `cmd/internal/run/run.go`
  - config 読み込み
  - 共通 validation helper
  - JSON validation / marshalling helper
  - response rendering helper
- `cmd/internal/run/resource.go`
  - 汎用 CRUD command 生成
  - sub-resource / arbitrary-path 共通 helper
- `internal/snipeit`
  - HTTP client と API error handling
- `internal/output`
  - output format の解釈と rendering

## 保守の優先順位

- cleverness より可読性
- 局所的な都合より一貫性
- 短いコードは、人間にとって明白なままのときだけ良い
- 小さな private helper は、重複除去と責務の明確化に効くなら許容する。単なる式の言い換えや無意味な制御フロー隠蔽は避ける

## 応答スタイル規則

- 技術的中身は残し、装飾だけ削る
- 敬語・挨拶・前置き・クッション言葉を省く
- ぼかしを避ける。必要なときだけ補足する
- まず結論を答える。追加説明は正確性・安全性・実行性を上げるときだけ足す
- 簡潔さで意味が落ちるなら、短さより明確さを優先する

## テスト規則

- 特に non-CRUD command では validation path の focused test を足す。
- `cmd/internal/run` に新しい shared helper の挙動を入れたら test を追加・更新する。
- client layer まで到達させたい command path は local HTTP test server を優先する。
- 制限環境では Go build cache 権限で `go test ./...` が落ちることがある。必要なら書き込み可能な `GOCACHE` を使う。

## ドキュメント規則

- user 向けの機能概要は `README.md` と `docs/api-coverage.md` に置く。
- agent 向けの再利用ガイドは `docs/agents/` に置く。
- 再利用 workflow は `docs/agents/skills/` に 1 回だけ書き、tool 固有の skill wrapper から参照させる。

## 便利コマンド

```bash
go test ./...
gofmt -w ./...
golangci-lint run
bash scripts/snipeit-local-e2e.sh
```
