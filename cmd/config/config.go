// config パッケージは snipe config サブコマンド群を提供する。
package config

import "github.com/spf13/cobra"

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
