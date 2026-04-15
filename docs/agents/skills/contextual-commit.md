# contextual-commit

Write commit messages in contextual commit format when the user explicitly asks for a commit.

## Format

```text
type(scope): 簡潔な説明

## なぜ

## 何を

## 影響

Co-Authored-By: ...
```

## Rules

- Emphasize why the change exists.
- Keep one logical change per commit.
- Use `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, or `ci`.
- Add `Co-Authored-By` when AI materially contributed.
