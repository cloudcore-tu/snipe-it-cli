// snipeit パッケージは Snipe-IT REST API クライアントを提供する。
//
// 公式 Go SDK が存在しないため（既存ライブラリは star 数基準未達）、
// 直接 HTTP で実装する（ADR-001 参照）。
//
// Snipe-IT API は全リソースで一貫したパターンを持つため、
// 汎用メソッド（List/GetByID/Create/Update/Delete）で全 CRUD をカバーできる。
package snipeit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client は Snipe-IT API クライアント。
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient は Snipe-IT API クライアントを初期化する。
// baseURL は "https://snipeit.example.com" 形式（末尾スラッシュなし）。
func NewClient(baseURL, token string, timeoutSec int) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("Snipe-IT URL is not configured (set SNIPEIT_URL or url in config file)")
	}
	if token == "" {
		return nil, fmt.Errorf("API token is not configured (set SNIPEIT_TOKEN or token in config file)")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Snipe-IT URL: %w", err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid Snipe-IT URL: must include scheme and host")
	}

	// 末尾の /api/v1 や / を除去して正規化する
	normalized := normalizeBaseURL(parsedURL)

	return &Client{
		baseURL: normalized,
		token:   token,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
			// loggingTransport で HTTP リクエスト/レスポンスをデバッグログに記録する。
			// Authorization ヘッダは transport 層で自動的にマスクされる。
			Transport: newLoggingTransport(http.DefaultTransport),
		},
	}, nil
}

// normalizeBaseURL は URL パスから末尾スラッシュと /api/v1 サフィックスを除去する。
// ユーザーが "https://snipeit.example.com/api/v1" を入力しても "/api/v1/api/v1/..." を防ぐ。
func normalizeBaseURL(u *url.URL) string {
	path := u.Path
	for {
		prev := path
		if len(path) > 0 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		if len(path) >= 7 && path[len(path)-7:] == "/api/v1" {
			path = path[:len(path)-7]
		}
		if path == prev {
			break
		}
	}
	u.Path = path
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

// apiURL は API エンドポイント URL を組み立てる。
func (c *Client) apiURL(path string) string {
	return c.baseURL + "/api/v1/" + path
}

// doRequest は HTTP リクエストを送信し、レスポンスボディと HTTP ステータスを返す。
func (c *Client) doRequest(ctx context.Context, method, urlStr string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respBody, resp.StatusCode, nil
}

// ListParams は list 操作のクエリパラメータを保持する。
type ListParams struct {
	Limit   int
	Offset  int
	Filters map[string][]string
}

// List は GET /api/v1/{path} でリソース一覧を取得する。
// レスポンスは {"total": N, "rows": [...]} 形式の生 JSON を返す。
func (c *Client) List(ctx context.Context, path string, params ListParams) ([]byte, error) {
	slog.Info("listing resources", "path", path, "limit", params.Limit, "offset", params.Offset)
	u, err := url.Parse(c.apiURL(path))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	for k, vs := range params.Filters {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	body, status, err := c.doRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newAPIError(status, body)
	}
	return body, nil
}

// GetByID は GET /api/v1/{path}/{id} でリソース単体を取得する。
func (c *Client) GetByID(ctx context.Context, path string, id int) ([]byte, error) {
	slog.Info("getting resource", "path", path, "id", id)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id)
	body, status, err := c.doRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newAPIError(status, body)
	}
	return body, nil
}

// Create は POST /api/v1/{path} でリソースを作成する。
// Snipe-IT のレスポンスは {"status": "success", "payload": {...}} のため、payload を取り出して返す。
func (c *Client) Create(ctx context.Context, path string, data []byte) ([]byte, error) {
	slog.Info("creating resource", "path", path)
	body, status, err := c.doRequest(ctx, http.MethodPost, c.apiURL(path), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return nil, newAPIError(status, body)
	}
	return extractPayload(body)
}

// Update は PATCH /api/v1/{path}/{id} でリソースを部分更新する。
// Snipe-IT のレスポンスは {"status": "success", "payload": {...}} のため、payload を取り出して返す。
func (c *Client) Update(ctx context.Context, path string, id int, data []byte) ([]byte, error) {
	slog.Info("updating resource", "path", path, "id", id)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id)
	body, status, err := c.doRequest(ctx, http.MethodPatch, urlStr, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newAPIError(status, body)
	}
	return extractPayload(body)
}

