// manufacturers パッケージは snipe manufacturers コマンド（/api/v1/manufacturers）を提供する。
package manufacturers

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は manufacturers コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "manufacturers",
		Short:   "メーカーを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/manufacturers",
		APIPath: "manufacturers",
		ActionFns: []run.ActionDef{
			{
				Use:       "restore",
				Short:     "削除済みメーカーを復元する",
				Action:    "restore",
				NeedsData: false,
			},
		},
	}
	return def.BuildCmd()
}
