// filedownload.go はバイナリレスポンスのファイル保存ヘルパーを提供する。
// PDF 等、JSON パースを行わないバイナリレスポンスの取得・保存に使用する。
package run

import (
	"context"
	"fmt"
	"io"
	"os"
)

// saveToFile はバイナリデータをファイルに書き込む。パーミッションは 0600。
func saveToFile(path string, raw []byte) error {
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("failed to write to %s: %w", path, err)
	}
	return nil
}

// writeRaw は生バイトを writer に書き出す。
func writeRaw(w io.Writer, raw []byte) error {
	_, err := w.Write(raw)
	return err
}

// printSaveResult は保存完了メッセージ（バイト数とパス）を出力する。
func printSaveResult(w io.Writer, size int, path string) error {
	_, err := fmt.Fprintf(w, "Saved %d bytes to %s\n", size, path)
	return err
}

// routeBinaryOutput はバイナリデータを outputFile に保存するか、空の場合は w に書き出す。
// 出力先のルーティングのみを責務とし、HTTP 取得とは切り離す。
func routeBinaryOutput(w io.Writer, raw []byte, outputFile string) error {
	if outputFile == "" {
		return writeRaw(w, raw)
	}
	if err := saveToFile(outputFile, raw); err != nil {
		return err
	}
	return printSaveResult(w, len(raw), outputFile)
}

// DownloadAndSave は GET /api/v1/{urlPath} のレスポンスをバイナリとして保存する。
// outputFile が空の場合は標準出力に書き出す（パイプ利用を想定）。
// JSON パースを行わないため、PDF 等のバイナリレスポンスに使用する。
func DownloadAndSave(ctx context.Context, o *BaseOptions, urlPath, outputFile string) error {
	raw, err := o.client.GetByPath(ctx, urlPath)
	if err != nil {
		return err
	}
	return routeBinaryOutput(o.Stdout(), raw, outputFile)
}

// DownloadBySegmentsAndSave は path segment を安全に連結してバイナリレスポンスを保存する。
func DownloadBySegmentsAndSave(ctx context.Context, o *BaseOptions, outputFile string, segments ...string) error {
	return DownloadAndSave(ctx, o, JoinPathSegments(segments...), outputFile)
}
