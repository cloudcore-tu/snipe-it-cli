# changelog-update

`CHANGELOG.md` を Keep a Changelog 形式で更新する。

## 手順

1. `git log` と `git diff` で最近の変更を確認する。
2. `CHANGELOG.md` を開いて `[Unreleased]` を探す。
3. 正しい section に項目を追加する。

## ルール

- section 見出しは英語のまま保つ: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`。
- 項目本文は日本語で書く。
- 内部 refactor より user-visible な変更を優先する。
- 同じ機能の重複記載を避ける。
