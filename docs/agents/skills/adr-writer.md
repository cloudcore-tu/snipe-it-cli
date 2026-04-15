# adr-writer

Create a new ADR under `docs/adr/` when the task introduces or records an architectural decision.

## Workflow

1. List existing ADRs in `docs/adr/` and pick the next number.
2. Create `NNN-kebab-case-title.md`.
3. Use this template:

```markdown
# ADR-NNN: Title

- **ステータス**: 提案 | 承認 | 却下 | 廃止
- **日付**: YYYY-MM-DD

## コンテキスト

## 検討した選択肢

## 決定

## 結果

## 未決事項
```

## Rules

- Default status is `提案`.
- Write factual context, not speculation.
- Record why the chosen option won, not only what was chosen.
