// licenses パッケージは snipe licenses コマンド（/api/v1/licenses）を提供する。
package licenses

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は licenses コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "licenses",
		Short:   "ソフトウェアライセンスを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/licenses",
		APIPath: "licenses",
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "ライセンスシートを checkout する",
				Action:    "checkout",
				NeedsData: true,
			},
			{
				Use:       "checkin",
				Short:     "ライセンスシートを checkin する",
				Action:    "checkin",
				NeedsData: false,
			},
		},
	}
	return def.BuildCmd()
}
