// run パッケージは全コマンドの Options パターン共通処理を提供する。
// BaseOptions の Complete、フィルタ解析、API エラーラップ等。
package run

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/cloudcore-tu/snipe-it-cli/internal/output"
	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/spf13/cobra"
)

// BaseOptions は全コマンドの Options 構造体が埋め込む共通フィールド。
// kubectl の Options パターン（Complete → Validate → Run）を採用する。
type BaseOptions struct {
	Client     *snipeit.Client
	PrintFlags *output.PrintFlags
	// Out は出力先。nil の場合は os.Stdout を使う。
	// テスト時に bytes.Buffer を注入することでコマンド出力を検証できる。
	Out io.Writer
}

// Stdout は出力先を返す。Out が nil なら os.Stdout を返す。
func (o *BaseOptions) Stdout() io.Writer {
	if o.Out != nil {
		return o.Out
	}
	return os.Stdout
}

func (o *BaseOptions) printer() (*output.Printer, error) {
	return o.PrintFlags.NewPrinter(o.Stdout())
}

// PrintValue は初期化済みの出力設定で値を描画する。
func (o *BaseOptions) PrintValue(v any) error {
	printer, err := o.printer()
	if err != nil {
		return err
	}
	return printer.Print(v)
}

// PrintResponse は JSON レスポンスをデコードして出力する。
func (o *BaseOptions) PrintResponse(raw []byte) error {
	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return err
	}
	return o.PrintValue(result)
}

// Complete はグローバルフラグ・環境変数・設定ファイルをマージして BaseOptions を初期化する。
func (o *BaseOptions) Complete(cmd *cobra.Command) error {
	profile, _ := cmd.Root().PersistentFlags().GetString("profile")
	cfg, err := config.Load(profile)
	if err != nil {
		return err
	}

	// CLI フラグによる上書き（フラグが明示的に指定された場合のみ）
	root := cmd.Root()
	if root.PersistentFlags().Changed("url") {
		cfg.URL, _ = root.PersistentFlags().GetString("url")
	}
	if root.PersistentFlags().Changed("token") {
		cfg.Token, _ = root.PersistentFlags().GetString("token")
	}
	if root.PersistentFlags().Changed("timeout") {
		cfg.Timeout, _ = root.PersistentFlags().GetInt("timeout")
	}
	if root.PersistentFlags().Changed("output") {
		cfg.Output, _ = root.PersistentFlags().GetString("output")
	}

	o.PrintFlags = &output.PrintFlags{OutputFormat: cfg.Output}

	client, err := snipeit.NewClient(cfg.URL, cfg.Token, cfg.Timeout)
	if err != nil {
		return err
	}
	o.Client = client
	return nil
}

// ParseFilters は "key=value" 形式の文字列スライスを map[string][]string に変換する。
// --filter フラグの値を API クエリパラメータに変換するために使用する。
// 同一キーの複数指定に対応する（例: --filter status=1 --filter category_id=2）。
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
		result[k] = append(result[k], v)
	}
	return result, nil
}

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

// RequireDeleteConfirmation は --yes なしの削除操作をブロックする。
// 誤削除防止のためすべての delete コマンドで呼び出す。
// 対話的なプロンプトは出さず、--yes フラグで確認を明示させる（エージェント対応）。
func RequireDeleteConfirmation(yes bool) error {
	if yes {
		return nil
	}
	return fmt.Errorf("--yes flag is required to confirm deletion")
}

// FormatAPIError は snipeit.APIError をユーザー向けのエラーに変換する。
func FormatAPIError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w", err)
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
