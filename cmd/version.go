package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version は ldflags で注入する。
// go build -ldflags "-X github.com/cloudcore-tu/snipe-it-cli/cmd.version=v0.1.0"
var version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show snipe-it-cli version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("snipe-it-cli %s\n", version)
			fmt.Println("Snipe-IT API: v1 (compatible with Snipe-IT v8.x)")
		},
	}
}
