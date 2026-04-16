// licenses パッケージは snipe licenses コマンド（/api/v1/licenses）を提供する。
package licenses

import (
	"context"
	"fmt"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は licenses コマンドを返す。
func NewCmd() *cobra.Command {
	def := &run.ResourceDef{
		Use:     "licenses",
		Short:   "ソフトウェアライセンスを管理する",
		DocsURL: "https://snipe-it.readme.io/reference/licenses",
		APIPath: "licenses",
		ActionFns: []run.ActionDef{
			{
				Use:       "checkout",
				Short:     "ライセンスシートを checkout する",
				Action:    "checkout",
				NeedsData: true,
			},
			{
				Use:       "checkin",
				Short:     "ライセンスシートを checkin する",
				Action:    "checkin",
				NeedsData: false,
			},
		},
	}
	cmd := def.BuildCmd()

	// サブリソース: GET /api/v1/licenses/{id}/{sub}
	cmd.AddCommand(run.BuildSubReadCmd("history", "ライセンスの操作履歴を取得する", "licenses", "history"))
	cmd.AddCommand(buildSeatsCmd())

	return cmd
}

// buildSeatsCmd は licenses seats サブコマンドグループを生成する。
// list/get/update の 3 操作のみ（シートの作成・削除はライセンス本体の CRUD で管理）。
func buildSeatsCmd() *cobra.Command {
	seats := &cobra.Command{
		Use:   "seats",
		Short: "ライセンスシートを管理する",
	}
	seats.AddCommand(run.BuildSubReadCmd("list", "ライセンスのシート一覧を取得する", "licenses", "seats"))
	seats.AddCommand(buildSeatGetCmd())
	seats.AddCommand(buildSeatUpdateCmd())
	return seats
}

type seatGetOptions struct {
	run.BaseOptions
	licenseID int
	seatID    int
}

func buildSeatGetCmd() *cobra.Command {
	o := &seatGetOptions{}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ライセンスシートを ID で取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireAll(
					run.RequirePositiveInt("--id", o.licenseID),
					run.RequirePositiveInt("--seat-id", o.seatID),
				)
			}, func(ctx context.Context) error {
				return run.FetchBySegmentsAndPrint(ctx, &o.BaseOptions, "licenses", fmt.Sprintf("%d", o.licenseID), "seats", fmt.Sprintf("%d", o.seatID))
			})
		},
	}
	cmd.Flags().IntVar(&o.licenseID, "id", 0, "License ID (required)")
	cmd.Flags().IntVar(&o.seatID, "seat-id", 0, "Seat ID (required)")
	cmd.MarkFlagRequired("id")      //nolint:errcheck
	cmd.MarkFlagRequired("seat-id") //nolint:errcheck
	return cmd
}

type seatUpdateOptions struct {
	run.BaseOptions
	licenseID int
	seatID    int
	data      string
}

func buildSeatUpdateCmd() *cobra.Command {
	o := &seatUpdateOptions{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "ライセンスシートを更新する（PATCH）",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireAll(
					run.RequirePositiveInt("--id", o.licenseID),
					run.RequirePositiveInt("--seat-id", o.seatID),
					run.RequireValidJSON("--data", o.data),
				)
			}, func(ctx context.Context) error {
				return run.PatchBySegmentsAndPrint(ctx, &o.BaseOptions, o.data, "licenses", fmt.Sprintf("%d", o.licenseID), "seats", fmt.Sprintf("%d", o.seatID))
			})
		},
	}
	cmd.Flags().IntVar(&o.licenseID, "id", 0, "License ID (required)")
	cmd.Flags().IntVar(&o.seatID, "seat-id", 0, "Seat ID (required)")
	cmd.Flags().StringVar(&o.data, "data", "", "JSON data for fields to update (required)")
	cmd.MarkFlagRequired("id")      //nolint:errcheck
	cmd.MarkFlagRequired("seat-id") //nolint:errcheck
	cmd.MarkFlagRequired("data")    //nolint:errcheck
	return cmd
}
