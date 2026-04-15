# Shared Agent Guide

This file is the canonical agent-facing guide for both Claude Code and Codex.
Tool-specific files such as `.claude/CLAUDE.md` and `AGENTS.md` should stay thin and point here.

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
```
