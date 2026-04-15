// imports パッケージは snipe imports コマンド（/api/v1/imports）を提供する。
// CSV ファイルのアップロードによる一括インポートと処理を担う。
package imports

import (
	"context"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は imports コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "imports",
		Short:   "一括インポートを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/imports",
		APIPath: "imports",
		ActionFns: []run.ActionDef{
			{
				Use:       "process",
				Short:     "インポートを実行する（POST /imports/{id}/process）",
				Action:    "process",
				NeedsData: false,
			},
		},
	}
	cmd := def.BuildCmd()

	// create はファイルアップロードのため標準 create を上書きする
	// ResourceDef.BuildCmd() が生成した create サブコマンドを削除して差し替える
	removeSubCmd(cmd, "create")
	cmd.AddCommand(buildCreateCmd())

	return cmd
}

// removeSubCmd は cobra コマンドから指定名のサブコマンドを削除する。
func removeSubCmd(parent *cobra.Command, name string) {
	for _, sub := range parent.Commands() {
		if sub.Use == name {
			parent.RemoveCommand(sub)
			return
		}
	}
}

// buildCreateCmd は "snip imports create --file PATH --type TYPE" コマンドを生成する。
// POST /api/v1/imports に multipart/form-data でファイルをアップロードする。
func buildCreateCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var filePath, importType string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "CSV ファイルをアップロードしてインポートを作成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, o, func() error {
				return run.RequireFileExists("--file", filePath)
			}, func(ctx context.Context) error {
				return run.RunUpload(ctx, o, "imports", "file_contents", filePath, map[string]string{"import_type": importType})
			})
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "Path to CSV file (required)")
	cmd.Flags().StringVar(&importType, "type", "hardware",
		"Import type: hardware, accessories, licenses, consumables, components, users")
	cmd.MarkFlagRequired("file") //nolint:errcheck
	return cmd
}
