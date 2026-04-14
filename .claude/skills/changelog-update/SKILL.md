---
name: changelog-update
description: Keep a Changelog 形式で CHANGELOG.md を更新するスキル。「CHANGELOG 更新」「変更履歴に追加」「changelog」などの指示で発動する。コミット後やリリース準備時に使用する。
---

# CHANGELOG 更新

[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) 形式で `CHANGELOG.md` を更新する。

## フォーマット

```markdown
## [Unreleased]

### Added
- 新機能の説明（日本語）

### Changed
- 変更内容の説明（日本語）

### Deprecated
- 非推奨になった機能の説明（日本語）

### Removed
- 削除された機能の説明（日本語）

### Fixed
- 修正されたバグの説明（日本語）

### Security
- セキュリティ修正の説明（日本語）
```

ルール:
- セクション見出し（Added, Changed 等）は英語
- 記載内容は日本語
- Semantic Versioning に準拠
- `[Unreleased]` セクションに追記していく

## ワークフロー

### Step 1: 変更内容の特定

```bash
git log --oneline -10
git diff HEAD~1 --stat
```

### Step 2: 既存 CHANGELOG の確認

```bash
cat CHANGELOG.md
```

`[Unreleased]` セクションの現在の内容を確認し、重複を避ける。

### Step 3: カテゴリの判定

| コミット type | CHANGELOG カテゴリ |
|--------------|-------------------|
| feat | Added |
| fix | Fixed |
| refactor | Changed |
| chore（依存関係変更） | Changed |
| chore（削除） | Removed |
| security | Security |

### Step 4: CHANGELOG に追記

`[Unreleased]` セクションの適切なカテゴリに追記する。
カテゴリが存在しなければ新規作成する。

記述スタイル:
- 箇条書き（`-` で開始）
- 機能単位でまとめる（コミット単位ではない）
- ユーザーに見える変更を中心に
