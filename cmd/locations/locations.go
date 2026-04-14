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
	return def.BuildCmd()
}
