// resource.go は CRUD リソースコマンドの汎用フレームワークを提供する。
//
// Snipe-IT API は全リソースで一貫したパターンを持つため、
// ResourceDef に APIPath を宣言するだけで list/get/create/update/delete コマンドを自動生成できる。
// netbox-cli の ResourceDef より単純: 関数フィールドがなく APIPath のみ（ADR-002 参照）。
package run

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/spf13/cobra"
)

// ActionDef は list/get/create/update/delete 以外の追加アクションを定義する。
// checkout/checkin 等のリソース固有操作に使用する。
type ActionDef struct {
	Use   string
	Short string
	// Action は API アクションパス（"checkout", "checkin" 等）
	Action string
	// NeedsData は --data フラグを受け付けるか
	NeedsData bool
}

// ResourceDef はリソースの CRUD 操作定義を保持する。
// APIPath を設定するだけで標準 CRUD コマンドが生成される。
type ResourceDef struct {
	Use     string
	Short   string
	DocsURL string

	// APIPath は API のリソースパス（例: "hardware", "users", "licenses"）
	APIPath string

	// ActionFns は標準 CRUD 以外の追加アクションコマンドを定義する
	ActionFns []ActionDef
}

// BuildCmd は ResourceDef から cobra.Command（親コマンド＋サブコマンド群）を生成する。
func (r *ResourceDef) BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   r.Use,
		Short: r.Short,
	}
	if r.DocsURL != "" {
		cmd.Long = r.Short + "\n\nDocs: " + r.DocsURL
	}

	if r.APIPath != "" {
		cmd.AddCommand(r.buildListCmd())
		cmd.AddCommand(r.buildGetCmd())
		cmd.AddCommand(r.buildCreateCmd())
		cmd.AddCommand(r.buildUpdateCmd())
		cmd.AddCommand(r.buildDeleteCmd())
	}
	for _, action := range r.ActionFns {
		cmd.AddCommand(r.buildActionCmd(action))
	}

	return cmd
}

// --- list ---

type genericListOptions struct {
	BaseOptions
	limit   int
	offset  int
	filters []string
	apiPath string
}

func (r *ResourceDef) buildListCmd() *cobra.Command {
	o := &genericListOptions{apiPath: r.APIPath}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List " + r.Use,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if o.limit < 1 || o.limit > 1000 {
				return fmt.Errorf("--limit must be between 1 and 1000")
			}
			if o.offset < 0 {
				return fmt.Errorf("--offset must be 0 or greater")
			}
			return o.runList(cmd.Context())
		},
	}

	cmd.Flags().IntVar(&o.limit, "limit", 50, "Maximum number of results per page (max 1000)")
	cmd.Flags().IntVar(&o.offset, "offset", 0, "Starting position for results")
	cmd.Flags().StringArrayVar(&o.filters, "filter", nil, "Filter (key=value, can be specified multiple times)")

	return cmd
}

func (o *genericListOptions) runList(ctx context.Context) error {
	filters, err := ParseFilters(o.filters)
	if err != nil {
		return err
	}

	raw, err := o.Client.List(ctx, o.apiPath, snipeit.ListParams{
		Limit:   o.limit,
		Offset:  o.offset,
		Filters: filters,
	})
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// --- get ---

type genericGetOptions struct {
	BaseOptions
	id      int
	apiPath string
}

func (r *ResourceDef) buildGetCmd() *cobra.Command {
	o := &genericGetOptions{apiPath: r.APIPath}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get " + r.Use + " by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := RequirePositiveInt("--id", o.id); err != nil {
				return err
			}
			return o.runGet(cmd.Context())
		},
	}

	cmd.Flags().IntVar(&o.id, "id", 0, "Resource ID (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck

	return cmd
}

