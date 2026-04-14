// assets パッケージは snipe assets コマンド（/api/v1/hardware）を提供する。
// 標準 CRUD に加え、checkout/checkin/audit 操作もサポートする。
package assets

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は assets コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "assets",
		Short:   "IT 資産（ハードウェア）を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/hardware",
		APIPath: "hardware",
		// checkout/checkin は --data を受け付けるアクション
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "資産を checkout する（ユーザー/ロケーションへ割り当て）",
				Action:    "checkout",
				NeedsData: true,
			},
			{
				Use:       "checkin",
				Short:     "資産を checkin する（割り当て解除）",
				Action:    "checkin",
				NeedsData: false,
			},
			{
				Use:       "audit",
				Short:     "資産の監査ログを記録する",
				Action:    "audit",
				NeedsData: false,
			},
		},
	}
	return def.BuildCmd()
}
