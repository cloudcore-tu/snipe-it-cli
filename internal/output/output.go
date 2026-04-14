// output パッケージは CLI の出力フォーマット切り替えを担う。
// kubectl の PrintFlags パターンに倣い、PrintFlags を構造体として共通化する。
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/tidwall/gjson"
	"go.yaml.in/yaml/v3"
)

// Format は出力フォーマットを表す。
type Format string

const (
	FormatJSON          Format = "json"
	FormatYAML          Format = "yaml"
	FormatTable         Format = "table"
	FormatCustomColumns Format = "custom-columns"
	FormatJSONPath      Format = "jsonpath"
)

// PrintFlags は --output フラグを構造体として保持する。各コマンドに埋め込む。
type PrintFlags struct {
	OutputFormat string
}

// NewPrintFlags はデフォルト値（table）で PrintFlags を生成する。
func NewPrintFlags() *PrintFlags {
	return &PrintFlags{OutputFormat: string(FormatTable)}
}

// Printer は指定フォーマットで値を出力する。
type Printer struct {
	format  Format
	options string
	writer  io.Writer
}

// NewPrinter は PrintFlags から Printer を生成する。
// "custom-columns=ID:.id,NAME:.name" や "jsonpath={rows.#.name}" の形式も受け付ける。
// jsonpath の式は gjson のパス構文を使用する。
func (f *PrintFlags) NewPrinter(w io.Writer) (*Printer, error) {
	format, options, err := parseOutputFormat(f.OutputFormat)
	if err != nil {
		return nil, err
	}
	return &Printer{format: format, options: options, writer: w}, nil
}

// parseOutputFormat は "--output" フラグの値をフォーマットとオプションに分解する。
func parseOutputFormat(raw string) (Format, string, error) {
	if strings.HasPrefix(raw, "custom-columns=") {
		return FormatCustomColumns, strings.TrimPrefix(raw, "custom-columns="), nil
	}
	if strings.HasPrefix(raw, "jsonpath=") {
		return FormatJSONPath, strings.TrimPrefix(raw, "jsonpath="), nil
	}
	switch Format(raw) {
	case FormatJSON, FormatYAML, FormatTable:
		return Format(raw), "", nil
	default:
		return "", "", fmt.Errorf(
			"unknown output format: %q\navailable: json, yaml, table, custom-columns=HEADER:PATH,..., jsonpath={gjson expression} (e.g. jsonpath={rows.#.id})",
			raw,
		)
	}
}

// Print は値を指定フォーマットで出力する。
func (p *Printer) Print(v any) error {
	switch p.format {
	case FormatJSON:
		return p.printJSON(v)
	case FormatYAML:
		return p.printYAML(v)
	case FormatTable:
		return p.printTable(v)
	case FormatCustomColumns:
		return p.printCustomColumns(v)
	case FormatJSONPath:
		return p.printJSONPath(v)
	default:
		return fmt.Errorf("unknown format: %s", p.format)
	}
}

