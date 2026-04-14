// statuslabels パッケージは snipe statuslabels コマンド（/api/v1/statuslabels）を提供する。
package statuslabels

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は statuslabels コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "statuslabels",
		Short:   "ステータスラベルを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/status-labels",
		APIPath: "statuslabels",
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/statuslabels/{id}/assetlist
	cmd.AddCommand(run.BuildSubReadCmd("assetlist", "このステータスラベルを持つ資産一覧を取得する", "statuslabels", "assetlist"))

	// 集計エンドポイント（ID 不要）
	cmd.AddCommand(run.BuildPathReadCmd("counts-by-label", "ステータスラベル別の資産数を取得する", "statuslabels/assets/name"))
	cmd.AddCommand(run.BuildPathReadCmd("counts-by-type", "メタステータス種別別の資産数を取得する", "statuslabels/assets/type"))

	return cmd
}
