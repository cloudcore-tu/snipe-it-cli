// categories パッケージは snipe categories コマンド（/api/v1/categories）を提供する。
package categories

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は categories コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "categories",
		Short:   "カテゴリを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/categories",
		APIPath: "categories",
	}
	return def.BuildCmd()
}
