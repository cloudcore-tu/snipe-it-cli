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
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Client は Snipe-IT API クライアント。
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type requestOptions struct {
	method         string
	url            string
	body           io.Reader
	contentType    string
	okStatuses     []int
	extractPayload bool
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

func (c *Client) newRequest(ctx context.Context, opts requestOptions) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, opts.method, opts.url, opts.body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if opts.contentType != "" {
		req.Header.Set("Content-Type", opts.contentType)
	}
	return req, nil
}

// doRequest は HTTP リクエストを送信し、レスポンスボディと HTTP ステータスを返す。
func (c *Client) doRequest(req *http.Request) ([]byte, int, error) {
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

func isAllowedStatus(status int, allowed []int) bool {
	for _, candidate := range allowed {
		if status == candidate {
			return true
		}
	}
	return false
}

func jsonContentType(body io.Reader) string {
	if body == nil {
		return ""
	}
	return "application/json"
}

func (c *Client) doAPIRequest(ctx context.Context, opts requestOptions) ([]byte, error) {
	req, err := c.newRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	body, status, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if !isAllowedStatus(status, opts.okStatuses) {
		return nil, newAPIError(status, body)
	}
	if !opts.extractPayload {
		return body, nil
	}
	return extractPayload(body)
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

	return c.doAPIRequest(ctx, requestOptions{
		method:     http.MethodGet,
		url:        u.String(),
		okStatuses: []int{http.StatusOK},
	})
}

// GetByID は GET /api/v1/{path}/{id} でリソース単体を取得する。
func (c *Client) GetByID(ctx context.Context, path string, id int) ([]byte, error) {
	slog.Info("getting resource", "path", path, "id", id)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id)
	return c.doAPIRequest(ctx, requestOptions{
		method:     http.MethodGet,
		url:        urlStr,
		okStatuses: []int{http.StatusOK},
	})
}

// Create は POST /api/v1/{path} でリソースを作成する。
// Snipe-IT のレスポンスは {"status": "success", "payload": {...}} のため、payload を取り出して返す。
func (c *Client) Create(ctx context.Context, path string, data []byte) ([]byte, error) {
	slog.Info("creating resource", "path", path)
	return c.doAPIRequest(ctx, requestOptions{
		method:         http.MethodPost,
		url:            c.apiURL(path),
		body:           bytes.NewReader(data),
		contentType:    "application/json",
		okStatuses:     []int{http.StatusOK, http.StatusCreated},
		extractPayload: true,
	})
}

// Update は PATCH /api/v1/{path}/{id} でリソースを部分更新する。
// Snipe-IT のレスポンスは {"status": "success", "payload": {...}} のため、payload を取り出して返す。
func (c *Client) Update(ctx context.Context, path string, id int, data []byte) ([]byte, error) {
	slog.Info("updating resource", "path", path, "id", id)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id)
	return c.doAPIRequest(ctx, requestOptions{
		method:         http.MethodPatch,
		url:            urlStr,
		body:           bytes.NewReader(data),
		contentType:    "application/json",
		okStatuses:     []int{http.StatusOK},
		extractPayload: true,
	})
}

// Delete は DELETE /api/v1/{path}/{id} でリソースを削除する。
func (c *Client) Delete(ctx context.Context, path string, id int) error {
	slog.Info("deleting resource", "path", path, "id", id)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id)
	_, err := c.doAPIRequest(ctx, requestOptions{
		method:     http.MethodDelete,
		url:        urlStr,
		okStatuses: []int{http.StatusOK, http.StatusNoContent},
	})
	return err
}

// GetSub は GET /api/v1/{path}/{id}/{subPath} でサブリソースを取得する。
// 例: hardware/42/history, users/5/assets
func (c *Client) GetSub(ctx context.Context, path string, id int, subPath string) ([]byte, error) {
	slog.Info("getting sub-resource", "path", path, "id", id, "sub", subPath)
	urlStr := c.apiURL(path) + "/" + strconv.Itoa(id) + "/" + subPath
	return c.doAPIRequest(ctx, requestOptions{
		method:     http.MethodGet,
		url:        urlStr,
		okStatuses: []int{http.StatusOK},
	})
}

