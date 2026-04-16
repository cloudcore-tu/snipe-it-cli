package run_test

import (
	"bytes"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/testutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newResourceTestRoot(cmd *cobra.Command) *cobra.Command {
	root := &cobra.Command{Use: "snip"}
	root.PersistentFlags().String("profile", "", "")
	root.PersistentFlags().String("url", "", "")
	root.PersistentFlags().String("token", "", "")
	root.PersistentFlags().Int("timeout", 0, "")
	root.PersistentFlags().String("output", "", "")
	root.AddCommand(cmd)
	return root
}

func executeResourceCommand(t *testing.T, cmd *cobra.Command, baseURL string, args ...string) (string, error) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SNIPEIT_URL", baseURL)
	t.Setenv("SNIPEIT_TOKEN", "test-token")
	t.Setenv("SNIPEIT_OUTPUT", "json")

	root := newResourceTestRoot(cmd)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs(append([]string{cmd.Name()}, args...))
	err := root.Execute()
	return stdout.String(), err
}

func newTestResourceDef() *run.ResourceDef {
	return &run.ResourceDef{
		Use:     "assets",
		Short:   "IT assets",
		APIPath: "hardware",
		ActionFns: []run.ActionDef{
			{Use: "checkout", Short: "checkout", Action: "checkout", NeedsData: true},
			{Use: "checkin", Short: "checkin", Action: "checkin", NeedsData: false},
		},
	}
}

func TestResourceListCommand_ReturnsRows(t *testing.T) {
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/hardware", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"total":2,"rows":[{"id":1,"name":"Laptop-001"},{"id":2,"name":"Laptop-002"}]}`))
		require.NoError(t, err)
	}))

	out, err := executeResourceCommand(t, newTestResourceDef().BuildCmd(), srv.URL, "list")
	require.NoError(t, err)
	assert.Contains(t, out, "Laptop-001")
	assert.Contains(t, out, "Laptop-002")
}

func TestResourceListCommand_PropagatesFilters(t *testing.T) {
	var query string
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"total":0,"rows":[]}`))
		require.NoError(t, err)
	}))

	_, err := executeResourceCommand(t, newTestResourceDef().BuildCmd(), srv.URL, "list", "--filter", "status_id=2")
	require.NoError(t, err)
	assert.Contains(t, query, "status_id=2")
}

func TestResourceGetCommand_ReturnsResource(t *testing.T) {
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/hardware/42", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"id":42,"name":"Laptop-042"}`))
		require.NoError(t, err)
	}))

	out, err := executeResourceCommand(t, newTestResourceDef().BuildCmd(), srv.URL, "get", "--id", "42")
	require.NoError(t, err)
	assert.Contains(t, out, `"id": 42`)
	assert.Contains(t, out, "Laptop-042")
}

func TestResourceCreateCommand_SendsPostAndReturnsPayload(t *testing.T) {
	var capturedBody []byte
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := ioReadAll(r)
		require.NoError(t, err)
		capturedBody = body
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"status":"success","payload":{"id":1,"name":"Laptop-001"}}`))
		require.NoError(t, err)
	}))

	out, err := executeResourceCommand(
		t,
		newTestResourceDef().BuildCmd(),
		srv.URL,
		"create",
		"--data",
		`{"name":"Laptop-001","asset_tag":"A001","model_id":1,"status_id":2}`,
	)
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"Laptop-001","asset_tag":"A001","model_id":1,"status_id":2}`, string(capturedBody))
	assert.Contains(t, out, `"id": 1`)
	assert.NotContains(t, out, "status")
}

