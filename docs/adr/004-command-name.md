# ADR-004: CLI コマンド名を `snip` に決定

- **日付**: 2026-04-14
- **ステータス**: 採用済み

## コンテキスト

CLI ツールのコマンド名を決める必要があった。手でシェルに打つ頻度が高いため、短さが重要。

## 検討した選択肢

| 候補 | 却下理由 |
|------|---------|
| `snipe` | Snipe-IT との紐付きが弱い（"snipe" 単体では別の意味）|
| `sit` | 英語として問題のある印象を与える可能性がある |
| `snipeit` | 7文字、手入力には長い |
| `snip` | **採用** |

## 決定

**`snip`** を採用する。

- 4文字。`nbox`（netbox-cli）と同じ長さ感
- "snipe" の先頭から自然に取れる
- 英語として問題なし（"to snip" = 切る、自然な単語）
- 現環境（Linux/WSL2）で競合コマンドなし（`which snip` = not found）
- Homebrew に `snip` cask（GUI スクリーンキャプチャ）が存在するが、GUI アプリのため PATH に入らず競合しない
- Linux 標準コマンド（coreutils / util-linux / bash builtins）に `snip` は存在しない

## 結果

```text
snip assets list
snip assets get --id 123
snip config init --url https://snipeit.example.com --token xxx
```
