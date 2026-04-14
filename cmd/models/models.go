// models パッケージは snipe models コマンド（/api/v1/models）を提供する。
package models

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は models コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "models",
		Short:   "機器モデルを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/models",
		APIPath: "models",
	}
	return def.BuildCmd()
}
