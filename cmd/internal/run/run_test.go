package run_test

import (
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- ParseFilters ---

func TestParseFilters_Nil(t *testing.T) {
	result, err := run.ParseFilters(nil)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseFilters_Single(t *testing.T) {
	result, err := run.ParseFilters([]string{"status=2"})
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"status": {"2"}}, result)
}

func TestParseFilters_MultipleKeys(t *testing.T) {
	result, err := run.ParseFilters([]string{"status=2", "category_id=5"})
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"status":      {"2"},
		"category_id": {"5"},
	}, result)
}

func TestParseFilters_SameKeyMultipleTimes(t *testing.T) {
	result, err := run.ParseFilters([]string{"tag=prod", "tag=core"})
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"tag": {"prod", "core"}}, result)
}

func TestParseFilters_ValueWithEquals(t *testing.T) {
	// value に = が含まれる場合（base64 等）
	result, err := run.ParseFilters([]string{"cf_data=a=b=c"})
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"cf_data": {"a=b=c"}}, result)
}

func TestParseFilters_WhitespaceTrimmed(t *testing.T) {
	result, err := run.ParseFilters([]string{" status = 2 "})
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"status": {"2"}}, result)
}

func TestParseFilters_MissingEquals(t *testing.T) {
	_, err := run.ParseFilters([]string{"no-equals-sign"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --filter format")
}

// --- UnmarshalJSON ---

func TestUnmarshalJSON_ValidObject(t *testing.T) {
	v, err := run.UnmarshalJSON(`{"name":"Laptop","status_id":2}`)
	require.NoError(t, err)
	assert.Equal(t, "Laptop", v["name"])
}

func TestUnmarshalJSON_EmptyObject(t *testing.T) {
	v, err := run.UnmarshalJSON(`{}`)
	require.NoError(t, err)
	assert.Empty(t, v)
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	_, err := run.UnmarshalJSON(`{not valid json}`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestUnmarshalJSON_Array(t *testing.T) {
	// --data に配列を渡す場合は map 変換に失敗する
	_, err := run.UnmarshalJSON(`[1,2,3]`)
	assert.Error(t, err)
}

func TestValidateJSON_Array(t *testing.T) {
	err := run.ValidateJSON(`[1,2,3]`)
	assert.NoError(t, err)
}

func TestJSONBytes_Invalid(t *testing.T) {
	_, err := run.JSONBytes(`{not valid json}`)
	assert.Error(t, err)
}

func TestMarshalJSONData(t *testing.T) {
	data, err := run.MarshalJSONData(map[string]int{"fieldset_id": 3})
	require.NoError(t, err)
	assert.JSONEq(t, `{"fieldset_id":3}`, string(data))
}

// --- RequireDeleteConfirmation ---

func TestRequireDeleteConfirmation_WithYes(t *testing.T) {
	assert.NoError(t, run.RequireDeleteConfirmation(true))
}

func TestRequireDeleteConfirmation_WithoutYes(t *testing.T) {
	err := run.RequireDeleteConfirmation(false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--yes")
}

// --- FormatAPIError ---

func TestFormatAPIError_Nil(t *testing.T) {
	assert.NoError(t, run.FormatAPIError(nil))
}

func TestFormatAPIError_WrapsError(t *testing.T) {
	err := run.FormatAPIError(assert.AnError)
	assert.Error(t, err)
}

func TestRequirePositiveInt(t *testing.T) {
	assert.NoError(t, run.RequirePositiveInt("--id", 1))
	err := run.RequirePositiveInt("--id", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--id")
}

func TestRequireNonEmpty(t *testing.T) {
	assert.NoError(t, run.RequireNonEmpty("--tag", "asset-001"))
	err := run.RequireNonEmpty("--tag", " ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--tag")
}
