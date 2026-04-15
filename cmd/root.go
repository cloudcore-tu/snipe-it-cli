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

type rootOptions struct {
	logLevel *slog.LevelVar
	verbose  bool
	debug    bool
}

func newRootOptions() *rootOptions {
	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn)
	return &rootOptions{logLevel: level}
}

func (o *rootOptions) installLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: o.logLevel,
	})))
}

func (o *rootOptions) applyLogFlags() {
	o.logLevel.Set(slog.LevelWarn)
	switch {
	case o.debug:
		o.logLevel.Set(slog.LevelDebug)
	case o.verbose:
		o.logLevel.Set(slog.LevelInfo)
	}
}

func (o *rootOptions) bindPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("url", "", "Snipe-IT URL (env: SNIPEIT_URL)")
	cmd.PersistentFlags().String("token", "", "API token (env: SNIPEIT_TOKEN)")
	cmd.PersistentFlags().String("profile", "", "Instance name to use (env: SNIPE_PROFILE)")
	cmd.PersistentFlags().Int("timeout", 0, "Request timeout in seconds (env: SNIPEIT_TIMEOUT)")
	cmd.PersistentFlags().StringP("output", "o", "", "Output format: table, json, yaml, custom-columns=..., jsonpath=... (env: SNIPEIT_OUTPUT)")
	cmd.PersistentFlags().BoolVarP(&o.verbose, "verbose", "v", false, "Show operational INFO logs")
	cmd.PersistentFlags().BoolVar(&o.debug, "debug", false, "Show debug logs (HTTP requests/responses)")
}

func addResourceCommands(cmd *cobra.Command) {
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(configcmd.NewCmd())
	cmd.AddCommand(assets.NewCmd())
	cmd.AddCommand(users.NewCmd())
	cmd.AddCommand(licenses.NewCmd())
	cmd.AddCommand(categories.NewCmd())
	cmd.AddCommand(locations.NewCmd())
	cmd.AddCommand(manufacturers.NewCmd())
	cmd.AddCommand(models.NewCmd())
	cmd.AddCommand(companies.NewCmd())
	cmd.AddCommand(departments.NewCmd())
	cmd.AddCommand(statuslabels.NewCmd())
	cmd.AddCommand(suppliers.NewCmd())
	cmd.AddCommand(fieldsets.NewCmd())
	cmd.AddCommand(accessories.NewCmd())
	cmd.AddCommand(components.NewCmd())
	cmd.AddCommand(consumables.NewCmd())
	cmd.AddCommand(maintenances.NewCmd())
	cmd.AddCommand(fields.NewCmd())
	cmd.AddCommand(depreciations.NewCmd())
	cmd.AddCommand(groups.NewCmd())
	cmd.AddCommand(reports.NewCmd())
	cmd.AddCommand(account.NewCmd())
	cmd.AddCommand(labels.NewCmd())
	cmd.AddCommand(imports.NewCmd())
	cmd.AddCommand(settings.NewCmd())
	cmd.AddCommand(notes.NewCmd())
}

func newRootCmd() *cobra.Command {
	options := newRootOptions()
	// installLogger は PersistentPreRunE で実行する。
	// newRootCmd() 呼び出し時点でグローバル slog.SetDefault を発火させないことで
	// テストでの parallel 実行時の global state 競合を防ぐ。

	cmd := &cobra.Command{
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
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			options.installLogger()
			options.applyLogFlags()
			return nil
		},
	}

	options.bindPersistentFlags(cmd)
	addResourceCommands(cmd)

	return cmd
}

// Execute はルートコマンドを実行する。main.go から呼ばれる。
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		output.PrintError(os.Stderr, err)
		os.Exit(1)
	}
}
