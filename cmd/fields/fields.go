// fields パッケージは snipe fields コマンド（/api/v1/fields）を提供する。
// カスタムフィールドの CRUD とフィールドセットへの関連付けを担う。
package fields

import (
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

func buildAssociateCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var id, fieldsetID int
	cmd := &cobra.Command{
		Use:   "associate",
		Short: "フィールドをフィールドセットに関連付ける",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := run.RequirePositiveInt("--id", id); err != nil {
				return err
			}
			if err := run.RequirePositiveInt("--fieldset-id", fieldsetID); err != nil {
				return err
			}
			body, err := run.MarshalJSONData(map[string]int{"fieldset_id": fieldsetID})
			if err != nil {
				return err
			}
			return run.RunPostByPath(cmd.Context(), o,
				fmt.Sprintf("fields/%d/associate", id), body)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Field ID (required)")
	cmd.Flags().IntVar(&fieldsetID, "fieldset-id", 0, "Fieldset ID to associate with (required)")
	cmd.MarkFlagRequired("id")          //nolint:errcheck
	cmd.MarkFlagRequired("fieldset-id") //nolint:errcheck
	return cmd
}

func buildDisassociateCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var id, fieldsetID int
	cmd := &cobra.Command{
		Use:   "disassociate",
		Short: "フィールドをフィールドセットから切り離す",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := run.RequirePositiveInt("--id", id); err != nil {
				return err
			}
			if err := run.RequirePositiveInt("--fieldset-id", fieldsetID); err != nil {
				return err
			}
			body, err := run.MarshalJSONData(map[string]int{"fieldset_id": fieldsetID})
			if err != nil {
				return err
			}
			return run.RunPostByPath(cmd.Context(), o,
				fmt.Sprintf("fields/%d/disassociate", id), body)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Field ID (required)")
	cmd.Flags().IntVar(&fieldsetID, "fieldset-id", 0, "Fieldset ID to disassociate from (required)")
	cmd.MarkFlagRequired("id")          //nolint:errcheck
	cmd.MarkFlagRequired("fieldset-id") //nolint:errcheck
	return cmd
}

// buildReorderCmd は "snip fields reorder --fieldset-id N --data '[1,2,3]'" コマンドを生成する。
// POST /api/v1/fields/fieldsets/{fieldset_id}/order
func buildReorderCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var fieldsetID int
	var data string
	cmd := &cobra.Command{
		Use:   "reorder",
		Short: "フィールドセット内のフィールド並び順を更新する",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := run.RequirePositiveInt("--fieldset-id", fieldsetID); err != nil {
				return err
			}
			payload, err := run.JSONBytes(data)
			if err != nil {
				return err
			}
			return run.RunPostByPath(cmd.Context(), o,
				fmt.Sprintf("fields/fieldsets/%d/order", fieldsetID), payload)
		},
	}
	cmd.Flags().IntVar(&fieldsetID, "fieldset-id", 0, "Fieldset ID (required)")
	cmd.Flags().StringVar(&data, "data", "", "JSON array of field IDs in desired order, e.g. [1,2,3] (required)")
	cmd.MarkFlagRequired("fieldset-id") //nolint:errcheck
	cmd.MarkFlagRequired("data")        //nolint:errcheck
	return cmd
}
