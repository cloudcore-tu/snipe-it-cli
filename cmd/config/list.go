package config

import (
	"fmt"
	"os"
	"sort"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "登録済みインスタンスを一覧表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			fc, err := config.ReadFile()
			if err != nil {
				return err
			}
			if fc == nil || len(fc.Instances) == 0 {
				fmt.Println("No instances configured. Run 'snip config init' to get started.")
				return nil
			}

			// アクティブなインスタンスを解決（環境変数 > config の current）
			active := os.Getenv("SNIPE_PROFILE")
			if active == "" {
				active = fc.Current
			}

			names := make([]string, 0, len(fc.Instances))
			for name := range fc.Instances {
				names = append(names, name)
			}
			sort.Strings(names)

			path, _ := config.ConfigFilePath()
			fmt.Printf("Config: %s\n\n", path)
			fmt.Printf("%-20s %-50s\n", "NAME", "URL")
			fmt.Printf("%-20s %-50s\n", "----", "---")
			for _, name := range names {
				inst := fc.Instances[name]
				marker := "  "
				if name == active {
					marker = "* "
				}
				fmt.Printf("%s%-18s %s\n", marker, name, inst.URL)
			}

			return nil
		},
	}
}
