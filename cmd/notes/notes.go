// notes パッケージは snipe notes コマンドを提供する。
// 資産（hardware）に紐づくノートの参照・追加を担う。
// API パスは /notes/{asset_id}/index|store（/hardware 配下ではない）。
package notes

import (
	"context"
	"fmt"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は notes コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notes",
		Short: "資産ノートを管理する",
	}
	cmd.AddCommand(buildListCmd())
	cmd.AddCommand(buildCreateCmd())
	return cmd
}

type listOptions struct {
	run.BaseOptions
	assetID int
}

// buildListCmd は "snip notes list --asset-id N" コマンドを生成する。
// GET /api/v1/notes/{N}/index
func buildListCmd() *cobra.Command {
	o := &listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "資産のノート一覧を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequirePositiveInt("--asset-id", o.assetID)
			}, func(ctx context.Context) error {
				return run.FetchBySegmentsAndPrint(ctx, &o.BaseOptions, "notes", fmt.Sprintf("%d", o.assetID), "index")
			})
		},
	}
	cmd.Flags().IntVar(&o.assetID, "asset-id", 0, "Asset (hardware) ID (required)")
	cmd.MarkFlagRequired("asset-id") //nolint:errcheck
	return cmd
}

type createOptions struct {
	run.BaseOptions
	assetID int
	data    string
}

// buildCreateCmd は "snip notes create --asset-id N --data JSON" コマンドを生成する。
// POST /api/v1/notes/{N}/store
func buildCreateCmd() *cobra.Command {
	o := &createOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "資産にノートを追加する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireAll(
					run.RequirePositiveInt("--asset-id", o.assetID),
					run.RequireValidJSON("--data", o.data),
				)
			}, func(ctx context.Context) error {
				return run.PostJSONBySegmentsAndPrint(ctx, &o.BaseOptions, o.data, "notes", fmt.Sprintf("%d", o.assetID), "store")
			})
		},
	}
	cmd.Flags().IntVar(&o.assetID, "asset-id", 0, "Asset (hardware) ID (required)")
	cmd.Flags().StringVar(&o.data, "data", "", `JSON data, e.g. {"note":"text"} (required)`)
	cmd.MarkFlagRequired("asset-id") //nolint:errcheck
	cmd.MarkFlagRequired("data")     //nolint:errcheck
	return cmd
}
