// suppliers パッケージは snipe suppliers コマンド（/api/v1/suppliers）を提供する。
package suppliers

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は suppliers コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "suppliers",
		Short:   "サプライヤーを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/suppliers",
		APIPath: "suppliers",
	}
	return def.BuildCmd()
}
