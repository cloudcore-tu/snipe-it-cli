package config

import (
	"fmt"
	"os"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var (
		name  string
		url   string
		token string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "初期設定ファイルを生成する",
		Long:  "設定ファイルを新規作成する。すでに存在する場合はエラーになる。",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 既存ファイル確認
			path, err := config.ConfigFilePath()
			if err != nil {
				return err
			}
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config file already exists: %s (use 'snipe config add' to add an instance)", path)
			}

			fc := &config.FileConfig{
				Current: name,
				Instances: map[string]config.Instance{
					name: {URL: url, Token: token},
				},
			}

			if err := config.WriteFile(fc); err != nil {
				return err
			}

			fmt.Printf("Config file created: %s\n", path)
			fmt.Printf("Active instance: %s\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "default", "Instance name")
	cmd.Flags().StringVar(&url, "url", "", "Snipe-IT URL (required)")
	cmd.Flags().StringVar(&token, "token", "", "API token (required)")
	cmd.MarkFlagRequired("url")   //nolint:errcheck
	cmd.MarkFlagRequired("token") //nolint:errcheck

	return cmd
}
