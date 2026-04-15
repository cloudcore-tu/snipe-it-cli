// settings パッケージは snipe settings コマンド（/api/v1/settings）を提供する。
// 管理者設定の参照・更新とサーバー操作（バックアップ等）を担う。
package settings

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は settings コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Snipe-IT サーバー設定を管理する",
	}

	// GET /api/v1/settings
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

func buildUpdateCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var data string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "設定を更新する（POST /api/v1/settings）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			payload, err := run.JSONBytes(data)
			if err != nil {
				return err
			}
			return run.RunPostByPath(cmd.Context(), o, "settings", payload)
		},
	}
	cmd.Flags().StringVar(&data, "data", "", "JSON data for settings fields to update (required)")
	cmd.MarkFlagRequired("data") //nolint:errcheck
	return cmd
}

// buildBackupDownloadCmd は最新バックアップまたは指定バックアップをダウンロードする。
func buildBackupDownloadCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var name, outputFile string
	cmd := &cobra.Command{
		Use:   "backup-download",
		Short: "バックアップをダウンロードする（省略時は最新）",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			var apiPath string
			if name != "" {
				apiPath = "settings/backups/download/" + name
			} else {
				apiPath = "settings/backups/download/latest"
			}
			return run.RunSaveBinary(cmd.Context(), o, apiPath, outputFile)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Backup file name (default: latest)")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "Save to file (default: stdout)")
	return cmd
}
