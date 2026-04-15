package config

import (
	"fmt"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var (
		url   string
		token string
	)

	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: "インスタンスを追加・更新する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			fc, err := config.ReadFile()
			if err != nil {
				return err
			}
			if fc == nil {
				fc = &config.FileConfig{
					Current:   name,
					Instances: make(map[string]config.Instance),
				}
			}
			if fc.Instances == nil {
				fc.Instances = make(map[string]config.Instance)
			}

			fc.Instances[name] = config.Instance{URL: url, Token: token}
			// current が未設定の場合は最初に追加したインスタンスをデフォルトにする
			if fc.Current == "" {
				fc.Current = name
			}

			if err := config.WriteFile(fc); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Instance %q added/updated.\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Snipe-IT URL (required)")
	cmd.Flags().StringVar(&token, "token", "", "API token (required)")
	cmd.MarkFlagRequired("url")   //nolint:errcheck
	cmd.MarkFlagRequired("token") //nolint:errcheck

	return cmd
}
