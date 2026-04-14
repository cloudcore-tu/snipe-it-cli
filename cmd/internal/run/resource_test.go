// resource_test.go は genericXxxOptions の Run 関数を直接テストする白箱テスト。
// package run（非 _test）にすることで unexported 型に直接アクセスする。
package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/internal/output"
	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startResourceTestServer はテスト用 HTTP サーバーを起動する。
// ループバックポートのバインドが不可能な制限環境ではテストをスキップする。
func startResourceTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback listener unavailable: %v", err)
		return nil
	}
	srv := httptest.NewUnstartedServer(handler)
	srv.Listener = l
	srv.Start()
	t.Cleanup(srv.Close)
	return srv
}

// newResourceBaseOptions はリソーステスト用の BaseOptions を生成する。出力は JSON 固定。
func newResourceBaseOptions(client *snipeit.Client, out *bytes.Buffer) BaseOptions {
	return BaseOptions{
		Client:     client,
		PrintFlags: &output.PrintFlags{OutputFormat: "json"},
		Out:        out,
	}
}

// newResourceTestClient はテストサーバーに向けた snipeit.Client を生成する。
func newResourceTestClient(t *testing.T, srv *httptest.Server) *snipeit.Client {
	t.Helper()
	c, err := snipeit.NewClient(srv.URL, "test-token", 5)
	require.NoError(t, err)
	return c
}

// --- runList ---

func TestRunList_ReturnsRows(t *testing.T) {
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/hardware", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"total":2,"rows":[{"id":1,"name":"Laptop-001"},{"id":2,"name":"Laptop-002"}]}`)
	})

	var buf bytes.Buffer
	o := &genericListOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		limit:       50,
	}

	require.NoError(t, o.runList(context.Background()))

	out := buf.String()
	assert.Contains(t, out, "Laptop-001")
	assert.Contains(t, out, "Laptop-002")
}

func TestRunList_FilterPropagated(t *testing.T) {
	var capturedQuery string
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"total":0,"rows":[]}`)
	})

	var buf bytes.Buffer
	o := &genericListOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		limit:       50,
		filters:     []string{"status_id=2"},
	}

	require.NoError(t, o.runList(context.Background()))
	assert.Contains(t, capturedQuery, "status_id=2")
}

func TestRunList_APIError_PropagatesError(t *testing.T) {
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":"Unauthorized"}`, http.StatusUnauthorized)
	})

	var buf bytes.Buffer
	o := &genericListOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		limit:       50,
	}

	assert.Error(t, o.runList(context.Background()))
}

// --- runGet ---

func TestRunGet_ReturnsResource(t *testing.T) {
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/hardware/42", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":42,"name":"Laptop-042"}`)
	})

	var buf bytes.Buffer
	o := &genericGetOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		id:          42,
	}

	require.NoError(t, o.runGet(context.Background()))
	assert.Contains(t, buf.String(), `"id": 42`)
	assert.Contains(t, buf.String(), "Laptop-042")
}

func TestRunGet_NotFound_ReturnsError(t *testing.T) {
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":"Not found"}`, http.StatusNotFound)
	})

	var buf bytes.Buffer
	o := &genericGetOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		id:          99999,
	}

	assert.Error(t, o.runGet(context.Background()))
}

// --- runCreate ---

