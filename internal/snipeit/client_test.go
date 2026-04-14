package snipeit_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startServer はテスト用 HTTP サーバーを起動する。
// ループバックポートのバインドが不可能な制限環境ではテストをスキップする。
func startServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
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

// newTestClient はテストサーバーに向けたクライアントを生成する。
func newTestClient(t *testing.T, srv *httptest.Server) *snipeit.Client {
	t.Helper()
	c, err := snipeit.NewClient(srv.URL, "test-token", 5)
	require.NoError(t, err)
	return c
}

// --- NewClient ---

func TestNewClient_Valid(t *testing.T) {
	c, err := snipeit.NewClient("https://snip.example.com", "token123", 30)
	require.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewClient_EmptyURL(t *testing.T) {
	_, err := snipeit.NewClient("", "token123", 30)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL is not configured")
}

func TestNewClient_EmptyToken(t *testing.T) {
	_, err := snipeit.NewClient("https://snip.example.com", "", 30)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is not configured")
}

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := snipeit.NewClient("not-a-url", "token123", 30)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Snipe-IT URL")
}

func TestNewClient_NormalizesAPIV1Suffix(t *testing.T) {
	// ユーザーが /api/v1 を含む URL を入力しても二重パスにならないことを確認する
	c, err := snipeit.NewClient("https://snip.example.com/api/v1", "token", 30)
	require.NoError(t, err)
	assert.NotNil(t, c)
}

// --- List ---

func TestList_Success(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/hardware", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"total":1,"rows":[{"id":1,"name":"Laptop-001"}]}`)) //nolint:errcheck
	})

	c := newTestClient(t, srv)
	data, err := c.List(context.Background(), "hardware", snipeit.ListParams{Limit: 50})
	require.NoError(t, err)
	assert.Contains(t, string(data), "Laptop-001")
}

func TestList_ServerError(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":"Unauthorized"}`, http.StatusUnauthorized)
	})

	c := newTestClient(t, srv)
	_, err := c.List(context.Background(), "hardware", snipeit.ListParams{Limit: 50})
	assert.Error(t, err)
}

// --- GetByID ---

func TestGetByID_Success(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/hardware/42", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":42,"name":"Laptop-001"}`)) //nolint:errcheck
	})

	c := newTestClient(t, srv)
	data, err := c.GetByID(context.Background(), "hardware", 42)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"id":42`)
}

func TestGetByID_NotFound(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":"No asset found"}`, http.StatusNotFound)
	})

	c := newTestClient(t, srv)
	_, err := c.GetByID(context.Background(), "hardware", 99999)
	assert.Error(t, err)
}

// --- Create ---

func TestCreate_Success(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/hardware", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success","payload":{"id":1,"name":"Laptop-001"}}`)) //nolint:errcheck
	})

	c := newTestClient(t, srv)
	data, err := c.Create(context.Background(), "hardware", []byte(`{"name":"Laptop-001"}`))
	require.NoError(t, err)
	// payload が取り出されていること
	assert.Contains(t, string(data), `"id":1`)
	assert.NotContains(t, string(data), "status")
}

func TestCreate_APIErrorStatus(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"error","messages":"Validation failed"}`)) //nolint:errcheck
	})

	c := newTestClient(t, srv)
	_, err := c.Create(context.Background(), "hardware", []byte(`{"name":""}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Validation failed")
}

// --- Update ---

func TestUpdate_Success(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/v1/hardware/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success","payload":{"id":1,"status_id":3}}`)) //nolint:errcheck
	})

	c := newTestClient(t, srv)
	data, err := c.Update(context.Background(), "hardware", 1, []byte(`{"status_id":3}`))
	require.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v1/hardware/1", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})

	c := newTestClient(t, srv)
	err := c.Delete(context.Background(), "hardware", 1)
	assert.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"messages":"Not found"}`, http.StatusNotFound)
	})

	c := newTestClient(t, srv)
	err := c.Delete(context.Background(), "hardware", 99999)
	assert.Error(t, err)
}

// --- PostAction ---

func TestPostAction_Checkout(t *testing.T) {
	srv := startServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/hardware/1/checkout", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success","payload":{"id":1}}`)) //nolint:errcheck
	})

	c := newTestClient(t, srv)
	data, err := c.PostAction(context.Background(), "hardware", 1, "checkout",
		[]byte(`{"checkout_to_type":"user","assigned_user":5}`))
	require.NoError(t, err)
	assert.Contains(t, string(data), `"id":1`)
}

// --- APIError ---

func TestAPIError_Error_WithBody(t *testing.T) {
	err := &snipeit.APIError{StatusCode: 404, Body: []byte(`{"messages":"Not found"}`)}
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "Not found")
}

func TestAPIError_Error_EmptyBody(t *testing.T) {
	err := &snipeit.APIError{StatusCode: 500, Body: nil}
	assert.Contains(t, err.Error(), "500")
}
