// assets パッケージは snipe assets コマンド（/api/v1/hardware）を提供する。
// 標準 CRUD に加え、checkout/checkin/audit/restore と各種サブリソース参照もサポートする。
package assets

import (
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は assets コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "assets",
		Short:   "IT 資産（ハードウェア）を管理する",
		DocsURL: "https://snipe-it.readme.io/reference/hardware",
		APIPath: "hardware",
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "資産を checkout する（ユーザー/ロケーションへ割り当て）",
				Action:    "checkout",
				NeedsData: true,
			},
			{
				Use:       "checkin",
				Short:     "資産を checkin する（割り当て解除）",
				Action:    "checkin",
				NeedsData: false,
			},
			{
				Use:       "audit",
				Short:     "資産の監査ログを記録する",
				Action:    "audit",
				NeedsData: false,
			},
			{
				Use:       "restore",
				Short:     "削除済み資産を復元する",
				Action:    "restore",
				NeedsData: false,
			},
		},
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/hardware/{id}/{sub}
	cmd.AddCommand(run.BuildSubReadCmd("history", "資産の操作履歴を取得する", "hardware", "history"))
	cmd.AddCommand(run.BuildSubReadCmd("licenses", "資産に紐づくライセンスを取得する", "hardware", "licenses"))
	cmd.AddCommand(run.BuildSubReadCmd("assigned-assets", "資産に割り当てられたサブ資産を取得する", "hardware", "assigned/assets"))
	cmd.AddCommand(run.BuildSubReadCmd("assigned-accessories", "資産に割り当てられたアクセサリーを取得する", "hardware", "assigned/accessories"))
	cmd.AddCommand(run.BuildSubReadCmd("assigned-components", "資産に割り当てられたコンポーネントを取得する", "hardware", "assigned/components"))

	// 資産タグ・シリアル番号による検索: GET /api/v1/hardware/bytag/{tag} など
	cmd.AddCommand(buildByTagCmd())
	cmd.AddCommand(buildBySerialCmd())

	return cmd
}

// buildByTagCmd は "snip assets bytag --tag TAG" コマンドを生成する。
func buildByTagCmd() *cobra.Command {
	var tag string
	o := &run.BaseOptions{}
	cmd := &cobra.Command{
		Use:   "bytag",
		Short: "資産タグで資産を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := run.RequireNonEmpty("--tag", tag); err != nil {
				return err
			}
			return run.RunGetByPath(cmd.Context(), o, "hardware/bytag/"+tag)
		},
	}
	cmd.Flags().StringVar(&tag, "tag", "", "Asset tag (required)")
	cmd.MarkFlagRequired("tag") //nolint:errcheck
	return cmd
}

// buildBySerialCmd は "snip assets byserial --serial SERIAL" コマンドを生成する。
func buildBySerialCmd() *cobra.Command {
	var serial string
	o := &run.BaseOptions{}
	cmd := &cobra.Command{
		Use:   "byserial",
		Short: "シリアル番号で資産を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := run.RequireNonEmpty("--serial", serial); err != nil {
				return err
			}
			return run.RunGetByPath(cmd.Context(), o, "hardware/byserial/"+serial)
		},
	}
	cmd.Flags().StringVar(&serial, "serial", "", "Serial number (required)")
	cmd.MarkFlagRequired("serial") //nolint:errcheck
	return cmd
}
