// fieldsets パッケージは snipe fieldsets コマンド（/api/v1/fieldsets）を提供する。
package fieldsets

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は fieldsets コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "fieldsets",
		Short:   "カスタムフィールドセットを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/fieldsets",
		APIPath: "fieldsets",
	}
	return def.BuildCmd()
}
