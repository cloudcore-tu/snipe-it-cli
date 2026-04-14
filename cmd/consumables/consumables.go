// consumables パッケージは snipe consumables コマンド（/api/v1/consumables）を提供する。
package consumables

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は consumables コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "consumables",
		Short:   "消耗品を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/consumables",
		APIPath: "consumables",
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "消耗品をユーザーへ払い出す",
				Action:    "checkout",
				NeedsData: true,
			},
		},
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/consumables/{id}/{sub}
	cmd.AddCommand(run.BuildSubReadCmd("history", "消耗品の操作履歴を取得する", "consumables", "history"))
	cmd.AddCommand(run.BuildSubReadCmd("users", "消耗品を受け取ったユーザー一覧を取得する", "consumables", "users"))

	return cmd
}
