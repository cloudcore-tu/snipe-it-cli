package run_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/cmd/internal/run"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testContextKey string

func newCompleteValidateRunCommand(ctx context.Context) *cobra.Command {
	root := &cobra.Command{Use: "snip"}
	root.PersistentFlags().String("profile", "", "")
	root.PersistentFlags().String("url", "", "")
	root.PersistentFlags().String("token", "", "")
	root.PersistentFlags().Int("timeout", 0, "")
	root.PersistentFlags().String("output", "", "")

	cmd := &cobra.Command{Use: "test"}
	cmd.SetContext(ctx)
	root.AddCommand(cmd)
	return cmd
}

func TestCompleteValidateRun_CallsValidateThenRun(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SNIPEIT_URL", "https://example.invalid")
	t.Setenv("SNIPEIT_TOKEN", "test-token")

	const key testContextKey = "key"

	cmd := newCompleteValidateRunCommand(context.WithValue(context.Background(), key, "value"))
	var calls []string

	err := run.CompleteValidateRun(cmd, &run.BaseOptions{}, func() error {
		calls = append(calls, "validate")
		return nil
	}, func(ctx context.Context) error {
		calls = append(calls, "run")
		assert.Equal(t, "value", ctx.Value(key))
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"validate", "run"}, calls)
}

func TestCompleteValidateRun_StopsBeforeRunWhenValidateFails(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SNIPEIT_URL", "https://example.invalid")
	t.Setenv("SNIPEIT_TOKEN", "test-token")

	cmd := newCompleteValidateRunCommand(context.Background())
	wantErr := errors.New("validate failed")
	runCalled := false

	err := run.CompleteValidateRun(cmd, &run.BaseOptions{}, func() error {
		return wantErr
	}, func(context.Context) error {
		runCalled = true
		return nil
	})

	require.ErrorIs(t, err, wantErr)
	assert.False(t, runCalled)
}

func TestCompleteValidateRun_StopsOnCompleteError(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SNIPEIT_URL", "")
	t.Setenv("SNIPEIT_TOKEN", "")

	cmd := newCompleteValidateRunCommand(context.Background())
	validateCalled := false
	runCalled := false

	err := run.CompleteValidateRun(cmd, &run.BaseOptions{}, func() error {
		validateCalled = true
		return nil
	}, func(context.Context) error {
		runCalled = true
		return nil
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Snipe-IT URL is not configured")
	assert.False(t, validateCalled)
	assert.False(t, runCalled)
}

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

func TestParseFilters_EmptyKey(t *testing.T) {
	_, err := run.ParseFilters([]string{"=value"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key must not be empty")
}

func TestParseFilters_EmptyValue(t *testing.T) {
	_, err := run.ParseFilters([]string{"status="})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must not be empty")
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

func TestRequirePositiveInt(t *testing.T) {
	assert.NoError(t, run.RequirePositiveInt("--id", 1))
	err := run.RequirePositiveInt("--id", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--id")
}

func TestRequireAll(t *testing.T) {
	wantErr := errors.New("boom")
	err := run.RequireAll(nil, wantErr, errors.New("ignored"))
	require.ErrorIs(t, err, wantErr)
}

func TestRequireNonEmpty(t *testing.T) {
	assert.NoError(t, run.RequireNonEmpty("--tag", "asset-001"))
	err := run.RequireNonEmpty("--tag", " ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--tag")
}

func TestRequireValidJSON(t *testing.T) {
	assert.NoError(t, run.RequireValidJSON("--data", `{"name":"Laptop"}`))

	err := run.RequireValidJSON("--data", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--data")

	err = run.RequireValidJSON("--data", "{invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestValidateOptionalJSON(t *testing.T) {
	assert.NoError(t, run.ValidateOptionalJSON(""))
	assert.NoError(t, run.ValidateOptionalJSON(" "))
	assert.NoError(t, run.ValidateOptionalJSON(`{"id":1}`))

	err := run.ValidateOptionalJSON("{invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestJoinPathSegments(t *testing.T) {
	assert.Equal(t,
		"labels/label%2Fa%20b",
		run.JoinPathSegments("labels", "label/a b"),
	)
}

func TestRequireFileExists(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "fixture-*.txt")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	assert.NoError(t, run.RequireFileExists("--file", file.Name()))

	err = run.RequireFileExists("--file", file.Name()+"-missing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to access --file")
}
