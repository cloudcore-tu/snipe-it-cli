package run

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON は JSON 文字列を map[string]any に変換する。
// create/update コマンドの --data フラグの検証に使用する。
func UnmarshalJSON(data string) (map[string]any, error) {
	var v map[string]any
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return v, nil
}

// ValidateJSON は JSON 文字列として妥当かを検証する。
func ValidateJSON(data string) error {
	var v any
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

// JSONBytes は JSON 文字列を検証してそのまま []byte にして返す。
func JSONBytes(data string) ([]byte, error) {
	if err := ValidateJSON(data); err != nil {
		return nil, err
	}
	return []byte(data), nil
}

// MarshalJSONData は値を JSON に変換する。
func MarshalJSONData(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return data, nil
}
