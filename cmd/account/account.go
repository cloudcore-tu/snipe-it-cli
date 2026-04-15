// account パッケージは snipe account コマンドを提供する。
// ログインユーザー自身の資産リクエスト操作を扱う。
package account

import (
	"context"
	"strconv"

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

	// GET /api/v1/account/eulas — 同意済み EULA 一覧
	cmd.AddCommand(run.BuildPathReadCmd("eulas", "同意済み EULA 一覧を取得する", "account/eulas"))

	// GET /api/v1/account/personal-access-tokens — API トークン一覧
	cmd.AddCommand(run.BuildPathReadCmd("tokens", "自分の API トークン一覧を取得する", "account/personal-access-tokens"))

	// POST /api/v1/account/personal-access-tokens — API トークン作成
	cmd.AddCommand(buildTokenCreateCmd())

	// DELETE /api/v1/account/personal-access-tokens/{id} — API トークン削除
	cmd.AddCommand(buildTokenDeleteCmd())

	return cmd
}

func buildRequestCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var id int
	cmd := &cobra.Command{
		Use:   "request",
		Short: "資産のリクエストを送信する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, o, func() error {
				return run.RequirePositiveInt("--id", id)
			}, func(ctx context.Context) error {
				return run.RunPostByPath(ctx, o, "account/request/"+strconv.Itoa(id), nil)
			})
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
			return run.CompleteValidateRun(cmd, o, func() error {
				return run.RequirePositiveInt("--id", id)
			}, func(ctx context.Context) error {
				return run.RunPostByPath(ctx, o, "account/request/"+strconv.Itoa(id)+"/cancel", nil)
			})
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "Asset ID to cancel request for (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck
	return cmd
}

func buildTokenCreateCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var data string
	cmd := &cobra.Command{
		Use:   "token-create",
		Short: "API トークンを作成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, o, nil, func(ctx context.Context) error {
				return run.RunPostJSONByPath(ctx, o, "account/personal-access-tokens", data)
			})
		},
	}
	cmd.Flags().StringVar(&data, "data", "", `JSON data, e.g. {"name":"my-token"} (required)`)
	cmd.MarkFlagRequired("data") //nolint:errcheck
	return cmd
}

func buildTokenDeleteCmd() *cobra.Command {
	o := &run.BaseOptions{}
	var tokenID int
	cmd := &cobra.Command{
		Use:   "token-delete",
		Short: "API トークンを削除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, o, func() error {
				return run.RequirePositiveInt("--token-id", tokenID)
			}, func(ctx context.Context) error {
				return run.RunDeleteByPath(ctx, o, "account/personal-access-tokens/"+strconv.Itoa(tokenID), tokenID)
			})
		},
	}
	cmd.Flags().IntVar(&tokenID, "token-id", 0, "Token ID to delete (required)")
	cmd.MarkFlagRequired("token-id") //nolint:errcheck
	return cmd
}
