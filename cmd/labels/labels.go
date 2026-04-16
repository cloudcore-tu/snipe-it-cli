// labels パッケージは snipe labels コマンド（/api/v1/labels）を提供する。
// list はラベルテンプレート一覧を JSON で返す。
// get はラベルを PDF/バイナリで取得し --output-file に保存する（省略時は stdout）。
package labels

import (
	"context"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は labels コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "labels",
		Short: "ラベルテンプレートを管理・出力する",
	}
	cmd.AddCommand(run.BuildPathReadCmd("list", "ラベルテンプレート一覧を取得する", "labels"))
	cmd.AddCommand(buildGetCmd())
	return cmd
}

// buildGetCmd は "snip labels get --name NAME [--output-file PATH]" コマンドを生成する。
// GET /api/v1/labels/{name} の応答はバイナリ（PDF 等）のため --output-file に保存する。
func buildGetCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var name, outputFile string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ラベルを取得してファイルに保存する（省略時は stdout）",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, o, func() error {
				return run.RequireNonEmpty("--name", name)
			}, func(ctx context.Context) error {
				return run.DownloadBySegmentsAndSave(ctx, o, outputFile, "labels", name)
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Label template name (required)")
	cmd.Flags().StringVar(&outputFile, "output-file", "", "Save to file (default: stdout)")
	cmd.MarkFlagRequired("name") //nolint:errcheck
	return cmd
}
