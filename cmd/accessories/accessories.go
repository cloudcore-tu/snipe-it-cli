// accessories パッケージは snipe accessories コマンド（/api/v1/accessories）を提供する。
package accessories

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は accessories コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "accessories",
		Short:   "アクセサリーを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/accessories",
		APIPath: "accessories",
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "アクセサリーを checkout する",
				Action:    "checkout",
				NeedsData: true,
			},
			{
				Use:       "checkin",
				Short:     "アクセサリーを checkin する",
				Action:    "checkin",
				NeedsData: false,
			},
		},
	}
	return def.BuildCmd()
}
