// consumables パッケージは snipe consumables コマンド（/api/v1/consumables）を提供する。
package consumables

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は consumables コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "consumables",
		Short:   "消耗品を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/consumables",
		APIPath: "consumables",
	}
	return def.BuildCmd()
}
