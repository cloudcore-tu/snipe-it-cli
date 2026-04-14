// locations パッケージは snipe locations コマンド（/api/v1/locations）を提供する。
package locations

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は locations コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "locations",
		Short:   "ロケーションを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/locations",
		APIPath: "locations",
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/locations/{id}/{sub}
	cmd.AddCommand(run.BuildSubReadCmd("users", "ロケーションに所属するユーザーを取得する", "locations", "users"))
	cmd.AddCommand(run.BuildSubReadCmd("assets", "ロケーションの資産一覧を取得する", "locations", "assets"))
	cmd.AddCommand(run.BuildSubReadCmd("assigned-assets", "ロケーションに割り当てられた資産を取得する", "locations", "assigned/assets"))
	cmd.AddCommand(run.BuildSubReadCmd("assigned-accessories", "ロケーションに割り当てられたアクセサリーを取得する", "locations", "assigned/accessories"))
	cmd.AddCommand(run.BuildSubReadCmd("history", "ロケーションの操作履歴を取得する", "locations", "history"))

	return cmd
}
