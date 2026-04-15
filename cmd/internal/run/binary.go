// binary.go はバイナリレスポンスの保存ヘルパーを提供する。
// PDF 等、JSON パースを行わないバイナリレスポンスの取得・保存に使用する。
package run

import (
	"context"
	"fmt"
	"io"
	"os"
)

func saveBinaryToFile(path string, raw []byte) error {
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("failed to write to %s: %w", path, err)
	}
	return nil
}

func writeBinary(w io.Writer, raw []byte) error {
	_, err := w.Write(raw)
	return err
}

func reportBinarySave(w io.Writer, size int, path string) error {
	_, err := fmt.Fprintf(w, "Saved %d bytes to %s\n", size, path)
	return err
}

// writeBinaryOutput はバイナリデータを outputFile に保存するか、空の場合は w に書き出す。
// 出力先のルーティングのみを責務とし、HTTP 取得とは切り離す。
func writeBinaryOutput(w io.Writer, raw []byte, outputFile string) error {
	if outputFile == "" {
		return writeBinary(w, raw)
	}
	if err := saveBinaryToFile(outputFile, raw); err != nil {
		return err
	}
	return reportBinarySave(w, len(raw), outputFile)
}

// RunSaveBinary は GET /api/v1/{urlPath} のレスポンスをバイナリとして保存する。
// outputFile が空の場合は標準出力に書き出す（パイプ利用を想定）。
// JSON パースを行わないため、PDF 等のバイナリレスポンスに使用する。
func RunSaveBinary(ctx context.Context, o *BaseOptions, urlPath, outputFile string) error {
	raw, err := o.Client.GetByPath(ctx, urlPath)
	if err != nil {
		return err
	}
	return writeBinaryOutput(o.Stdout(), raw, outputFile)
}

// RunSaveBinaryBySegments は path segment を安全に連結してバイナリレスポンスを保存する。
func RunSaveBinaryBySegments(ctx context.Context, o *BaseOptions, outputFile string, segments ...string) error {
	return RunSaveBinary(ctx, o, JoinPathSegments(segments...), outputFile)
}
