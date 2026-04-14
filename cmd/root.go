// cmd パッケージは snipe-it-cli のコマンド体系を定義する。
package cmd

import (
	"log/slog"
	"os"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/accessories"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/account"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/assets"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/categories"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/companies"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/components"
	configcmd "github.com/cloudcore-tu/snipe-it-cli/cmd/config"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/consumables"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/departments"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/depreciations"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/fields"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/fieldsets"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/groups"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/imports"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/labels"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/licenses"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/locations"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/maintenances"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/manufacturers"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/models"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/notes"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/reports"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/settings"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/statuslabels"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/suppliers"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/users"
	"github.com/cloudcore-tu/snipe-it-cli/internal/output"
	"github.com/spf13/cobra"
)

// logLevel は slog のデフォルトハンドラで使うログレベル変数。
// --verbose で INFO、--debug で DEBUG に切り替わる。
// LevelVar はアトミックに変更できるため並行実行時のデータ競合がない。
var logLevel = new(slog.LevelVar) // デフォルト: INFO

var (
	verbose bool
	debug   bool
)

func init() {
	// デフォルトを WARN に設定し、通常実行では slog 出力を抑制する。
	// --verbose で INFO、--debug で DEBUG に切り替える。
	logLevel.Set(slog.LevelWarn)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// グローバルフラグ
	rootCmd.PersistentFlags().String("url", "", "Snipe-IT URL (env: SNIPEIT_URL)")
	rootCmd.PersistentFlags().String("token", "", "API token (env: SNIPEIT_TOKEN)")
	rootCmd.PersistentFlags().String("profile", "", "Instance name to use (env: SNIPE_PROFILE)")
	rootCmd.PersistentFlags().Int("timeout", 0, "Request timeout in seconds (env: SNIPEIT_TIMEOUT)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format: table, json, yaml, custom-columns=..., jsonpath=... (env: SNIPEIT_OUTPUT)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show operational INFO logs")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Show debug logs (HTTP requests/responses)")

	// サブコマンド登録
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(configcmd.NewCmd())
	rootCmd.AddCommand(assets.NewCmd())
	rootCmd.AddCommand(users.NewCmd())
	rootCmd.AddCommand(licenses.NewCmd())
	rootCmd.AddCommand(categories.NewCmd())
	rootCmd.AddCommand(locations.NewCmd())
	rootCmd.AddCommand(manufacturers.NewCmd())
	rootCmd.AddCommand(models.NewCmd())
	rootCmd.AddCommand(companies.NewCmd())
	rootCmd.AddCommand(departments.NewCmd())
	rootCmd.AddCommand(statuslabels.NewCmd())
	rootCmd.AddCommand(suppliers.NewCmd())
	rootCmd.AddCommand(fieldsets.NewCmd())
	rootCmd.AddCommand(accessories.NewCmd())
	rootCmd.AddCommand(components.NewCmd())
	rootCmd.AddCommand(consumables.NewCmd())
	rootCmd.AddCommand(maintenances.NewCmd())
	rootCmd.AddCommand(fields.NewCmd())
	rootCmd.AddCommand(depreciations.NewCmd())
	rootCmd.AddCommand(groups.NewCmd())
	rootCmd.AddCommand(reports.NewCmd())
	rootCmd.AddCommand(account.NewCmd())
	rootCmd.AddCommand(labels.NewCmd())
	rootCmd.AddCommand(imports.NewCmd())
	rootCmd.AddCommand(settings.NewCmd())
	rootCmd.AddCommand(notes.NewCmd())
}

var rootCmd = &cobra.Command{
	Use:   "snip",
	Short: "Snipe-IT CLI — IT 資産管理ツール",
	Long: `snip は Snipe-IT（IT 資産管理 OSS）を操作する CLI ツールです。

Usage:
  snip [global flags] {resource} {verb} [flags]

Examples:
  snip assets list --filter status_id=2
  snip assets get --id 123
  snip assets create --data '{"name":"Laptop-001","asset_tag":"ASSET-001","model_id":1,"status_id":2}'
  snip users list`,
	SilenceUsage: true,
	// SilenceErrors: エラーは Execute() 側で PrintError して表示するため cobra の自動出力を抑制する
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// フラグの優先順位: --debug > --verbose > デフォルト(WARN)
		switch {
		case debug:
			logLevel.Set(slog.LevelDebug)
		case verbose:
			logLevel.Set(slog.LevelInfo)
		}
		return nil
	},
}

// Execute はルートコマンドを実行する。main.go から呼ばれる。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.PrintError(os.Stderr, err)
		os.Exit(1)
	}
}
