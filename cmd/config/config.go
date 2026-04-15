// config パッケージは snipe config サブコマンド群を提供する。
package config

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

func validateInstanceInput(name, url, token string) error {
	return run.RequireAll(
		run.RequireNonEmpty("--name", name),
		run.RequireNonEmpty("--url", url),
		run.RequireNonEmpty("--token", token),
	)
}

// NewCmd は config コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Snipe-IT インスタンスの設定を管理する",
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newListCmd())

	return cmd
}
