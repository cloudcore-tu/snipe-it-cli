package cmd_test

import (
	"bytes"
	"testing"

	accountcmd "github.com/cloudcore-tu/snipe-it-cli/cmd/account"
	assetscmd "github.com/cloudcore-tu/snipe-it-cli/cmd/assets"
	configcmd "github.com/cloudcore-tu/snipe-it-cli/cmd/config"
	fieldscmd "github.com/cloudcore-tu/snipe-it-cli/cmd/fields"
	licensescmd "github.com/cloudcore-tu/snipe-it-cli/cmd/licenses"
	notescmd "github.com/cloudcore-tu/snipe-it-cli/cmd/notes"
	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRequest_RejectsNonPositiveID(t *testing.T) {
	cmd := accountcmd.NewCmd()
	err := executeCommandForTest(t, cmd, "request", "--id", "0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--id")
}

func TestAssetsByTag_RejectsEmptyTag(t *testing.T) {
	cmd := assetscmd.NewCmd()
	err := executeCommandForTest(t, cmd, "bytag", "--tag", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--tag")
}

func TestLicensesSeatsGet_RejectsNonPositiveSeatID(t *testing.T) {
	cmd := licensescmd.NewCmd()
	err := executeCommandForTest(t, cmd, "seats", "get", "--id", "1", "--seat-id", "0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--seat-id")
}

func TestNotesCreate_RejectsNonPositiveAssetID(t *testing.T) {
	cmd := notescmd.NewCmd()
	err := executeCommandForTest(t, cmd, "create", "--asset-id", "0", "--data", `{"note":"test"}`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--asset-id")
}

func TestFieldsReorder_RejectsInvalidJSON(t *testing.T) {
	cmd := fieldscmd.NewCmd()
	err := executeCommandForTest(t, cmd, "reorder", "--fieldset-id", "1", "--data", "{invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
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
