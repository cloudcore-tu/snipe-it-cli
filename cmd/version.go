package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version は ldflags で注入する。
// go build -ldflags "-X github.com/cloudcore-tu/snipe-it-cli/cmd.version=v0.1.0"
var version = "dev"

// versionInfo は version コマンドの出力構造体。
type versionInfo struct {
	ClientVersion string `json:"clientVersion"`
	SnipeITAPI    string `json:"snipeItAPI"`
}

func newVersionCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show snipe-it-cli and Snipe-IT API version",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := versionInfo{
				ClientVersion: version,
				SnipeITAPI:    "v1 (compatible with Snipe-IT v8.x)",
			}

			switch outputFormat {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			default:
				fmt.Printf("snipe-it-cli %s\n", info.ClientVersion)
				fmt.Printf("Snipe-IT API: %s\n", info.SnipeITAPI)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json")

	return cmd
}
