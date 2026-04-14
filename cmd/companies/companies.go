// companies パッケージは snipe companies コマンド（/api/v1/companies）を提供する。
package companies

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は companies コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "companies",
		Short:   "会社を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/companies",
		APIPath: "companies",
	}
	return def.BuildCmd()
}
