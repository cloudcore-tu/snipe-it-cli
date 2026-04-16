// fields パッケージは snipe fields コマンド（/api/v1/fields）を提供する。
// カスタムフィールドの CRUD とフィールドセットへの関連付けを担う。
package fields

import (
	"context"
	"fmt"

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
	cmd := def.BuildCmd()

	// POST /api/v1/fields/{id}/associate — フィールドをフィールドセットに関連付ける
	cmd.AddCommand(buildAssociateCmd())
	// POST /api/v1/fields/{id}/disassociate — フィールドをフィールドセットから切り離す
	cmd.AddCommand(buildDisassociateCmd())
	// POST /api/v1/fields/fieldsets/{fieldset_id}/order — フィールドの並び順を更新する
	cmd.AddCommand(buildReorderCmd())

	return cmd
}

type associateOptions struct {
	run.BaseOptions
	id         int
	fieldsetID int
}

func buildAssociateCmd() *cobra.Command {
	o := &associateOptions{}
	cmd := &cobra.Command{
		Use:   "associate",
		Short: "フィールドをフィールドセットに関連付ける",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireAll(
					run.RequirePositiveInt("--id", o.id),
					run.RequirePositiveInt("--fieldset-id", o.fieldsetID),
				)
			}, func(ctx context.Context) error {
				return run.RunPostValueBySegments(ctx, &o.BaseOptions, map[string]int{"fieldset_id": o.fieldsetID}, "fields", fmt.Sprintf("%d", o.id), "associate")
			})
		},
	}
	cmd.Flags().IntVar(&o.id, "id", 0, "Field ID (required)")
	cmd.Flags().IntVar(&o.fieldsetID, "fieldset-id", 0, "Fieldset ID to associate with (required)")
	cmd.MarkFlagRequired("id")          //nolint:errcheck
	cmd.MarkFlagRequired("fieldset-id") //nolint:errcheck
	return cmd
}

type disassociateOptions struct {
	run.BaseOptions
	id         int
	fieldsetID int
}

func buildDisassociateCmd() *cobra.Command {
	o := &disassociateOptions{}
	cmd := &cobra.Command{
		Use:   "disassociate",
		Short: "フィールドをフィールドセットから切り離す",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireAll(
					run.RequirePositiveInt("--id", o.id),
					run.RequirePositiveInt("--fieldset-id", o.fieldsetID),
				)
			}, func(ctx context.Context) error {
				return run.RunPostValueBySegments(ctx, &o.BaseOptions, map[string]int{"fieldset_id": o.fieldsetID}, "fields", fmt.Sprintf("%d", o.id), "disassociate")
			})
		},
	}
	cmd.Flags().IntVar(&o.id, "id", 0, "Field ID (required)")
	cmd.Flags().IntVar(&o.fieldsetID, "fieldset-id", 0, "Fieldset ID to disassociate from (required)")
	cmd.MarkFlagRequired("id")          //nolint:errcheck
	cmd.MarkFlagRequired("fieldset-id") //nolint:errcheck
	return cmd
}

type reorderOptions struct {
	run.BaseOptions
	fieldsetID int
	data       string
}

// buildReorderCmd は "snip fields reorder --fieldset-id N --data '[1,2,3]'" コマンドを生成する。
// POST /api/v1/fields/fieldsets/{fieldset_id}/order
func buildReorderCmd() *cobra.Command {
	o := &reorderOptions{}
	cmd := &cobra.Command{
		Use:   "reorder",
		Short: "フィールドセット内のフィールド並び順を更新する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireAll(
					run.RequirePositiveInt("--fieldset-id", o.fieldsetID),
					run.RequireValidJSON("--data", o.data),
				)
			}, func(ctx context.Context) error {
				return run.RunPostJSONBySegments(ctx, &o.BaseOptions, o.data, "fields", "fieldsets", fmt.Sprintf("%d", o.fieldsetID), "order")
			})
		},
	}
	cmd.Flags().IntVar(&o.fieldsetID, "fieldset-id", 0, "Fieldset ID (required)")
	cmd.Flags().StringVar(&o.data, "data", "", "JSON array of field IDs in desired order, e.g. [1,2,3] (required)")
	cmd.MarkFlagRequired("fieldset-id") //nolint:errcheck
	cmd.MarkFlagRequired("data")        //nolint:errcheck
	return cmd
}
