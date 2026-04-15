package run

import (
	"fmt"
	"os"
	"strings"
)

// ParseFilters は "key=value" 形式の文字列スライスを map[string][]string に変換する。
// --filter フラグの値を API クエリパラメータに変換するために使用する。
//
// 前提条件:
//   - 各要素は "key=value" 形式（key・value ともに非空）
//
// 挙動:
//   - value に "=" が含まれる場合は最初の "=" のみ分割点とする（SplitN 2分割）
//   - 同一キーを複数指定した場合は値をスライスに追加する
//
// 失敗条件:
//   - "=" を含まない → error
//   - key が空 → error
//   - value が空 → error
func ParseFilters(rawFilters []string) (map[string][]string, error) {
	if len(rawFilters) == 0 {
		return nil, nil
	}
	result := make(map[string][]string, len(rawFilters))
	for _, f := range rawFilters {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --filter format: %q (expected: key=value)", f)
		}
		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if k == "" {
			return nil, fmt.Errorf("invalid --filter format: %q (key must not be empty)", f)
		}
		if v == "" {
			return nil, fmt.Errorf("invalid --filter format: %q (value must not be empty)", f)
		}
		result[k] = append(result[k], v)
	}
	return result, nil
}

// RequireAll は複数の validation error を順に評価し、最初のエラーを返す。
func RequireAll(validations ...error) error {
	for _, err := range validations {
		if err != nil {
			return err
		}
	}
	return nil
}

// RequireDeleteConfirmation は --yes なしの削除操作をブロックする。
// 誤削除防止のためすべての delete コマンドで呼び出す。
// 対話的なプロンプトは出さず、--yes フラグで確認を明示させる（エージェント対応）。
func RequireDeleteConfirmation(yes bool) error {
	if yes {
		return nil
	}
	return fmt.Errorf("--yes flag is required to confirm deletion")
}

// RequirePositiveInt は整数フラグが正の値であることを保証する。
func RequirePositiveInt(flagName string, value int) error {
	if value > 0 {
		return nil
	}
	return fmt.Errorf("%s must be a positive integer", flagName)
}

// RequireNonEmpty は文字列フラグが空でないことを保証する。
func RequireNonEmpty(flagName, value string) error {
	if strings.TrimSpace(value) != "" {
		return nil
	}
	return fmt.Errorf("%s must not be empty", flagName)
}

// RequireValidJSON は必須 JSON 文字列フラグが非空かつ妥当な JSON であることを保証する。
func RequireValidJSON(flagName, data string) error {
	return RequireAll(
		RequireNonEmpty(flagName, data),
		ValidateJSON(data),
	)
}

// ValidateOptionalJSON は空文字を許容し、それ以外は妥当な JSON であることを保証する。
func ValidateOptionalJSON(data string) error {
	if strings.TrimSpace(data) == "" {
		return nil
	}
	return ValidateJSON(data)
}

// RequireFileExists はファイルパスが空でなく、既存の通常ファイルを指すことを保証する。
func RequireFileExists(flagName, path string) error {
	if err := RequireNonEmpty(flagName, path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to access %s: %w", flagName, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s must point to a file", flagName)
	}
	return nil
}
