---
name: contextual-commit
description: Contextual Commits 形式でコミットするスキル。コードを変更した後の git commit 時に使用する。「コミットして」「コミット作って」「commit」など、git commit に関連する明示的な指示で発動する。通常の Conventional Commits よりも豊富なコンテキスト（なぜ・何を・影響）をコミットメッセージに埋め込み、AI や将来の開発者がコミット履歴からコンテキストを回収できるようにする。
---

# Contextual Commits

コミットメッセージに「なぜその変更をしたか」のコンテキストを埋め込む。
AI がコミット履歴から文脈を回収できるようにするため、通常の Conventional Commits よりも情報量を多くする。

## コミットメッセージのフォーマット

```text
type(scope): 簡潔な説明（体言止め、日本語50文字以内）

## なぜ
この変更が必要な背景・理由。ビジネス上の動機や技術的な制約を書く。
「何を変えたか」ではなく「なぜ変える必要があったか」に集中する。

## 何を
具体的な変更内容。箇条書きで主要な変更点を列挙する。
diff から読み取れない意図や判断を書く。diff で自明なことは省略してよい。

## 影響
他モジュールへの影響、破壊的変更、マイグレーションの必要性など。

Refs: #issue-number

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

セクション見出しに `##` を使うことで、正規表現 `^## (なぜ|何を|影響)` で各セクションを機械的に抽出できる。

### type 一覧

| type | 用途 |
|------|------|
| feat | 新機能 |
| fix | バグ修正 |
| docs | ドキュメントのみの変更 |
| refactor | リファクタリング（機能変更なし） |
| test | テストの追加・修正 |
| chore | ビルド、CI、依存関係等の雑務 |
| ci | CI/CD 設定の変更 |

### scope の決め方

変更の主なドメインまたはモジュール名を使う（例: `ipam`, `dcim`, `config`, `output`）。

- 特定できない場合や横断的な変更は scope を省略してよい
- 複数 scope にまたがる場合は最も影響の大きいものを1つ選ぶ

### 省略ルール

- subject 行だけで十分に意図が伝わる場合（typo 修正、依存更新等）はボディを省略してよい
- 「何を」は diff から自明な場合は省略してよい
- 「影響」は影響がなければセクションごと省略する
- 「なぜ」は原則書く。ただし subject 行に理由が含まれている場合は省略可
- `Refs:` は関連 Issue がある場合のみ
- `Co-Authored-By:` は AI が関与したコミットに必ずつける

## ワークフロー

### 1. 変更内容の分析

```bash
git status
git diff --staged
git diff
```

### 2. 論理的な単位に分割

1つのコミットは1つの論理的な変更単位にする。

分割の基準:
- ドメインモデルの追加とそのテスト → 同じコミット
- リファクタリングと機能追加 → 別コミット
- フォーマッター/リンター適用のみ → 別コミット（type: chore）

### 3. コミットメッセージの作成

HEREDOC で渡す。

```bash
git commit -m "$(cat <<'EOF'
feat(ipam): prefix 一覧取得コマンドを追加

## なぜ
運用作業でよく使う prefix 一覧表示を CLI から即座に行えるようにするため。
Web UI を開かずにターミナルで確認できることで作業効率が向上する。

## 何を
- `netbox-cli ipam prefix list` コマンドを追加
- フィルターオプション: --site, --tenant, --status
- 出力形式: table（デフォルト）/ json / yaml

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
EOF
)"
```