func (p *Printer) printJSON(v any) error {
	enc := json.NewEncoder(p.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (p *Printer) printYAML(v any) error {
	enc := yaml.NewEncoder(p.writer)
	enc.SetIndent(2)
	return enc.Encode(v)
}

// columnDef は列定義を表す（ヘッダー名と gjson パスのペア）。
type columnDef struct {
	header string
	path   string
}

// tableColumnCandidates は --output table 時に表示を試みる列の候補。
// Snipe-IT API のフィールド名に対応している。
// 優先順位の高い順に定義し、最初のアイテムで値が存在する列のみ表示する。
var tableColumnCandidates = []columnDef{
	{header: "ID", path: "id"},
	{header: "NAME", path: "name"},
	{header: "ASSET TAG", path: "asset_tag"},
	{header: "SERIAL", path: "serial"},
	{header: "MODEL", path: "model.name"},
	{header: "CATEGORY", path: "category.name"},
	{header: "STATUS", path: "status_label.name"},
	{header: "LOCATION", path: "location.name"},
	{header: "ASSIGNED TO", path: "assigned_to.name"},
	{header: "EMAIL", path: "email"},
	{header: "DEPARTMENT", path: "department.name"},
	{header: "SEATS", path: "seats"},
	{header: "AVAILABLE SEATS", path: "available_actions.checkout"},
	{header: "MANUFACTURER", path: "manufacturer.name"},
	{header: "TYPE", path: "category_type"},
	{header: "DESCRIPTION", path: "notes"},
}

// extractTableItems は JSON バイト列からテーブル描画用のアイテムスライスを取り出す。
// Snipe-IT のリストレスポンス（rows 配列）と単一オブジェクトの両方に対応する。
func extractTableItems(jsonBytes []byte) []gjson.Result {
	// Snipe-IT は list レスポンスに "rows" を使う（NetBox の "results" とは異なる）
	rows := gjson.GetBytes(jsonBytes, "rows")
	if rows.IsArray() {
		var items []gjson.Result
		rows.ForEach(func(_, v gjson.Result) bool {
			items = append(items, v)
			return true
		})
		return items
	}
	// rows 配列がなければトップレベル全体を単一アイテムとして扱う
	return []gjson.Result{gjson.ParseBytes(jsonBytes)}
}

// detectTableColumns は最初のアイテムから表示する列を決定する。
func detectTableColumns(firstItem gjson.Result) []columnDef {
	var cols []columnDef
	for _, c := range tableColumnCandidates {
		val := firstItem.Get(c.path)
		if val.Exists() && val.Type != gjson.Null {
			cols = append(cols, c)
		}
	}
	return cols
}

// tableCellValue はテーブルセルに表示する文字列を返す。
func tableCellValue(item gjson.Result, path string) string {
	val := item.Get(path)
	if !val.Exists() || val.Type == gjson.Null {
		return "<nil>"
	}
	text := val.String()
	if text == "" {
		return "<nil>"
	}
	return text
}

// printTable は値をテーブル形式で出力する。
// 最初のアイテムを調べて表示する列を自動検出し、tabwriter で整形する。
// 既知フィールドが見つからない場合は JSON にフォールバックする。
func (p *Printer) printTable(v any) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	items := extractTableItems(jsonBytes)
	if len(items) == 0 {
		return p.printJSON(v)
	}

	activeCols := detectTableColumns(items[0])
	if len(activeCols) == 0 {
		return p.printJSON(v)
	}

	w := tabwriter.NewWriter(p.writer, 0, 0, 2, ' ', 0)

	headers := make([]string, len(activeCols))
	for i, c := range activeCols {
		headers[i] = c.header
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, item := range items {
		vals := make([]string, len(activeCols))
		for i, c := range activeCols {
			vals[i] = tableCellValue(item, c.path)
		}
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	}

	return w.Flush()
}

// printCustomColumns は "HEADER:PATH,HEADER:PATH" 形式で列を定義して table 出力する。
func (p *Printer) printCustomColumns(v any) error {
	cols, err := parseColumnDefs(p.options)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	w := tabwriter.NewWriter(p.writer, 0, 0, 2, ' ', 0)

	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.header
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, item := range extractTableItems(jsonBytes) {
		vals := make([]string, len(cols))
		for i, c := range cols {
			vals[i] = tableCellValue(item, c.path)
		}
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	}

	return w.Flush()
}

// printJSONPath は gjson のパス式で値を抽出して出力する。
// {expr} 形式の外側の {} と先頭の . は互換のために除去して評価する。
func (p *Printer) printJSONPath(v any) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	expr := strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(p.options, "{"), "}"), ".")

	result := gjson.GetBytes(jsonBytes, expr)
	fmt.Fprintln(p.writer, result.String())
	return nil
}

// parseColumnDefs は "ID:.id,NAME:.name" 形式の文字列を columnDef スライスに変換する。
func parseColumnDefs(raw string) ([]columnDef, error) {
	parts := strings.Split(raw, ",")
	cols := make([]columnDef, 0, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid custom-columns format: %q (expected: HEADER:PATH)", part)
		}
		cols = append(cols, columnDef{
			header: strings.TrimSpace(kv[0]),
			path:   strings.TrimPrefix(strings.TrimSpace(kv[1]), "."),
		})
	}
	return cols, nil
}

// PrintError はエラーメッセージを stderr にプレーンテキストで出力する。
func PrintError(w io.Writer, err error) {
	fmt.Fprintf(w, "Error: %s\n", err)
}
