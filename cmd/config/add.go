package config

import (
	"fmt"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
)

type addOptions struct {
	name  string
	url   string
	token string
}

func (o *addOptions) validate() error {
	return validateInstanceInput(o.name, o.url, o.token)
}

func (o *addOptions) run() error {
	return config.UpsertInstance(o.name, config.Instance{URL: o.url, Token: o.token})
}

func newAddCmd() *cobra.Command {
	o := &addOptions{}

	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: "インスタンスを追加・更新する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Instance %q added/updated.\n", o.name)
			return err
		},
	}

	cmd.Flags().StringVar(&o.url, "url", "", "Snipe-IT URL (required)")
	cmd.Flags().StringVar(&o.token, "token", "", "API token (required)")
	cmd.MarkFlagRequired("url")   //nolint:errcheck
	cmd.MarkFlagRequired("token") //nolint:errcheck

	return cmd
}