func TestResourceUpdateCommand_SendsPatchRequest(t *testing.T) {
	var (
		capturedMethod string
		capturedPath   string
		capturedBody   []byte
	)
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		body, err := ioReadAll(r)
		require.NoError(t, err)
		capturedBody = body
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"status":"success","payload":{"id":42,"status_id":3}}`))
		require.NoError(t, err)
	}))

	_, err := executeResourceCommand(
		t,
		newTestResourceDef().BuildCmd(),
		srv.URL,
		"update",
		"--id",
		"42",
		"--data",
		`{"status_id":3}`,
	)
	require.NoError(t, err)
	assert.Equal(t, http.MethodPatch, capturedMethod)
	assert.Equal(t, "/api/v1/hardware/42", capturedPath)
	assert.JSONEq(t, `{"status_id":3}`, string(capturedBody))
}

func TestResourceDeleteCommand_WithoutYesDoesNotCallAPI(t *testing.T) {
	var calls atomic.Int32
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusOK)
	}))

	_, err := executeResourceCommand(t, newTestResourceDef().BuildCmd(), srv.URL, "delete", "--id", "1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--yes")
	assert.Zero(t, calls.Load())
}

func TestResourceDeleteCommand_WithYesCallsAPIAndOutputs(t *testing.T) {
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/hardware/5", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))

	out, err := executeResourceCommand(t, newTestResourceDef().BuildCmd(), srv.URL, "delete", "--id", "5", "--yes")
	require.NoError(t, err)
	assert.Contains(t, out, `"deleted": true`)
	assert.Contains(t, out, `"id": 5`)
}

func TestResourceActionCommand_CheckoutSendsPostWithData(t *testing.T) {
	var (
		capturedPath string
		capturedBody []byte
	)
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := ioReadAll(r)
		require.NoError(t, err)
		capturedBody = body
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"status":"success","payload":{"id":1}}`))
		require.NoError(t, err)
	}))

	_, err := executeResourceCommand(
		t,
		newTestResourceDef().BuildCmd(),
		srv.URL,
		"checkout",
		"--id",
		"1",
		"--data",
		`{"checkout_to_type":"user","assigned_user":5}`,
	)
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/hardware/1/checkout", capturedPath)
	assert.JSONEq(t, `{"checkout_to_type":"user","assigned_user":5}`, string(capturedBody))
}

func TestResourceActionCommand_CheckinSendsPostWithoutData(t *testing.T) {
	var capturedPath string
	srv := testutil.StartLoopbackServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status":"success","payload":{"id":1}}`))
		require.NoError(t, err)
	}))

	_, err := executeResourceCommand(t, newTestResourceDef().BuildCmd(), srv.URL, "checkin", "--id", "1")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/hardware/1/checkin", capturedPath)
}

func TestBuildCmd_HasExpectedSubcommands(t *testing.T) {
	cmd := newTestResourceDef().BuildCmd()
	assert.Equal(t, "assets", cmd.Use)

	subCmds := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subCmds[sub.Use] = true
	}

	for _, expected := range []string{"list", "get", "create", "update", "delete"} {
		assert.True(t, subCmds[expected], "subcommand %q should exist", expected)
	}
}

func TestBuildCmd_WithActionFns(t *testing.T) {
	cmd := newTestResourceDef().BuildCmd()

	subCmds := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subCmds[sub.Use] = true
	}

	assert.True(t, subCmds["checkout"])
	assert.True(t, subCmds["checkin"])
}

func TestBuildCmd_ExcludeSubCmds(t *testing.T) {
	cmd := (&run.ResourceDef{
		Use:            "imports",
		Short:          "imports",
		APIPath:        "imports",
		ExcludeSubCmds: []string{"create"},
	}).BuildCmd()

	subCmds := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subCmds[sub.Use] = true
	}

	assert.False(t, subCmds["create"])
	assert.True(t, subCmds["list"])
	assert.True(t, subCmds["get"])
	assert.True(t, subCmds["update"])
	assert.True(t, subCmds["delete"])
}

func ioReadAll(r *http.Request) ([]byte, error) {
	defer r.Body.Close() //nolint:errcheck
	return io.ReadAll(r.Body)
}
