// settings パッケージは snipe settings コマンド（/api/v1/settings）を提供する。
// 管理者設定の参照・更新とサーバー操作（バックアップ等）を担う。
package settings

import (
	"context"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は settings コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Snipe-IT サーバー設定を管理する",
	}

	// GET /api/v1/settings/{setting}
	// Note: Snipe-IT v8.3.7 currently routes this to Api\SettingsController::show,
	// but that controller method does not exist upstream and the API returns HTTP 500.
	// Keep the command for route parity, but do not rely on it in smoke tests.
	cmd.AddCommand(run.BuildPathReadCmd("get", "現在の設定を取得する", "settings/general"))

	// POST /api/v1/settings — JSON で設定を更新する（PATCH でなく POST）
	cmd.AddCommand(buildUpdateCmd())

	// ログイン試行履歴
	cmd.AddCommand(run.BuildPathReadCmd("login-attempts", "ログイン試行履歴を取得する", "settings/login-attempts"))

	// バックアップ一覧
	cmd.AddCommand(run.BuildPathReadCmd("backups", "バックアップ一覧を取得する", "settings/backups"))

	// バックアップダウンロード
	cmd.AddCommand(buildBackupDownloadCmd())

	return cmd
}

type updateOptions struct {
	run.BaseOptions
	data string
}

func buildUpdateCmd() *cobra.Command {
	o := &updateOptions{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "設定を更新する（POST /api/v1/settings）",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireValidJSON("--data", o.data)
			}, func(ctx context.Context) error {
				return run.RunPostJSONByPath(ctx, &o.BaseOptions, "settings", o.data)
			})
		},
	}
	cmd.Flags().StringVar(&o.data, "data", "", "JSON data for settings fields to update (required)")
	cmd.MarkFlagRequired("data") //nolint:errcheck
	return cmd
}

type backupDownloadOptions struct {
	run.BaseOptions
	name       string
	outputFile string
}

// buildBackupDownloadCmd は最新バックアップまたは指定バックアップをダウンロードする。
func buildBackupDownloadCmd() *cobra.Command {
	o := &backupDownloadOptions{}
	cmd := &cobra.Command{
		Use:   "backup-download",
		Short: "バックアップをダウンロードする（省略時は最新）",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, nil, func(ctx context.Context) error {
				if o.name != "" {
					return run.RunSaveBinaryBySegments(ctx, &o.BaseOptions, o.outputFile, "settings", "backups", "download", o.name)
				}
				return run.RunSaveBinaryBySegments(ctx, &o.BaseOptions, o.outputFile, "settings", "backups", "download", "latest")
			})
		},
	}
	cmd.Flags().StringVar(&o.name, "name", "", "Backup file name (default: latest)")
	cmd.Flags().StringVar(&o.outputFile, "output-file", "", "Save to file (default: stdout)")
	return cmd
}
