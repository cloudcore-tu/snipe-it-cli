// run パッケージは全コマンドの Options パターン共通処理を提供する。
// BaseOptions の Complete と Complete → Validate → Run の制御フローを提供する。
// validation helpers は validate.go、JSON helpers は json.go、CRUD フレームワークは resource.go、
// HTTP 実行ヘルパーは helpers.go、バイナリ保存は binary.go に分離している。
package run

import (
	"context"
	"encoding/json"
	"io"
	"os"

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
// Out が未設定の場合は cmd.OutOrStdout() を使う（cobra の SetOut でテスト注入が可能）。
// 副作用: 設定ファイルのパーミッションが 0600 以外の場合、slog.Warn を出力する。
func (o *BaseOptions) Complete(cmd *cobra.Command) error {
	o.resolveOutput(cmd)
	cfg, err := loadConfig(cmd.Root())
	if err != nil {
		return err
	}
	config.WarnInsecurePermissions()
	return o.initFromConfig(cfg)
}

func (o *BaseOptions) resolveOutput(cmd *cobra.Command) {
	if o.Out == nil {
		o.Out = cmd.OutOrStdout()
	}
}

// loadConfig はプロファイル読み込みと CLI フラグ上書きを行い設定を返す。
func loadConfig(root *cobra.Command) (*config.Config, error) {
	profile, _ := root.PersistentFlags().GetString("profile")
	cfg, err := config.Load(profile)
	if err != nil {
		return nil, err
	}
	applyFlagOverrides(cfg, root)
	return cfg, nil
}

// initFromConfig は設定から PrintFlags と Client を初期化する。
func (o *BaseOptions) initFromConfig(cfg *config.Config) error {
	o.PrintFlags = &output.PrintFlags{OutputFormat: cfg.Output}
	client, err := snipeit.NewClient(cfg.URL, cfg.Token, cfg.Timeout)
	if err != nil {
		return err
	}
	o.Client = client
	return nil
}

// applyFlagOverrides はルートコマンドの永続フラグで cfg を上書きする。
// フラグが明示的に指定された（Changed == true）場合のみ上書きする。
// 優先順位: CLI フラグ > 環境変数（config.Load が適用済み）
func applyFlagOverrides(cfg *config.Config, root *cobra.Command) {
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
}

// CompleteValidateRun は BaseOptions を使う command の共通制御。
// Complete -> Validate -> Run の契約を 1 か所に閉じ込める。
func CompleteValidateRun(
	cmd *cobra.Command,
	o *BaseOptions,
	validate func() error,
	run func(context.Context) error,
) error {
	if err := o.Complete(cmd); err != nil {
		return err
	}
	if validate != nil {
		if err := validate(); err != nil {
			return err
		}
	}
	return run(cmd.Context())
}
