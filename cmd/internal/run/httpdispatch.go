// httpdispatch.go は HTTP 実行 + 出力の dispatch ヘルパーを提供する。
// 各 HTTP 動詞に対応する *AndPrint 関数と、
// path segment を安全に連結する JoinPathSegments を含む。
package run

import (
	"context"
	"net/url"
	"strings"
)

// JoinPathSegments は API パス segment を "/" で連結し、各 segment を個別に URL エスケープする。
// 可変入力は固定 segment と分けて渡すことで、path injection を防ぐ。
func JoinPathSegments(segments ...string) string {
	escaped := make([]string, len(segments))
	for i, segment := range segments {
		escaped[i] = url.PathEscape(strings.Trim(segment, "/"))
	}
	return strings.Join(escaped, "/")
}

// FetchAndPrint は GET /api/v1/{urlPath} を実行し結果を出力する。
// 各パッケージで文字列フラグからパスを組み立てるコマンドに使用する（bytag, byserial 等）。
func FetchAndPrint(ctx context.Context, o *BaseOptions, urlPath string) error {
	raw, err := o.client.GetByPath(ctx, urlPath)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}

// FetchBySegmentsAndPrint は path segment を安全に連結して GET /api/v1/{urlPath} を実行する。
func FetchBySegmentsAndPrint(ctx context.Context, o *BaseOptions, segments ...string) error {
	return FetchAndPrint(ctx, o, JoinPathSegments(segments...))
}

// PatchAndPrint は PATCH /api/v1/{urlPath} を実行し結果を出力する。
// ライセンスシート等の入れ子 PATCH エンドポイントに使用する。
func PatchAndPrint(ctx context.Context, o *BaseOptions, urlPath, data string) error {
	payload, err := ParseJSONBytes(data)
	if err != nil {
		return err
	}
	raw, err := o.client.PatchByPath(ctx, urlPath, payload)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}

// PatchBySegmentsAndPrint は path segment を安全に連結して PATCH を実行する。
func PatchBySegmentsAndPrint(ctx context.Context, o *BaseOptions, data string, segments ...string) error {
	return PatchAndPrint(ctx, o, JoinPathSegments(segments...), data)
}

// PostJSONAndPrint は JSON 文字列を検証して POST /api/v1/{urlPath} を実行し結果を出力する。
func PostJSONAndPrint(ctx context.Context, o *BaseOptions, urlPath, data string) error {
	payload, err := ParseJSONBytes(data)
	if err != nil {
		return err
	}
	return PostAndPrint(ctx, o, urlPath, payload)
}

// PostJSONBySegmentsAndPrint は path segment を安全に連結して POST を実行する。
func PostJSONBySegmentsAndPrint(ctx context.Context, o *BaseOptions, data string, segments ...string) error {
	return PostJSONAndPrint(ctx, o, JoinPathSegments(segments...), data)
}

// PostValueAndPrint は値を JSON 化して POST /api/v1/{urlPath} を実行し結果を出力する。
func PostValueAndPrint(ctx context.Context, o *BaseOptions, urlPath string, value any) error {
	payload, err := EncodeJSON(value)
	if err != nil {
		return err
	}
	return PostAndPrint(ctx, o, urlPath, payload)
}

// PostValueBySegmentsAndPrint は path segment を安全に連結して POST を実行する。
func PostValueBySegmentsAndPrint(ctx context.Context, o *BaseOptions, value any, segments ...string) error {
	return PostValueAndPrint(ctx, o, JoinPathSegments(segments...), value)
}

// PostAndPrint は POST /api/v1/{urlPath} を実行し結果を出力する。
// account/request 等の非 CRUD POST エンドポイントに使用する。
func PostAndPrint(ctx context.Context, o *BaseOptions, urlPath string, data []byte) error {
	raw, err := o.client.PostByPath(ctx, urlPath, data)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}

// PostBySegmentsAndPrint は path segment を安全に連結して POST を実行する。
func PostBySegmentsAndPrint(ctx context.Context, o *BaseOptions, data []byte, segments ...string) error {
	return PostAndPrint(ctx, o, JoinPathSegments(segments...), data)
}

// DeleteAndPrint は DELETE /api/v1/{urlPath} を実行し結果を出力する。
// account/personal-access-tokens/{id} 等の非 CRUD DELETE に使用する。
func DeleteAndPrint(ctx context.Context, o *BaseOptions, urlPath string) error {
	if err := o.client.DeleteByPath(ctx, urlPath); err != nil {
		return err
	}
	return o.PrintValue(map[string]any{"deleted": true})
}

// DeleteBySegmentsAndPrint は path segment を安全に連結して DELETE を実行する。
func DeleteBySegmentsAndPrint(ctx context.Context, o *BaseOptions, segments ...string) error {
	return DeleteAndPrint(ctx, o, JoinPathSegments(segments...))
}

// UploadAndPrint は multipart/form-data でファイルをアップロードし結果を出力する。
// Snipe-IT のインポート API に使用する。
func UploadAndPrint(ctx context.Context, o *BaseOptions, urlPath, fieldName, filePath string, extraFields map[string]string) error {
	raw, err := o.client.Upload(ctx, urlPath, fieldName, filePath, extraFields)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}
