// resource.go は CRUD リソースコマンドの汎用フレームワークを提供する。
//
// Snipe-IT API は全リソースで一貫したパターンを持つため、
// ResourceDef に APIPath を宣言するだけで list/get/create/update/delete コマンドを自動生成できる。
// netbox-cli の ResourceDef より単純: 関数フィールドがなく APIPath のみ（ADR-002 参照）。
package run

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/spf13/cobra"
)

// ActionDef は list/get/create/update/delete 以外の追加アクションを定義する。
// checkout/checkin 等のリソース固有操作に使用する。
type ActionDef struct {
	Use       string
	Short     string
	// Action は API アクションパス（"checkout", "checkin" 等）
	Action    string
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

	// 生 JSON を any として表現してプリンターに渡す
	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return err
	}

	printer, err := o.PrintFlags.NewPrinter(o.Stdout())
	if err != nil {
		return err
	}
	return printer.Print(result)
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
			if o.id <= 0 {
				return fmt.Errorf("--id must be a positive integer")
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

	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return err
	}

	printer, err := o.PrintFlags.NewPrinter(o.Stdout())
	if err != nil {
		return err
	}
	return printer.Print(result)
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
	// JSON の妥当性チェック
	if _, err := UnmarshalJSON(o.data); err != nil {
		return err
	}

	raw, err := o.Client.Create(ctx, o.apiPath, []byte(o.data))
	if err != nil {
		return FormatAPIError(err)
	}

	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return err
	}

	printer, err := o.PrintFlags.NewPrinter(o.Stdout())
	if err != nil {
		return err
	}
	return printer.Print(result)
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
			if o.id <= 0 {
				return fmt.Errorf("--id must be a positive integer")
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
	if _, err := UnmarshalJSON(o.data); err != nil {
		return err
	}

	raw, err := o.Client.Update(ctx, o.apiPath, o.id, []byte(o.data))
	if err != nil {
		return FormatAPIError(err)
	}

	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return err
	}

	printer, err := o.PrintFlags.NewPrinter(o.Stdout())
	if err != nil {
		return err
	}
	return printer.Print(result)
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
			if o.id <= 0 {
				return fmt.Errorf("--id must be a positive integer")
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
	fmt.Fprintf(o.Stdout(), "{\"deleted\":true,\"id\":%d}\n", o.id)
	return nil
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
			if o.id <= 0 {
				return fmt.Errorf("--id must be a positive integer")
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
		if _, err := UnmarshalJSON(o.data); err != nil {
			return err
		}
		dataBytes = []byte(o.data)
	}

	raw, err := o.Client.PostAction(ctx, o.apiPath, o.id, o.action, dataBytes)
	if err != nil {
		return FormatAPIError(err)
	}

	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return err
	}

	printer, err := o.PrintFlags.NewPrinter(o.Stdout())
	if err != nil {
		return err
	}
	return printer.Print(result)
}

// ParseIDs は "--id 1,2,3" や "--id 1 --id 2" 形式の ID 文字列リストを int スライスに変換する。
func parseIDsFromStrings(raw []string) ([]int, error) {
	var ids []int
	for _, s := range raw {
		for _, part := range strings.Split(s, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			n, err := strconv.Atoi(part)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid ID value: %q (must be a positive integer)", part)
			}
			ids = append(ids, n)
		}
	}
	return ids, nil
}
