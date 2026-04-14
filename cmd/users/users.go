// users パッケージは snipe users コマンド（/api/v1/users）を提供する。
package users

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は users コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "users",
		Short:   "ユーザーを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/users",
		APIPath: "users",
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/users/{id}/{sub}
	cmd.AddCommand(run.BuildSubReadCmd("assets", "ユーザーに割り当てられた資産を取得する", "users", "assets"))
	cmd.AddCommand(run.BuildSubReadCmd("licenses", "ユーザーに割り当てられたライセンスを取得する", "users", "licenses"))
	cmd.AddCommand(run.BuildSubReadCmd("accessories", "ユーザーに割り当てられたアクセサリーを取得する", "users", "accessories"))
	cmd.AddCommand(run.BuildSubReadCmd("consumables", "ユーザーに割り当てられた消耗品を取得する", "users", "consumables"))

	return cmd
}
