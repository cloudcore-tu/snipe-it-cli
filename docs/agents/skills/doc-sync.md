# doc-sync

コード変更に対してドキュメント更新が必要か確認する。

## 確認対象

- `docs/agents/project.md`
- `.claude/CLAUDE.md`
- `AGENTS.md`
- `docs/adr/`
- `CHANGELOG.md`
- `README.md`
- `docs/api-coverage.md`

## 手順

1. diff を確認する。
2. 変更を分類する: dependency、configuration、architecture、command surface、user-visible behavior。
3. 対応する doc が更新されているか確認する。

## ルール

- shared agent guidance はまず `docs/agents/` を更新する。
- `.claude/CLAUDE.md`、`AGENTS.md`、tool 固有 skill wrapper は薄く保つ。
- command surface が変わる場合は `README.md` と `docs/api-coverage.md` の両方を確認する。
