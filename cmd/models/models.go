// models パッケージは snipe models コマンド（/api/v1/models）を提供する。
package models

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は models コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "models",
		Short:   "機器モデルを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/models",
		APIPath: "models",
		ActionFns: []run.ActionDef{
			{
				Use:       "restore",
				Short:     "削除済みモデルを復元する",
				Action:    "restore",
				NeedsData: false,
			},
		},
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/models/{id}/history
	cmd.AddCommand(run.BuildSubReadCmd("history", "機器モデルの操作履歴を取得する", "models", "history"))

	return cmd
}
