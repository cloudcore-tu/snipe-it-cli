// testutil パッケージはテスト用の共通ユーティリティを提供する。
// cmd/ 配下の全テストパッケージから利用する。本番コードでは使用しない。
package testutil

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/cloudcore-tu/snipe-it-cli/internal/output"
	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/stretchr/testify/require"
)

// NewBaseOptions はテスト用の BaseOptions を生成する。出力フォーマットは JSON 固定。
func NewBaseOptions(client *snipeit.Client, out *bytes.Buffer) run.BaseOptions {
	return run.BaseOptions{
		Client:     client,
		PrintFlags: &output.PrintFlags{OutputFormat: "json"},
		Out:        out,
	}
}

// NewServer は指定パスにのみレスポンスを返す HTTP テストサーバーを生成する。
// パスが一致しない場合は 404 を返しテストを失敗させる。
// ループバックポートのバインドが不可能な制限環境ではテストをスキップする。
func NewServer(t *testing.T, path, body string) *httptest.Server {
	t.Helper()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("unexpected path: got %q, want %q", r.URL.Path, path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	})
	return startTestServer(t, handler)
}

// NewClientFromServer はテストサーバーに接続する snipeit.Client を生成する。
func NewClientFromServer(t *testing.T, srv *httptest.Server) *snipeit.Client {
	t.Helper()
	client, err := snipeit.NewClient(srv.URL, "test-token", 5)
	require.NoError(t, err)
	return client
}

// startTestServer はループバックポートでテストサーバーを起動する。
func startTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback listener unavailable (restricted environment): %v", err)
		return nil
	}
	srv := httptest.NewUnstartedServer(handler)
	srv.Listener = l
	srv.Start()
	t.Cleanup(srv.Close)
	return srv
}