func (o *genericGetOptions) runGet(ctx context.Context) error {
	raw, err := o.Client.GetByID(ctx, o.apiPath, o.id)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// --- create ---

type genericCreateOptions struct {
	BaseOptions
	data    string
	apiPath string
}

func (r *ResourceDef) buildCreateCmd() *cobra.Command {
	o := &genericCreateOptions{apiPath: r.APIPath}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create " + r.Use,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			return o.runCreate(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&o.data, "data", "", "JSON data for the resource to create (required)")
	cmd.MarkFlagRequired("data") //nolint:errcheck

	return cmd
}

func (o *genericCreateOptions) runCreate(ctx context.Context) error {
	data, err := JSONBytes(o.data)
	if err != nil {
		return err
	}

	raw, err := o.Client.Create(ctx, o.apiPath, data)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// --- update ---

type genericUpdateOptions struct {
	BaseOptions
	id      int
	data    string
	apiPath string
}

func (r *ResourceDef) buildUpdateCmd() *cobra.Command {
	o := &genericUpdateOptions{apiPath: r.APIPath}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update " + r.Use + " (PATCH)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := RequirePositiveInt("--id", o.id); err != nil {
				return err
			}
			return o.runUpdate(cmd.Context())
		},
	}

	cmd.Flags().IntVar(&o.id, "id", 0, "Resource ID (required)")
	cmd.Flags().StringVar(&o.data, "data", "", "JSON data for fields to update (required)")
	cmd.MarkFlagRequired("id")   //nolint:errcheck
	cmd.MarkFlagRequired("data") //nolint:errcheck

	return cmd
}

func (o *genericUpdateOptions) runUpdate(ctx context.Context) error {
	data, err := JSONBytes(o.data)
	if err != nil {
		return err
	}

	raw, err := o.Client.Update(ctx, o.apiPath, o.id, data)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// --- delete ---

type genericDeleteOptions struct {
	BaseOptions
	id      int
	yes     bool
	apiPath string
}

func (r *ResourceDef) buildDeleteCmd() *cobra.Command {
	o := &genericDeleteOptions{apiPath: r.APIPath}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete " + r.Use,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := RequirePositiveInt("--id", o.id); err != nil {
				return err
			}
			return o.runDelete(cmd.Context())
		},
	}

	cmd.Flags().IntVar(&o.id, "id", 0, "Resource ID (required)")
	// 誤削除防止のため --yes を必須とする。エージェントはこのフラグを明示的に渡す。
	cmd.Flags().BoolVar(&o.yes, "yes", false, "Confirm deletion")
	cmd.MarkFlagRequired("id") //nolint:errcheck

	return cmd
}

func (o *genericDeleteOptions) runDelete(ctx context.Context) error {
	if err := RequireDeleteConfirmation(o.yes); err != nil {
		return err
	}
	if err := o.Client.Delete(ctx, o.apiPath, o.id); err != nil {
		return FormatAPIError(err)
	}
	return o.PrintValue(map[string]any{"deleted": true, "id": o.id})
}

// --- action (checkout/checkin 等) ---

type genericActionOptions struct {
	BaseOptions
	id      int
	data    string
	apiPath string
	action  string
}

func (r *ResourceDef) buildActionCmd(actionDef ActionDef) *cobra.Command {
	o := &genericActionOptions{apiPath: r.APIPath, action: actionDef.Action}

	cmd := &cobra.Command{
		Use:   actionDef.Use,
		Short: actionDef.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := RequirePositiveInt("--id", o.id); err != nil {
				return err
			}
			return o.runAction(cmd.Context())
		},
	}

	cmd.Flags().IntVar(&o.id, "id", 0, "Resource ID (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck

	if actionDef.NeedsData {
		cmd.Flags().StringVar(&o.data, "data", "", "JSON data for the action (required)")
		cmd.MarkFlagRequired("data") //nolint:errcheck
	}

	return cmd
}

func (o *genericActionOptions) runAction(ctx context.Context) error {
	var dataBytes []byte
	if o.data != "" {
		var err error
		dataBytes, err = JSONBytes(o.data)
		if err != nil {
			return err
		}
	}

	raw, err := o.Client.PostAction(ctx, o.apiPath, o.id, o.action, dataBytes)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// --- サブリソース取得ヘルパー ---
// BuildSubReadCmd と BuildPathReadCmd は ResourceDef を拡張せず、
// 各リソースパッケージから明示的に呼ばれるヘルパー関数として提供する。

// subReadOptions は "GET /api/v1/{parentPath}/{id}/{subPath}" の共通 Options。
type subReadOptions struct {
	BaseOptions
	id         int
	parentPath string
	subPath    string
}

func (o *subReadOptions) run(ctx context.Context) error {
	raw, err := o.Client.GetSub(ctx, o.parentPath, o.id, o.subPath)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// BuildSubReadCmd は "snip {resource} {use} --id N" コマンドを生成する。
// GET /api/v1/{parentAPIPath}/{N}/{subPath} を呼び出す。
// 例: BuildSubReadCmd("history", "資産の操作履歴", "hardware", "history")
//
//	→ snip assets history --id 42
func BuildSubReadCmd(use, short, parentAPIPath, subPath string) *cobra.Command {
	o := &subReadOptions{parentPath: parentAPIPath, subPath: subPath}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := RequirePositiveInt("--id", o.id); err != nil {
				return err
			}
			return o.run(cmd.Context())
		},
	}
	cmd.Flags().IntVar(&o.id, "id", 0, "Resource ID (required)")
	cmd.MarkFlagRequired("id") //nolint:errcheck
	return cmd
}

