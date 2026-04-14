// transport.go は HTTP RoundTripper のデバッグログラッパーを提供する。
// Authorization ヘッダを伏字にすることでデバッグログへのトークン漏洩を防ぐ。
package snipeit

import (
	"log/slog"
	"net/http"
	"strings"
)

// loggingTransport はリクエスト/レスポンスを DEBUG レベルでログに記録する RoundTripper。
// --debug フラグが有効な場合にのみログが出力される（slog のレベルフィルタによる）。
type loggingTransport struct {
	base http.RoundTripper
}

func newLoggingTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &loggingTransport{base: base}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	slog.Debug("HTTP request",
		"method", req.Method,
		"url", req.URL.String(),
		// セキュリティ: Bearer トークンを *** でマスク
		"authorization", MaskAuthHeader(req.Header.Get("Authorization")),
	)

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		slog.Debug("HTTP request failed", "error", err)
		return nil, err
	}

	slog.Debug("HTTP response", "status", resp.StatusCode)
	return resp, nil
}

// MaskAuthHeader は "Bearer <token>" の <token> 部分を *** でマスクする。
// デバッグログへの API トークン漏洩を防ぐ。
// エクスポートしてテストから検証できるようにする。
func MaskAuthHeader(header string) string {
	if header == "" {
		return ""
	}
	const bearerPrefix = "Bearer "
	if strings.HasPrefix(header, bearerPrefix) {
		return bearerPrefix + "***"
	}
	// Bearer 以外の認証方式も念のためマスクする
	return "***REDACTED***"
}
