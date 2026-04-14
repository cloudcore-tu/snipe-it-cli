package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/internal/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- NewPrinter / parseOutputFormat ---

func TestNewPrinter_Table(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "table"}
	p, err := pf.NewPrinter(&bytes.Buffer{})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewPrinter_JSON(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "json"}
	p, err := pf.NewPrinter(&bytes.Buffer{})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewPrinter_YAML(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "yaml"}
	p, err := pf.NewPrinter(&bytes.Buffer{})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewPrinter_CustomColumns(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "custom-columns=ID:.id,NAME:.name"}
	p, err := pf.NewPrinter(&bytes.Buffer{})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewPrinter_JSONPath(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "jsonpath={rows.#.id}"}
	p, err := pf.NewPrinter(&bytes.Buffer{})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewPrinter_UnknownFormat(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "xml"}
	_, err := pf.NewPrinter(&bytes.Buffer{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "xml")
}

// --- JSON 出力 ---

func TestPrint_JSON(t *testing.T) {
	buf := &bytes.Buffer{}
	pf := &output.PrintFlags{OutputFormat: "json"}
	p, err := pf.NewPrinter(buf)
	require.NoError(t, err)

	data := map[string]any{"id": 1, "name": "Laptop-001"}
	require.NoError(t, p.Print(data))

	assert.Contains(t, buf.String(), `"id"`)
	assert.Contains(t, buf.String(), `"Laptop-001"`)
}

// --- YAML 出力 ---

func TestPrint_YAML(t *testing.T) {
	buf := &bytes.Buffer{}
	pf := &output.PrintFlags{OutputFormat: "yaml"}
	p, err := pf.NewPrinter(buf)
	require.NoError(t, err)

	data := map[string]any{"id": 1, "name": "Laptop-001"}
	require.NoError(t, p.Print(data))

	out := buf.String()
	assert.Contains(t, out, "id:")
	assert.Contains(t, out, "Laptop-001")
}

// --- テーブル出力（Snipe-IT リストレスポンス） ---

func TestPrint_Table_WithRows(t *testing.T) {
	buf := &bytes.Buffer{}
	pf := &output.PrintFlags{OutputFormat: "table"}
	p, err := pf.NewPrinter(buf)
	require.NoError(t, err)

	// Snipe-IT のリストレスポンス形式
	data := map[string]any{
		"total": 1,
		"rows": []any{
			map[string]any{"id": 1, "name": "Laptop-001", "asset_tag": "ASSET-001"},
		},
	}
	require.NoError(t, p.Print(data))

	out := buf.String()
	assert.Contains(t, out, "ID")
	assert.Contains(t, out, "NAME")
	assert.Contains(t, out, "Laptop-001")
}

func TestPrint_Table_EmptyRows_FallsbackToJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	pf := &output.PrintFlags{OutputFormat: "table"}
	p, err := pf.NewPrinter(buf)
	require.NoError(t, err)

	data := map[string]any{"total": 0, "rows": []any{}}
	require.NoError(t, p.Print(data))

	// 行がない場合は JSON フォールバック
	assert.Contains(t, buf.String(), "total")
}

// --- custom-columns 出力 ---

func TestPrint_CustomColumns(t *testing.T) {
	buf := &bytes.Buffer{}
	pf := &output.PrintFlags{OutputFormat: "custom-columns=ID:.id,TAG:.asset_tag"}
	p, err := pf.NewPrinter(buf)
	require.NoError(t, err)

	data := map[string]any{
		"total": 1,
		"rows": []any{
			map[string]any{"id": 1, "asset_tag": "ASSET-001", "name": "Laptop"},
		},
	}
	require.NoError(t, p.Print(data))

	out := buf.String()
	assert.Contains(t, out, "ID")
	assert.Contains(t, out, "TAG")
	assert.Contains(t, out, "ASSET-001")
	// name は columns 指定にないので出力されない
	assert.NotContains(t, out, "Laptop")
}

func TestPrint_CustomColumns_InvalidFormat(t *testing.T) {
	pf := &output.PrintFlags{OutputFormat: "custom-columns=NOCORON"}
	_, err := pf.NewPrinter(&bytes.Buffer{})
	// custom-columns= 形式として解析される（エラーは Print 時）
	require.NoError(t, err)
	p, _ := pf.NewPrinter(&bytes.Buffer{})
	err = p.Print(map[string]any{})
	assert.Error(t, err)
}

// --- jsonpath 出力 ---

func TestPrint_JSONPath(t *testing.T) {
	buf := &bytes.Buffer{}
	pf := &output.PrintFlags{OutputFormat: "jsonpath={rows.#.id}"}
	p, err := pf.NewPrinter(buf)
	require.NoError(t, err)

	data := map[string]any{
		"total": 2,
		"rows": []any{
			map[string]any{"id": 1},
			map[string]any{"id": 2},
		},
	}
	require.NoError(t, p.Print(data))

	out := strings.TrimSpace(buf.String())
	assert.Equal(t, "[1,2]", out)
}

// --- PrintError ---

func TestPrintError(t *testing.T) {
	buf := &bytes.Buffer{}
	output.PrintError(buf, assert.AnError)
	assert.Contains(t, buf.String(), "Error:")
}