// pathReadOptions は "GET /api/v1/{urlPath}" の共通 Options（ID なし）。
type pathReadOptions struct {
	BaseOptions
	urlPath string
}

func (o *pathReadOptions) run(ctx context.Context) error {
	raw, err := o.Client.GetByPath(ctx, o.urlPath)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// RunGetByPath は初期化済みの BaseOptions を使って GET /api/v1/{urlPath} を実行し結果を出力する。
// 各パッケージで文字列フラグからパスを組み立てるコマンドに使用する（bytag, byserial 等）。
func RunGetByPath(ctx context.Context, o *BaseOptions, urlPath string) error {
	raw, err := o.Client.GetByPath(ctx, urlPath)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// BuildPathReadCmd は固定 API パスに GET する引数なしコマンドを生成する。
// 例: BuildPathReadCmd("activity", "アクティビティレポート", "reports/activity")
//
//	→ snip reports activity
func BuildPathReadCmd(use, short, apiPath string) *cobra.Command {
	o := &pathReadOptions{urlPath: apiPath}
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			return o.run(cmd.Context())
		},
	}
}

// RunPatchByPath は初期化済みの BaseOptions を使って PATCH /api/v1/{urlPath} を実行し結果を出力する。
// ライセンスシート等の入れ子 PATCH エンドポイントに使用する。
func RunPatchByPath(ctx context.Context, o *BaseOptions, urlPath, data string) error {
	payload, err := JSONBytes(data)
	if err != nil {
		return err
	}
	raw, err := o.Client.PatchByPath(ctx, urlPath, payload)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// RunPostByPath は初期化済みの BaseOptions を使って POST /api/v1/{urlPath} を実行し結果を出力する。
// account/request 等の非 CRUD POST エンドポイントに使用する。
func RunPostByPath(ctx context.Context, o *BaseOptions, urlPath string, data []byte) error {
	raw, err := o.Client.PostByPath(ctx, urlPath, data)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// RunDeleteByPath は初期化済みの BaseOptions を使って DELETE /api/v1/{urlPath} を実行する。
// account/personal-access-tokens/{id} 等の非 CRUD DELETE に使用する。
func RunDeleteByPath(ctx context.Context, o *BaseOptions, urlPath string, id int) error {
	if err := o.Client.DeleteByPath(ctx, urlPath); err != nil {
		return FormatAPIError(err)
	}
	return o.PrintValue(map[string]any{"deleted": true, "id": id})
}

// RunUpload は multipart/form-data でファイルをアップロードし結果を出力する。
// Snipe-IT のインポート API に使用する。
func RunUpload(ctx context.Context, o *BaseOptions, urlPath, fieldName, filePath string, extraFields map[string]string) error {
	raw, err := o.Client.Upload(ctx, urlPath, fieldName, filePath, extraFields)
	if err != nil {
		return FormatAPIError(err)
	}
	return o.PrintResponse(raw)
}

// RunSaveBinary は GET /api/v1/{urlPath} のレスポンスをバイナリとして保存する。
// outputFile が空の場合は標準出力に書き出す（パイプ利用を想定）。
// JSON パースを行わないため、PDF 等のバイナリレスポンスに使用する。
func RunSaveBinary(ctx context.Context, o *BaseOptions, urlPath, outputFile string) error {
	raw, err := o.Client.GetByPath(ctx, urlPath)
	if err != nil {
		return FormatAPIError(err)
	}
	if outputFile != "" {
		if err := os.WriteFile(outputFile, raw, 0o600); err != nil {
			return fmt.Errorf("failed to write to %s: %w", outputFile, err)
		}
		fmt.Fprintf(o.Stdout(), "Saved %d bytes to %s\n", len(raw), outputFile)
		return nil
	}
	_, err = o.Stdout().Write(raw)
	return err
}
