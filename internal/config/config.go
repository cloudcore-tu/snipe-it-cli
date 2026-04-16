// config パッケージは snipe-it-cli の設定管理を担う。
//
// 設定ファイルは XDG Base Directory Specification に従い
// $XDG_CONFIG_HOME/snipe-it-cli/config.yaml（未設定時は ~/.config/snipe-it-cli/config.yaml）に置く。
//
// 優先順位: CLI フラグ > 環境変数 > 設定ファイル（選択インスタンス） > デフォルト値
package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"go.yaml.in/yaml/v3"
)

const (
	DefaultTimeout = 30
	DefaultOutput  = "table"
	configDirName  = "snipe-it-cli"
	configFileName = "config.yaml"
)

// FileConfig はディスク上の設定ファイル構造。
// 複数インスタンスを管理し、current で現在アクティブなインスタンスを指定する。
type FileConfig struct {
	Current   string              `yaml:"current"`
	Instances map[string]Instance `yaml:"instances"`
	Timeout   int                 `yaml:"timeout,omitempty"`
	Output    string              `yaml:"output,omitempty"`
}

// Instance は1つの Snipe-IT インスタンスへの接続情報を保持する。
type Instance struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

// Config は解決済みの設定。単一インスタンス用のフラットな構造体。
// Load() が FileConfig + 環境変数をマージして返す。
type Config struct {
	URL     string
	Token   string
	Timeout int
	Output  string
}

// ResolveProfile はアクティブインスタンス名を解決する。
// profile フラグ > SNIPE_PROFILE 環境変数 > fc.Current の優先順位で決定する。
// fc が nil の場合は環境変数と profile のみを参照する。
func ResolveProfile(fc *FileConfig, profile string) string {
	if profile != "" {
		return profile
	}
	if v := os.Getenv("SNIPE_PROFILE"); v != "" {
		return v
	}
	if fc != nil {
		return fc.Current
	}
	return ""
}

// InitFile は設定ファイルを新規作成する。
// すでに存在する場合はエラーを返す（上書き禁止）。
// 成功時は書き込まれたファイルのパスを返す。
func InitFile(fc *FileConfig) (string, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("config file already exists: %s (use 'snip config add' to add an instance)", path)
	}
	if err := WriteFile(fc); err != nil {
		return "", err
	}
	return path, nil
}

// UpsertInstance は設定ファイルのインスタンスを追加または更新する。
// ファイルが存在しない場合は新規作成する。
// current が未設定のときは name を current に設定する。
func UpsertInstance(name string, inst Instance) error {
	fc, err := ReadFile()
	if err != nil {
		return err
	}
	if fc == nil {
		fc = &FileConfig{
			Current:   name,
			Instances: make(map[string]Instance),
		}
	}
	if fc.Instances == nil {
		fc.Instances = make(map[string]Instance)
	}
	fc.Instances[name] = inst
	// current が未設定の場合は最初に追加したインスタンスをデフォルトにする
	if fc.Current == "" {
		fc.Current = name
	}
	return WriteFile(fc)
}

// Load は設定ファイル・環境変数から設定を読み込んで返す。
// profile が空の場合は設定ファイルの current インスタンスを使う。
// 優先順位: CLI フラグ（呼び出し元が適用） > 環境変数 > 設定ファイル > デフォルト値
// 副作用なし。パーミッション警告が必要な場合は呼び出し元で WarnInsecurePermissions を呼ぶ。
func Load(profile string) (*Config, error) {
	cfg := &Config{Timeout: DefaultTimeout, Output: DefaultOutput}

	fc, err := ReadFile()
	if err != nil {
		return nil, err
	}
	if fc != nil {
		applyFileConfig(cfg, fc)
		name := ResolveProfile(fc, profile)
		if name != "" {
			inst, ok := fc.Instances[name]
			if !ok {
				return nil, fmt.Errorf("instance %q not found in config file (check with snip config list)", name)
			}
			cfg.URL = inst.URL
			cfg.Token = inst.Token
		}
	}

	return cfg, applyEnvOverrides(cfg)
}

// applyFileConfig はファイル設定値を cfg に適用する。ゼロ値はデフォルトを上書きしない。
func applyFileConfig(cfg *Config, fc *FileConfig) {
	if fc.Timeout > 0 {
		cfg.Timeout = fc.Timeout
	}
	if fc.Output != "" {
		cfg.Output = fc.Output
	}
}

// WarnInsecurePermissions は設定ファイルのパーミッションを検査して警告する。
// 0600 以外の場合、API トークンを第三者が読める可能性があるため slog.Warn を出力する。
// 副作用（ログ出力）を持つため、CLI 境界（BaseOptions.Complete 等）から呼ぶ。
func WarnInsecurePermissions() {
	path, err := ConfigFilePath()
	if err != nil {
		return
	}
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		slog.Warn("config file has insecure permissions; recommend chmod 0600",
			"path", path, "permissions", fmt.Sprintf("%04o", perm))
	}
}

// applyEnvOverrides は SNIPEIT_* 環境変数で cfg を上書きする。設定ファイルより優先される。
func applyEnvOverrides(cfg *Config) error {
	if v := os.Getenv("SNIPEIT_URL"); v != "" {
		cfg.URL = v
	}
	if v := os.Getenv("SNIPEIT_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("SNIPEIT_TIMEOUT"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid SNIPEIT_TIMEOUT value (must be an integer): %w", err)
		}
		cfg.Timeout = t
	}
	if v := os.Getenv("SNIPEIT_OUTPUT"); v != "" {
		cfg.Output = v
	}
	return nil
}

// ReadFile は設定ファイルを読み込んで FileConfig を返す。
// ファイルが存在しない場合は nil を返す（エラーではない）。
// セキュリティ検査（パーミッション確認）は Load() 内の checkConfigPermissions() が担う。
func ReadFile() (*FileConfig, error) {
	path, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var fc FileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &fc, nil
}

// WriteFile は FileConfig を設定ファイルに書き込む。
func WriteFile(fc *FileConfig) error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(fc)
	if err != nil {
		return fmt.Errorf("failed to serialize config file: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ConfigDir は snipe-it-cli の設定ディレクトリのパスを返す。
// XDG Base Directory Specification に従い $XDG_CONFIG_HOME/snipe-it-cli を優先する。
// $XDG_CONFIG_HOME が未設定の場合は ~/.config/snipe-it-cli にフォールバックする。
// os.UserConfigDir() は macOS で ~/Library/Application Support を返すため使用しない。
func ConfigDir() (string, error) {
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		return filepath.Join(xdgHome, configDirName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", configDirName), nil
}

// ConfigFilePath は設定ファイルのフルパスを返す。
func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}
