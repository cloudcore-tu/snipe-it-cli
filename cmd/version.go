package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// version は ldflags で注入する。
// go build -ldflags "-X github.com/cloudcore-tu/snipe-it-cli/cmd.version=v1.0.0"
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
			case "text", "":
				fmt.Fprintf(cmd.OutOrStdout(), "snipe-it-cli %s\n", info.ClientVersion)
				fmt.Fprintf(cmd.OutOrStdout(), "Snipe-IT API: %s\n", info.SnipeITAPI)
			case "json":
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			default:
				// 不正フォーマットは exit 1（ADR-006 準拠: 引数エラーは非ゼロ終了）
				return fmt.Errorf("unknown output format: %q (available: text, json)", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json")

	return cmd
}
