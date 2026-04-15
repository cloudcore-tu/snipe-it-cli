package config

import (
	"fmt"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
)

type initOptions struct {
	name  string
	url   string
	token string
}

func (o *initOptions) validate() error {
	return validateInstanceInput(o.name, o.url, o.token)
}

func (o *initOptions) fileConfig() *config.FileConfig {
	return &config.FileConfig{
		Current: o.name,
		Instances: map[string]config.Instance{
			o.name: {URL: o.url, Token: o.token},
		},
	}
}

// run は設定ファイルを初期化し、作成されたパスを返す。
func (o *initOptions) run() (string, error) {
	return config.InitFile(o.fileConfig())
}

func newInitCmd() *cobra.Command {
	o := &initOptions{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "初期設定ファイルを生成する",
		Long:  "設定ファイルを新規作成する。すでに存在する場合はエラーになる。",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validate(); err != nil {
				return err
			}
			path, err := o.run()
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Config file created: %s\n", path); err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Active instance: %s\n", o.name)
			return err
		},
	}

	cmd.Flags().StringVar(&o.name, "name", "default", "Instance name")
	cmd.Flags().StringVar(&o.url, "url", "", "Snipe-IT URL (required)")
	cmd.Flags().StringVar(&o.token, "token", "", "API token (required)")
	cmd.MarkFlagRequired("url")   //nolint:errcheck
	cmd.MarkFlagRequired("token") //nolint:errcheck

	return cmd
}
