# Shared Agent Guide

This file is the canonical agent-facing guide for both Claude Code and Codex.
Tool-specific files such as `.claude/CLAUDE.md` and `AGENTS.md` should stay thin and point here.

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

## Project Summary

- Project: `snipe-it-cli`
- Binary name: `snip`
- Purpose: a Go CLI for operating Snipe-IT from terminals, scripts, and coding agents
- Primary users: human operators first, agents second

## Product Shape

- Default output is human-readable table output.
- Agents and scripts should explicitly request `-o json` when machine-readable output matters.
- Commands follow `snip [global flags] {resource} {verb} [flags]`.
- Core CRUD resources are implemented through `cmd/internal/run.ResourceDef`.

## Code Design Rules

- Keep resource packages declarative. Put reusable mechanics in `cmd/internal/run`, not in each `cmd/<resource>` package.
- Prefer `Complete -> Validate -> Run` structure even when using inline `RunE`.
- Validate IDs and required string values before making API calls.
- Validate JSON input centrally and avoid hand-written JSON strings when `json.Marshal` or shared helpers suffice.
- Keep output behavior consistent by routing JSON response printing through `run.BaseOptions`.
- Avoid expanding `ResourceDef` with one-off features when explicit helper commands keep the framework simpler.

## Design Guardrails

- Prefer encapsulation. Hide mutable state and imperative steps behind small helpers or option structs instead of leaking them across packages.
- Separate concerns strictly. Keep command wiring, validation, API transport, rendering, release automation, and local-dev scripts in distinct layers.
- Design by contract. Make inputs, outputs, and failure conditions explicit; validate preconditions early and fail loudly when contracts are not met.
- Isolate side effects. Keep filesystem writes, environment-variable dependence, network calls, release publication, and tap updates in narrow, easy-to-audit boundaries.
- Treat workflows and shell scripts as production code. Apply the same rules of encapsulation, separation of concerns, contracts, and side-effect isolation there too.
- Prefer no-op behavior over hidden mutation. If a step has nothing to change, exit cleanly instead of rewriting state anyway.

## Release Notes Contract

- GitHub Release notes are derived from `CHANGELOG.md`.
- Only the body of the matching version section is used.
- Exclude the `## [x.y.z] - YYYY-MM-DD` heading itself.
- Keep subsection headings such as `### Added` and `### Changed`.
- Fail the workflow if the matching section is missing or renders empty.

## Current Shared Abstractions

- `cmd/internal/run/run.go`
  - config loading
  - common validation helpers
  - JSON validation/marshalling helpers
  - response rendering helpers
- `cmd/internal/run/resource.go`
  - generic CRUD command generation
  - shared sub-resource and arbitrary-path helpers
- `internal/snipeit`
  - HTTP client and API error handling
- `internal/output`
  - output format parsing and rendering

## Maintenance Priorities

- Readability over cleverness.
- Consistency over local convenience.
- Compact code is good only when it remains obvious to a human reader.
- Small private helper functions are fine when they remove duplication and carry a clear responsibility. Avoid private helpers that only rename one expression or hide control flow for no gain.

## Response Style Rules

- 技術的中身は残し、装飾だけ削る
- 敬語・挨拶・前置き・クッション言葉を省く
- ぼかしを避ける。必要なときだけ補足する
- まず結論を答える。追加説明は正確性・安全性・実行性を上げるときだけ足す
- 簡潔さで意味が落ちるなら、短さより明確さを優先する

## Testing Rules

- Add focused tests for command validation paths, especially for non-CRUD commands.
- Add or update tests when introducing new shared helper behavior in `cmd/internal/run`.
- Prefer local HTTP test servers for command paths that should reach the client layer.
- In restricted environments, `go test ./...` may fail because of Go build cache permissions; use a writable `GOCACHE` if needed.

## Documentation Rules

- User-facing capability summaries live in `README.md` and `docs/api-coverage.md`.
- Agent-facing reusable guidance lives in `docs/agents/`.
- Reusable workflows should be authored once under `docs/agents/skills/` and referenced from tool-specific skill wrappers.

## Handy Commands

```bash
go test ./...
gofmt -w ./...
golangci-lint run
bash scripts/snipeit-local-e2e.sh
```
