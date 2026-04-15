package cmd_test

import (
	"bytes"
	"os"
	"testing"

	accountcmd "github.com/cloudcore-tu/snipe-it-cli/cmd/account"
	assetscmd "github.com/cloudcore-tu/snipe-it-cli/cmd/assets"
	configcmd "github.com/cloudcore-tu/snipe-it-cli/cmd/config"
	fieldscmd "github.com/cloudcore-tu/snipe-it-cli/cmd/fields"
	importscmd "github.com/cloudcore-tu/snipe-it-cli/cmd/imports"
	labelscmd "github.com/cloudcore-tu/snipe-it-cli/cmd/labels"
	licensescmd "github.com/cloudcore-tu/snipe-it-cli/cmd/licenses"
	notescmd "github.com/cloudcore-tu/snipe-it-cli/cmd/notes"
	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandValidation(t *testing.T) {
	existingFile := writeValidationFixture(t)

	testCases := []struct {
		name       string
		newCommand func() *cobra.Command
		args       []string
		want       string
	}{
		{
			name:       "account request rejects non-positive id",
			newCommand: accountcmd.NewCmd,
			args:       []string{"request", "--id", "0"},
			want:       "--id",
		},
		{
			name:       "assets bytag rejects empty tag",
			newCommand: assetscmd.NewCmd,
			args:       []string{"bytag", "--tag", ""},
			want:       "--tag",
		},
		{
			name:       "fields reorder rejects invalid json",
			newCommand: fieldscmd.NewCmd,
			args:       []string{"reorder", "--fieldset-id", "1", "--data", "{invalid"},
			want:       "failed to parse JSON",
		},
		{
			name:       "imports create rejects missing file",
			newCommand: importscmd.NewCmd,
			args:       []string{"create", "--file", "missing.csv"},
			want:       "failed to access --file",
		},
		{
			name:       "imports create accepts existing file before API call",
			newCommand: importscmd.NewCmd,
			args:       []string{"create", "--file", existingFile},
			want:       "Post",
		},
		{
			name:       "labels get rejects empty name",
			newCommand: labelscmd.NewCmd,
			args:       []string{"get", "--name", ""},
			want:       "--name",
		},
		{
			name:       "licenses seats get rejects non-positive seat id",
			newCommand: licensescmd.NewCmd,
			args:       []string{"seats", "get", "--id", "1", "--seat-id", "0"},
			want:       "--seat-id",
		},
		{
			name:       "notes create rejects non-positive asset id",
			newCommand: notescmd.NewCmd,
			args:       []string{"create", "--asset-id", "0", "--data", `{"note":"test"}`},
			want:       "--asset-id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := executeCommandForTest(t, tc.newCommand(), tc.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestConfigList_WritesToCommandOutput(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	require.NoError(t, config.WriteFile(&config.FileConfig{
		Current: "prod",
		Instances: map[string]config.Instance{
			"prod": {URL: "https://snip.example.com", Token: "secret"},
		},
	}))

	cmd := configcmd.NewCmd()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs([]string{"list"})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, stdout.String(), "prod")
	assert.Contains(t, stdout.String(), "https://snip.example.com")
}

func executeCommandForTest(t *testing.T, cmd *cobra.Command, args ...string) error {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SNIPEIT_URL", "https://snip.example.com")
	t.Setenv("SNIPEIT_TOKEN", "test-token")
	t.Setenv("SNIPEIT_OUTPUT", "json")

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func writeValidationFixture(t *testing.T) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "import-*.csv")
	require.NoError(t, err)
	t.Cleanup(func() {
		file.Close() //nolint:errcheck
	})
	_, err = file.WriteString("asset_tag,name\nA-001,Laptop\n")
	require.NoError(t, err)
	return file.Name()
}
