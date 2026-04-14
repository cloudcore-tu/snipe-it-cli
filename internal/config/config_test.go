package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tempConfigDir はテスト用の一時設定ディレクトリを作成し、XDG_CONFIG_HOME を向ける。
func tempConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	return dir
}

// --- Load ---

func TestLoad_NoFile_ReturnsDefaults(t *testing.T) {
	tempConfigDir(t)
	// 環境変数をクリア
	t.Setenv("SNIPEIT_URL", "")
	t.Setenv("SNIPEIT_TOKEN", "")

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Equal(t, config.DefaultTimeout, cfg.Timeout)
	assert.Equal(t, config.DefaultOutput, cfg.Output)
	assert.Empty(t, cfg.URL)
	assert.Empty(t, cfg.Token)
}

func TestLoad_EnvOverridesDefaults(t *testing.T) {
	tempConfigDir(t)
	t.Setenv("SNIPEIT_URL", "https://snip.example.com")
	t.Setenv("SNIPEIT_TOKEN", "mytoken")
	t.Setenv("SNIPEIT_TIMEOUT", "60")
	t.Setenv("SNIPEIT_OUTPUT", "json")

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Equal(t, "https://snip.example.com", cfg.URL)
	assert.Equal(t, "mytoken", cfg.Token)
	assert.Equal(t, 60, cfg.Timeout)
	assert.Equal(t, "json", cfg.Output)
}

func TestLoad_InvalidTimeout_ReturnsError(t *testing.T) {
	tempConfigDir(t)
	t.Setenv("SNIPEIT_TIMEOUT", "not-a-number")

	_, err := config.Load("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SNIPEIT_TIMEOUT")
}

func TestLoad_ProfileFromEnv(t *testing.T) {
	dir := tempConfigDir(t)
	t.Setenv("SNIPEIT_URL", "")
	t.Setenv("SNIPEIT_TOKEN", "")
	t.Setenv("SNIPE_PROFILE", "staging")

	fc := &config.FileConfig{
		Current: "prod",
		Instances: map[string]config.Instance{
			"prod":    {URL: "https://prod.example.com", Token: "prod-token"},
			"staging": {URL: "https://staging.example.com", Token: "stg-token"},
		},
	}
	writeTestConfig(t, dir, fc)

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Equal(t, "https://staging.example.com", cfg.URL)
	assert.Equal(t, "stg-token", cfg.Token)
}

func TestLoad_ProfileFlagOverridesEnv(t *testing.T) {
	dir := tempConfigDir(t)
	t.Setenv("SNIPEIT_URL", "")
	t.Setenv("SNIPEIT_TOKEN", "")
	t.Setenv("SNIPE_PROFILE", "staging")

	fc := &config.FileConfig{
		Current: "prod",
		Instances: map[string]config.Instance{
			"prod":    {URL: "https://prod.example.com", Token: "prod-token"},
			"staging": {URL: "https://staging.example.com", Token: "stg-token"},
		},
	}
	writeTestConfig(t, dir, fc)

	// --profile フラグで prod を指定 → SNIPE_PROFILE=staging より優先される
	cfg, err := config.Load("prod")
	require.NoError(t, err)
	assert.Equal(t, "https://prod.example.com", cfg.URL)
}

func TestLoad_InvalidProfile_ReturnsError(t *testing.T) {
	dir := tempConfigDir(t)
	t.Setenv("SNIPEIT_URL", "")
	t.Setenv("SNIPEIT_TOKEN", "")
	t.Setenv("SNIPE_PROFILE", "")

	fc := &config.FileConfig{
		Current:   "prod",
		Instances: map[string]config.Instance{"prod": {URL: "https://prod.example.com", Token: "token"}},
	}
	writeTestConfig(t, dir, fc)

	_, err := config.Load("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}

// --- ReadFile / WriteFile ---

func TestWriteFile_And_ReadFile(t *testing.T) {
	dir := tempConfigDir(t)

	fc := &config.FileConfig{
		Current: "prod",
		Instances: map[string]config.Instance{
			"prod": {URL: "https://prod.example.com", Token: "prod-token"},
		},
		Timeout: 45,
		Output:  "json",
	}
	require.NoError(t, config.WriteFile(fc))

	got, err := config.ReadFile()
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "prod", got.Current)
	assert.Equal(t, "https://prod.example.com", got.Instances["prod"].URL)
	assert.Equal(t, 45, got.Timeout)

	// セキュリティ: WriteFile は 0600 で書き込むこと
	path := filepath.Join(dir, "snipe-it-cli", "config.yaml")
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestReadFile_NotExist_ReturnsNil(t *testing.T) {
	tempConfigDir(t)

	fc, err := config.ReadFile()
	require.NoError(t, err)
	assert.Nil(t, fc)
}

// --- ConfigDir ---

func TestConfigDir_WithXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-test")
	dir, err := config.ConfigDir()
	require.NoError(t, err)
	assert.Equal(t, "/tmp/xdg-test/snipe-it-cli", dir)
}

func TestConfigDir_WithoutXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	dir, err := config.ConfigDir()
	require.NoError(t, err)
	assert.Contains(t, dir, "snipe-it-cli")
	// macOS の ~/Library/Application Support は使わず ~/.config を使うことを確認する
	assert.NotContains(t, dir, "Application Support")
}

// writeTestConfig はテスト用設定ファイルを作成するヘルパー。
func writeTestConfig(t *testing.T, xdgDir string, fc *config.FileConfig) {
	t.Helper()
	require.NoError(t, config.WriteFile(fc))
	// テスト後はファイルを自動削除（t.TempDir が処理）
	_ = xdgDir
}
