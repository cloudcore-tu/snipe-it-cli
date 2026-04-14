package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCmd_TextOutput(t *testing.T) {
	cmd := newVersionCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	require.NoError(t, cmd.Execute())

	out := buf.String()
	assert.Contains(t, out, "snipe-it-cli")
	assert.Contains(t, out, "Snipe-IT API")
}

func TestVersionCmd_JSONOutput(t *testing.T) {
	cmd := newVersionCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"-o", "json"})

	require.NoError(t, cmd.Execute())

	var info versionInfo
	require.NoError(t, json.Unmarshal(buf.Bytes(), &info))
	assert.NotEmpty(t, info.ClientVersion)
	assert.NotEmpty(t, info.SnipeITAPI)
}

func TestVersionCmd_UnknownFormat(t *testing.T) {
	cmd := newVersionCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{"-o", "xml"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown output format")
}
