// cmd パッケージは snipe-it-cli のコマンド体系を定義する。
package cmd

import (
	"log/slog"
	"os"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/accessories"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/assets"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/categories"
	configcmd "github.com/cloudcore-tu/snipe-it-cli/cmd/config"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/companies"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/components"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/consumables"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/departments"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/fieldsets"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/licenses"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/locations"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/maintenances"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/manufacturers"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/models"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/statuslabels"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/suppliers"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/users"
	"github.com/cloudcore-tu/snipe-it-cli/internal/output"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:          "snip",
	Short:        "Snipe-IT CLI — IT 資産管理ツール",
	SilenceUsage: true,
	// SilenceErrors: エラーは Execute() 側で PrintError して表示するため cobra の自動出力を抑制する
	SilenceErrors: true,
	Long: `snip は Snipe-IT（IT 資産管理 OSS）を操作する CLI ツールです。

Usage:
  snip [global flags] {resource} {verb} [flags]

Examples:
  snip assets list --filter status_id=2
  snip assets get --id 123
  snip assets create --data '{"name":"Laptop-001","asset_tag":"ASSET-001","model_id":1,"status_id":2}'
  snip users list`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// --debug が指定された場合はログレベルを DEBUG に変更する
		if debug {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})))
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

func init() {
	// グローバルフラグ
	rootCmd.PersistentFlags().String("url", "", "Snipe-IT URL (env: SNIPEIT_URL)")
	rootCmd.PersistentFlags().String("token", "", "API token (env: SNIPEIT_TOKEN)")
	rootCmd.PersistentFlags().String("profile", "", "Instance name to use (env: SNIPE_PROFILE)")
	rootCmd.PersistentFlags().Int("timeout", 0, "Request timeout in seconds (env: SNIPEIT_TIMEOUT)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format: table, json, yaml, custom-columns=..., jsonpath=... (env: SNIPEIT_OUTPUT)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

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
}
