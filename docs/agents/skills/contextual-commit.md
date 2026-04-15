# contextual-commit

ユーザーが commit を明示要求したときに、contextual commit 形式で commit message を作る。

## 形式

```text
type(scope): 簡潔な説明

## なぜ

## 何を

## 影響

Co-Authored-By: ...
```

## ルール

- なぜその変更が必要かを強く出す。
- 1 commit 1 論理変更を保つ。
- `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci` を使う。
- AI が実質的に関与した場合は `Co-Authored-By` を付ける。
