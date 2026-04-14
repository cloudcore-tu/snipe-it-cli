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
	return def.BuildCmd()
}
