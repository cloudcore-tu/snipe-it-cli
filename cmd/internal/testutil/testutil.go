// testutil パッケージは cmd/internal/run と cmd の両テストパッケージで共有するテストヘルパーを提供する。
package testutil

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// StartLoopbackServer はランダムポートの loopback HTTP テストサーバーを起動する。
// loopback listener が利用できない場合はテストをスキップする。
func StartLoopbackServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback listener unavailable: %v", err)
		return nil
	}

	srv := httptest.NewUnstartedServer(handler)
	srv.Listener = listener
	srv.Start()
	t.Cleanup(srv.Close)
	return srv
}