func TestRunCreate_SendsPostAndReturnsPayload(t *testing.T) {
	var capturedBody []byte
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"success","payload":{"id":1,"name":"Laptop-001"}}`)
	})

	var buf bytes.Buffer
	o := &genericCreateOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		data:        `{"name":"Laptop-001","asset_tag":"A001","model_id":1,"status_id":2}`,
	}

	require.NoError(t, o.runCreate(context.Background()))
	assert.JSONEq(t, `{"name":"Laptop-001","asset_tag":"A001","model_id":1,"status_id":2}`, string(capturedBody))
	// payload が取り出されている
	assert.Contains(t, buf.String(), `"id": 1`)
	assert.NotContains(t, buf.String(), "status")
}

func TestRunCreate_InvalidJSON_ReturnsError(t *testing.T) {
	// 不正 JSON は HTTP リクエストを送らずエラーを返す
	client, _ := snipeit.NewClient("http://127.0.0.1:1", "test-token", 1)
	var buf bytes.Buffer
	o := &genericCreateOptions{
		BaseOptions: newResourceBaseOptions(client, &buf),
		apiPath:     "hardware",
		data:        `{not valid}`,
	}

	err := o.runCreate(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

// --- runUpdate ---

func TestRunUpdate_SendsPatchRequest(t *testing.T) {
	var (
		capturedMethod string
		capturedPath   string
		capturedBody   []byte
	)
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"success","payload":{"id":42,"status_id":3}}`)
	})

	var buf bytes.Buffer
	o := &genericUpdateOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		id:          42,
		data:        `{"status_id":3}`,
	}

	require.NoError(t, o.runUpdate(context.Background()))
	assert.Equal(t, http.MethodPatch, capturedMethod)
	assert.Equal(t, "/api/v1/hardware/42", capturedPath)
	assert.JSONEq(t, `{"status_id":3}`, string(capturedBody))
}

// --- runDelete ---

func TestRunDelete_WithoutYes_ReturnsError(t *testing.T) {
	// --yes なしは削除しない（HTTP リクエストを送らない）
	client, _ := snipeit.NewClient("http://127.0.0.1:1", "test-token", 1)
	var buf bytes.Buffer
	o := &genericDeleteOptions{
		BaseOptions: newResourceBaseOptions(client, &buf),
		apiPath:     "hardware",
		id:          1,
		yes:         false,
	}

	err := o.runDelete(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--yes")
}

func TestRunDelete_WithYes_CallsAPIAndOutputs(t *testing.T) {
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/hardware/5", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})

	var buf bytes.Buffer
	o := &genericDeleteOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		id:          5,
		yes:         true,
	}

	require.NoError(t, o.runDelete(context.Background()))
	assert.Contains(t, buf.String(), `"deleted": true`)
	assert.Contains(t, buf.String(), `"id": 5`)
}

// --- runAction (checkout/checkin) ---

func TestRunAction_Checkout_SendsPostWithData(t *testing.T) {
	var capturedPath string
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"success","payload":{"id":1}}`)
	})

	var buf bytes.Buffer
	o := &genericActionOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		action:      "checkout",
		id:          1,
		data:        `{"checkout_to_type":"user","assigned_user":5}`,
	}

	require.NoError(t, o.runAction(context.Background()))
	assert.Equal(t, "/api/v1/hardware/1/checkout", capturedPath)
}

func TestRunAction_Checkin_SendsPostWithoutData(t *testing.T) {
	var capturedPath string
	srv := startResourceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"success","payload":{"id":1}}`)
	})

	var buf bytes.Buffer
	o := &genericActionOptions{
		BaseOptions: newResourceBaseOptions(newResourceTestClient(t, srv), &buf),
		apiPath:     "hardware",
		action:      "checkin",
		id:          1,
		data:        "", // checkin にはデータ不要
	}

	require.NoError(t, o.runAction(context.Background()))
	assert.Equal(t, "/api/v1/hardware/1/checkin", capturedPath)
}

// --- BuildCmd (ResourceDef) ---

func TestBuildCmd_HasExpectedSubcommands(t *testing.T) {
	def := &ResourceDef{
		Use:     "assets",
		Short:   "IT 資産を管理する",
		APIPath: "hardware",
	}

	cmd := def.BuildCmd()
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
	def := &ResourceDef{
		Use:     "assets",
		Short:   "IT 資産を管理する",
		APIPath: "hardware",
		ActionFns: []ActionDef{
			{Use: "checkout", Short: "チェックアウト", Action: "checkout", NeedsData: true},
			{Use: "checkin", Short: "チェックイン", Action: "checkin", NeedsData: false},
		},
	}

	cmd := def.BuildCmd()
	subCmds := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		subCmds[sub.Use] = true
	}

	assert.True(t, subCmds["checkout"])
	assert.True(t, subCmds["checkin"])
}
