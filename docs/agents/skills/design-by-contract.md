# design-by-contract

契約による設計をコード、テスト、workflow、script に適用する。

## 使う場面

- 暗黙の前提や広い副作用を持つコードを整理するとき
- shared helper や共通 abstraction を追加・変更するとき
- test が setup や assertion を重複しているとき
- workflow や script の失敗条件を明示したいとき

## 手順

1. 境界を特定する。
2. 契約を言語化する。
3. 前提条件を境界で検証する。
4. 副作用を境界の内側に閉じ込める。
5. 契約を直接検証する test を足す。

## 契約チェックリスト

- 入力が明示されている。
- 出力が明示されている。
- 不変条件が明示されている。
- 許される副作用が狭く定義されている。
- 失敗条件が決定的で利用者に見える。
- 呼び出し側が内部事情を知らなくても正しく使える。

## リファクタ規則

- request/response や setup/teardown の重複は、小さな helper に寄せる。
- validation、orchestration、transport、rendering を 1 つの helper に混ぜすぎない。
- shared global の mutation より constructor や options-based setup を優先する。
- nil や empty の扱いは偶然の挙動に頼らず、契約として明示する。
- test helper は重複排除のために使うが、各 case の意図までは隠さない。

## 出力規則

- 実装説明より先に契約を書く。
- 契約が欠けている場合は、どの境界が契約を持つべきかを先に示す。
- リファクタで挙動が変わる場合は before/after の契約差分を明示する。
