package run

import (
	"encoding/json"
	"fmt"
)

// ParseJSONObject は JSON 文字列を map[string]any に変換する。
// create/update コマンドの --data フラグの検証に使用する。
//
// 前提条件: data は JSON object 形式（{...}）であること。配列は error を返す。
func ParseJSONObject(data string) (map[string]any, error) {
	var v map[string]any
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return v, nil
}

// CheckJSONSyntax は JSON 文字列として構文が妥当かを検証する。
// object/array/primitive いずれも受け付ける。
func CheckJSONSyntax(data string) error {
	var v any
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

// ParseJSONBytes は JSON 文字列の構文を検証し []byte に変換して返す。
// HTTP ヘルパーに渡す前の入力検証に使用する。
func ParseJSONBytes(data string) ([]byte, error) {
	if err := CheckJSONSyntax(data); err != nil {
		return nil, err
	}
	return []byte(data), nil
}

// EncodeJSON は値を JSON []byte に変換する。
// エラーメッセージを統一するための薄いラッパー。
func EncodeJSON(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return data, nil
}