// Delete は DELETE /api/v1/{path}/{id} でリソースを削除する。
func (c *Client) Delete(ctx context.Context, path string, id int) error {
	slog.Info("deleting resource", "path", path, "id", id)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id)
	body, status, err := c.doRequest(ctx, http.MethodDelete, urlStr, nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK && status != http.StatusNoContent {
		return newAPIError(status, body)
	}
	return nil
}

// GetSub は GET /api/v1/{path}/{id}/{subPath} でサブリソースを取得する。
// 例: hardware/42/history, users/5/assets
func (c *Client) GetSub(ctx context.Context, path string, id int, subPath string) ([]byte, error) {
	slog.Info("getting sub-resource", "path", path, "id", id, "sub", subPath)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id) + "/" + subPath
	body, status, err := c.doRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newAPIError(status, body)
	}
	return body, nil
}

// GetByPath は GET /api/v1/{urlPath} で任意のパスを取得する。
// bytag/byserial/reports/account 等の非 CRUD パスに使用する。
func (c *Client) GetByPath(ctx context.Context, urlPath string) ([]byte, error) {
	slog.Info("getting by path", "path", urlPath)
	urlStr := c.apiURL(urlPath)
	body, status, err := c.doRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newAPIError(status, body)
	}
	return body, nil
}

// PatchByPath は PATCH /api/v1/{urlPath} で任意のパスを更新する。
// ライセンスシート更新など、CRUD 汎用メソッドが対応しない入れ子パスに使用する。
func (c *Client) PatchByPath(ctx context.Context, urlPath string, data []byte) ([]byte, error) {
	slog.Info("patching by path", "path", urlPath)
	body, status, err := c.doRequest(ctx, http.MethodPatch, c.apiURL(urlPath), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newAPIError(status, body)
	}
	return extractPayload(body)
}

// PostByPath は POST /api/v1/{urlPath} で任意のパスにアクションを送信する。
// account/request 等の非 CRUD POST に使用する。data が nil の場合は空ボディで送信する。
func (c *Client) PostByPath(ctx context.Context, urlPath string, data []byte) ([]byte, error) {
	slog.Info("posting by path", "path", urlPath)
	var bodyReader io.Reader
	if data != nil {
		bodyReader = bytes.NewReader(data)
	}
	body, status, err := c.doRequest(ctx, http.MethodPost, c.apiURL(urlPath), bodyReader)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return nil, newAPIError(status, body)
	}
	return extractPayload(body)
}

// PostAction は POST /api/v1/{path}/{id}/{action} を呼ぶ。
// checkout/checkin 等のリソース固有操作に使用する。
// data が nil の場合は空ボディで送信する。
func (c *Client) PostAction(ctx context.Context, path string, id int, action string, data []byte) ([]byte, error) {
	slog.Info("posting action", "path", path, "id", id, "action", action)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id) + "/" + action
	var bodyReader io.Reader
	if data != nil {
		bodyReader = bytes.NewReader(data)
	}
	body, status, err := c.doRequest(ctx, http.MethodPost, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return nil, newAPIError(status, body)
	}
	return extractPayload(body)
}

// extractPayload は Snipe-IT の create/update/action レスポンスから payload を取り出す。
// {"status": "success", "payload": {...}} → {...} の生 JSON を返す。
// payload フィールドがない場合はレスポンス全体を返す。
func extractPayload(body []byte) ([]byte, error) {
	var wrapper struct {
		Status   string          `json:"status"`
		Messages json.RawMessage `json:"messages"`
		Payload  json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		// JSON パース失敗時はそのまま返す
		return body, nil
	}

	if wrapper.Status == "error" {
		msg := string(wrapper.Messages)
		return nil, fmt.Errorf("API error: %s", msg)
	}

	if wrapper.Payload != nil {
		return wrapper.Payload, nil
	}
	return body, nil
}

// APIError は Snipe-IT API のエラーレスポンスを表す。
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	msg := string(e.Body)
	if msg == "" {
		return fmt.Sprintf("API error: HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("API error: HTTP %d: %s", e.StatusCode, msg)
}

func newAPIError(statusCode int, body []byte) *APIError {
	return &APIError{StatusCode: statusCode, Body: body}
}
