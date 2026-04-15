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
		// create はファイルアップロードのため標準実装を除外し buildCreateCmd で差し替える
		ExcludeSubCmds: []string{"create"},
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
	cmd.AddCommand(buildCreateCmd())
	return cmd
}

type createOptions struct {
	run.BaseOptions
	filePath   string
	importType string
}

// buildCreateCmd は "snip imports create --file PATH --type TYPE" コマンドを生成する。
// POST /api/v1/imports に multipart/form-data でファイルをアップロードする。
func buildCreateCmd() *cobra.Command {
	o := &createOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "CSV ファイルをアップロードしてインポートを作成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireFileExists("--file", o.filePath)
			}, func(ctx context.Context) error {
				return run.RunUpload(ctx, &o.BaseOptions, "imports", "file_contents", o.filePath, map[string]string{"import_type": o.importType})
			})
		},
	}
	cmd.Flags().StringVar(&o.filePath, "file", "", "Path to CSV file (required)")
	cmd.Flags().StringVar(&o.importType, "type", "hardware",
		"Import type: hardware, accessories, licenses, consumables, components, users")
	cmd.MarkFlagRequired("file") //nolint:errcheck
	return cmd
}
