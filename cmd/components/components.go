// components パッケージは snipe components コマンド（/api/v1/components）を提供する。
package components

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は components コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "components",
		Short:   "コンポーネントを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/components",
		APIPath: "components",
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "コンポーネントを checkout する",
				Action:    "checkout",
				NeedsData: true,
			},
			{
				Use:       "checkin",
				Short:     "コンポーネントを checkin する",
				Action:    "checkin",
				NeedsData: false,
			},
		},
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/components/{id}/{sub}
	cmd.AddCommand(run.BuildSubReadCmd("history", "コンポーネントの操作履歴を取得する", "components", "history"))
	cmd.AddCommand(run.BuildSubReadCmd("assets", "コンポーネントが割り当てられた資産を取得する", "components", "assets"))

	return cmd
}
