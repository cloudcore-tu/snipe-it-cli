// maintenances パッケージは snipe maintenances コマンド（/api/v1/maintenances）を提供する。
package maintenances

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は maintenances コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "maintenances",
		Short:   "メンテナンス記録を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/maintenances",
		APIPath: "maintenances",
	}
	return def.BuildCmd()
}
