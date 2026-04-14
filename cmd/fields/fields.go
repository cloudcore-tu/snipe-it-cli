// fields パッケージは snipe fields コマンド（/api/v1/fields）を提供する。
// カスタムフィールドの CRUD を担う。フィールドセットへの関連付けには fieldsets コマンドを使用する。
package fields

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は fields コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "fields",
		Short:   "カスタムフィールドを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/fields",
		APIPath: "fields",
	}
	return def.BuildCmd()
}
