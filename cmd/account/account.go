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

type requestOptions struct {
	run.BaseOptions
	id int
}

func buildRequestCmd() *cobra.Command {
	o := &requestOptions{}
	cmd := &cobra.Command{
		Use:   "request",
		Short: "資産のリクエストを送信する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequirePositiveInt("--id", o.id)
			}, func(ctx context.Context) error {
				return run.RunPostBySegments(ctx, &o.BaseOptions, nil, "account", "request", strconv.Itoa(o.id))
			})
		},
	}
	cmd.Flags().IntVar(&o.id, "id", 0, "Asset ID to request (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck
	return cmd
}

type cancelRequestOptions struct {
	run.BaseOptions
	id int
}

func buildCancelRequestCmd() *cobra.Command {
	o := &cancelRequestOptions{}
	cmd := &cobra.Command{
		Use:   "cancel-request",
		Short: "資産リクエストをキャンセルする",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequirePositiveInt("--id", o.id)
			}, func(ctx context.Context) error {
				return run.RunPostBySegments(ctx, &o.BaseOptions, nil, "account", "request", strconv.Itoa(o.id), "cancel")
			})
		},
	}
	cmd.Flags().IntVar(&o.id, "id", 0, "Asset ID to cancel request for (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck
	return cmd
}

type tokenCreateOptions struct {
	run.BaseOptions
	data string
}

func buildTokenCreateCmd() *cobra.Command {
	o := &tokenCreateOptions{}
	cmd := &cobra.Command{
		Use:   "token-create",
		Short: "API トークンを作成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequireValidJSON("--data", o.data)
			}, func(ctx context.Context) error {
				return run.RunPostJSONByPath(ctx, &o.BaseOptions, "account/personal-access-tokens", o.data)
			})
		},
	}
	cmd.Flags().StringVar(&o.data, "data", "", `JSON data, e.g. {"name":"my-token"} (required)`)
	cmd.MarkFlagRequired("data") //nolint:errcheck
	return cmd
}

type tokenDeleteOptions struct {
	run.BaseOptions
	tokenID int
}

func buildTokenDeleteCmd() *cobra.Command {
	o := &tokenDeleteOptions{}
	cmd := &cobra.Command{
		Use:   "token-delete",
		Short: "API トークンを削除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run.CompleteValidateRun(cmd, &o.BaseOptions, func() error {
				return run.RequirePositiveInt("--token-id", o.tokenID)
			}, func(ctx context.Context) error {
				return run.RunDeleteBySegments(ctx, &o.BaseOptions, "account", "personal-access-tokens", strconv.Itoa(o.tokenID))
			})
		},
	}
	cmd.Flags().IntVar(&o.tokenID, "token-id", 0, "Token ID to delete (required)")
	cmd.MarkFlagRequired("token-id") //nolint:errcheck
	return cmd
}
