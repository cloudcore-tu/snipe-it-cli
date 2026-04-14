// account パッケージは snipe account コマンドを提供する。
// ログインユーザー自身の資産リクエスト操作を扱う。
package account

import (
	"fmt"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
)

// NewCmd は account コマンドを返す。
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "アカウント（自分自身）の操作をする",
	}

	// GET /api/v1/account/requestable/hardware — リクエスト可能な資産一覧
	cmd.AddCommand(run.BuildPathReadCmd(
		"requestable",
		"リクエスト可能な資産一覧を取得する",
		"account/requestable/hardware",
	))

	// GET /api/v1/account/requests — 自分のリクエスト一覧
	cmd.AddCommand(run.BuildPathReadCmd(
		"requests",
		"自分が送信したリクエスト一覧を取得する",
		"account/requests",
	))

	// POST /api/v1/account/request/{id} — 資産をリクエストする
	cmd.AddCommand(buildRequestCmd())

	// POST /api/v1/account/request/{id}/cancel — リクエストをキャンセルする
	cmd.AddCommand(buildCancelRequestCmd())

	return cmd
}

func buildRequestCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var id int
	cmd := &cobra.Command{
		Use:   "request",
		Short: "資産のリクエストを送信する",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			return run.RunPostByPath(cmd.Context(), o,
				"account/request/"+itoa(id), nil)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Asset ID to request (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck
	return cmd
}

func buildCancelRequestCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var id int
	cmd := &cobra.Command{
		Use:   "cancel-request",
		Short: "資産リクエストをキャンセルする",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			return run.RunPostByPath(cmd.Context(), o,
				"account/request/"+itoa(id)+"/cancel", nil)
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Asset ID to cancel request for (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck
	return cmd
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
