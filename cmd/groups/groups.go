// groups パッケージは snipe groups コマンド（/api/v1/groups）を提供する。
package groups

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は groups コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "groups",
		Short:   "権限グループを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/groups",
		APIPath: "groups",
	}
	return def.BuildCmd()
}