// GetByPath は GET /api/v1/{urlPath} で任意のパスを取得する。
// bytag/byserial/reports/account 等の非 CRUD パスに使用する。
func (c *Client) GetByPath(ctx context.Context, urlPath string) ([]byte, error) {
	slog.Info("getting by path", "path", urlPath)
	urlStr := c.apiURL(urlPath)
	return c.doAPIRequest(ctx, requestOptions{
		method:     http.MethodGet,
		url:        urlStr,
		okStatuses: []int{http.StatusOK},
	})
}

// PatchByPath は PATCH /api/v1/{urlPath} で任意のパスを更新する。
// ライセンスシート更新など、CRUD 汎用メソッドが対応しない入れ子パスに使用する。
func (c *Client) PatchByPath(ctx context.Context, urlPath string, data []byte) ([]byte, error) {
	slog.Info("patching by path", "path", urlPath)
	return c.doAPIRequest(ctx, requestOptions{
		method:         http.MethodPatch,
		url:            c.apiURL(urlPath),
		body:           bytes.NewReader(data),
		contentType:    "application/json",
		okStatuses:     []int{http.StatusOK},
		extractPayload: true,
	})
}

// PostByPath は POST /api/v1/{urlPath} で任意のパスにアクションを送信する。
// account/request 等の非 CRUD POST に使用する。data が nil の場合は空ボディで送信する。
func (c *Client) PostByPath(ctx context.Context, urlPath string, data []byte) ([]byte, error) {
	slog.Info("posting by path", "path", urlPath)
	var bodyReader io.Reader
	if data != nil {
		bodyReader = bytes.NewReader(data)
	}
	return c.doAPIRequest(ctx, requestOptions{
		method:         http.MethodPost,
		url:            c.apiURL(urlPath),
		body:           bodyReader,
		contentType:    jsonContentType(bodyReader),
		okStatuses:     []int{http.StatusOK, http.StatusCreated},
		extractPayload: true,
	})
}

// DeleteByPath は DELETE /api/v1/{urlPath} で任意のパスを削除する。
// account/personal-access-tokens/{id} 等の非 CRUD DELETE に使用する。
func (c *Client) DeleteByPath(ctx context.Context, urlPath string) error {
	slog.Info("deleting by path", "path", urlPath)
	_, err := c.doAPIRequest(ctx, requestOptions{
		method:     http.MethodDelete,
		url:        c.apiURL(urlPath),
		okStatuses: []int{http.StatusOK, http.StatusNoContent},
	})
	return err
}

// Upload は multipart/form-data で POST /api/v1/{urlPath} にファイルをアップロードする。
// Snipe-IT のインポート API（POST /api/v1/imports）に使用する。
// extraFields は追加フォームフィールド（例: {"import_type": "hardware"}）。
func (c *Client) Upload(ctx context.Context, urlPath, fieldName, filePath string, extraFields map[string]string) ([]byte, error) {
	slog.Info("uploading file", "path", urlPath, "file", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close() //nolint:errcheck

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	// ファイルフィールドを追加
	part, err := mw.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 追加フィールドを書き込む
	for k, v := range extraFields {
		if err := mw.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", k, err)
		}
	}
	mw.Close() //nolint:errcheck

	return c.doAPIRequest(ctx, requestOptions{
		method:         http.MethodPost,
		url:            c.apiURL(urlPath),
		body:           &buf,
		contentType:    mw.FormDataContentType(),
		okStatuses:     []int{http.StatusOK, http.StatusCreated},
		extractPayload: true,
	})
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
	return c.doAPIRequest(ctx, requestOptions{
		method:         http.MethodPost,
		url:            urlStr,
		body:           bodyReader,
		contentType:    jsonContentType(bodyReader),
		okStatuses:     []int{http.StatusOK, http.StatusCreated},
		extractPayload: true,
	})
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
