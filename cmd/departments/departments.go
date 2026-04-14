// departments パッケージは snipe departments コマンド（/api/v1/departments）を提供する。
package departments

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は departments コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "departments",
		Short:   "部門を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/departments",
		APIPath: "departments",
	}
	return def.BuildCmd()
}
