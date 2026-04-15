// helpers.go は HTTP 実行ヘルパーを提供する。
// GET/POST/PATCH/DELETE の各 HTTP 動詞に対応する Run* 関数と、
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

// RunGetByPath は初期化済みの BaseOptions を使って GET /api/v1/{urlPath} を実行し結果を出力する。
// 各パッケージで文字列フラグからパスを組み立てるコマンドに使用する（bytag, byserial 等）。
func RunGetByPath(ctx context.Context, o *BaseOptions, urlPath string) error {
	raw, err := o.Client.GetByPath(ctx, urlPath)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}

// RunGetBySegments は path segment を安全に連結して GET /api/v1/{urlPath} を実行する。
func RunGetBySegments(ctx context.Context, o *BaseOptions, segments ...string) error {
	return RunGetByPath(ctx, o, JoinPathSegments(segments...))
}

// RunPatchByPath は初期化済みの BaseOptions を使って PATCH /api/v1/{urlPath} を実行し結果を出力する。
// ライセンスシート等の入れ子 PATCH エンドポイントに使用する。
func RunPatchByPath(ctx context.Context, o *BaseOptions, urlPath, data string) error {
	payload, err := JSONBytes(data)
	if err != nil {
		return err
	}
	raw, err := o.Client.PatchByPath(ctx, urlPath, payload)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}

// RunPatchBySegments は path segment を安全に連結して PATCH /api/v1/{urlPath} を実行する。
func RunPatchBySegments(ctx context.Context, o *BaseOptions, data string, segments ...string) error {
	return RunPatchByPath(ctx, o, JoinPathSegments(segments...), data)
}

// RunPostJSONByPath は JSON 文字列を検証して POST /api/v1/{urlPath} を実行し結果を出力する。
func RunPostJSONByPath(ctx context.Context, o *BaseOptions, urlPath, data string) error {
	payload, err := JSONBytes(data)
	if err != nil {
		return err
	}
	return RunPostByPath(ctx, o, urlPath, payload)
}

// RunPostJSONBySegments は path segment を安全に連結して POST /api/v1/{urlPath} を実行する。
func RunPostJSONBySegments(ctx context.Context, o *BaseOptions, data string, segments ...string) error {
	return RunPostJSONByPath(ctx, o, JoinPathSegments(segments...), data)
}

// RunPostValueByPath は値を JSON 化して POST /api/v1/{urlPath} を実行し結果を出力する。
func RunPostValueByPath(ctx context.Context, o *BaseOptions, urlPath string, value any) error {
	payload, err := MarshalJSONData(value)
	if err != nil {
		return err
	}
	return RunPostByPath(ctx, o, urlPath, payload)
}

// RunPostValueBySegments は path segment を安全に連結して POST /api/v1/{urlPath} を実行する。
func RunPostValueBySegments(ctx context.Context, o *BaseOptions, value any, segments ...string) error {
	return RunPostValueByPath(ctx, o, JoinPathSegments(segments...), value)
}

// RunPostByPath は初期化済みの BaseOptions を使って POST /api/v1/{urlPath} を実行し結果を出力する。
// account/request 等の非 CRUD POST エンドポイントに使用する。
func RunPostByPath(ctx context.Context, o *BaseOptions, urlPath string, data []byte) error {
	raw, err := o.Client.PostByPath(ctx, urlPath, data)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}

// RunPostBySegments は path segment を安全に連結して POST /api/v1/{urlPath} を実行する。
func RunPostBySegments(ctx context.Context, o *BaseOptions, data []byte, segments ...string) error {
	return RunPostByPath(ctx, o, JoinPathSegments(segments...), data)
}

// RunDeleteByPath は初期化済みの BaseOptions を使って DELETE /api/v1/{urlPath} を実行する。
// account/personal-access-tokens/{id} 等の非 CRUD DELETE に使用する。
func RunDeleteByPath(ctx context.Context, o *BaseOptions, urlPath string) error {
	if err := o.Client.DeleteByPath(ctx, urlPath); err != nil {
		return err
	}
	return o.PrintValue(map[string]any{"deleted": true})
}

// RunDeleteBySegments は path segment を安全に連結して DELETE /api/v1/{urlPath} を実行する。
func RunDeleteBySegments(ctx context.Context, o *BaseOptions, segments ...string) error {
	return RunDeleteByPath(ctx, o, JoinPathSegments(segments...))
}

// RunUpload は multipart/form-data でファイルをアップロードし結果を出力する。
// Snipe-IT のインポート API に使用する。
func RunUpload(ctx context.Context, o *BaseOptions, urlPath, fieldName, filePath string, extraFields map[string]string) error {
	raw, err := o.Client.Upload(ctx, urlPath, fieldName, filePath, extraFields)
	if err != nil {
		return err
	}
	return o.PrintResponse(raw)
}
