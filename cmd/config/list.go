package config

import (
	"fmt"
	"io"
	"sort"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
)

type listView struct {
	path      string
	active    string
	instances map[string]config.Instance
	names     []string
}

func sortedNames(instances map[string]config.Instance) []string {
	names := make([]string, 0, len(instances))
	for name := range instances {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func loadListView() (*listView, error) {
	fc, err := config.ReadFile()
	if err != nil {
		return nil, err
	}
	if fc == nil || len(fc.Instances) == 0 {
		return nil, nil
	}

	path, err := config.ConfigFilePath()
	if err != nil {
		return nil, err
	}

	return &listView{
		path:      path,
		active:    config.ResolveProfile(fc, ""),
		instances: fc.Instances,
		names:     sortedNames(fc.Instances),
	}, nil
}

func renderListView(w io.Writer, view *listView) error {
	if view == nil {
		_, err := fmt.Fprintln(w, "No instances configured. Run 'snip config init' to get started.")
		return err
	}

	if _, err := fmt.Fprintf(w, "Config: %s\n\n", view.path); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%-20s %-50s\n", "NAME", "URL"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%-20s %-50s\n", "----", "---"); err != nil {
		return err
	}
	for _, name := range view.names {
		marker := "  "
		if name == view.active {
			marker = "* "
		}
		if _, err := fmt.Fprintf(w, "%s%-18s %s\n", marker, name, view.instances[name].URL); err != nil {
			return err
		}
	}
	return nil
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "登録済みインスタンスを一覧表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			view, err := loadListView()
			if err != nil {
				return err
			}
			return renderListView(cmd.OutOrStdout(), view)
		},
	}
}
