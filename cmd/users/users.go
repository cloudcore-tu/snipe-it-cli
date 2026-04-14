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
	return def.BuildCmd()
}
