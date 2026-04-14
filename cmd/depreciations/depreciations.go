// depreciations パッケージは snipe depreciations コマンド（/api/v1/depreciations）を提供する。
package depreciations

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は depreciations コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "depreciations",
		Short:   "償却設定を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/depreciations",
		APIPath: "depreciations",
	}
	return def.BuildCmd()
}
