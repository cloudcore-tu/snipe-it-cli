// reports パッケージは snipe reports コマンドを提供する。
// activity（操作履歴）と depreciation（償却）の 2 種類のレポートを扱う。
package reports

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は reports コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reports",
		Short: "各種レポートを取得する",
	}

	// GET /api/v1/reports/activity
	cmd.AddCommand(run.BuildPathReadCmd(
		"activity",
		"全リソースのアクティビティレポートを取得する",
		"reports/activity",
	))

	// GET /api/v1/reports/depreciation
	cmd.AddCommand(run.BuildPathReadCmd(
		"depreciation",
		"資産の償却レポートを取得する",
		"reports/depreciation",
	))

	return cmd
}
