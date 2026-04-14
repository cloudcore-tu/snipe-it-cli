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

// Load は設定ファイル・環境変数から設定を読み込んで返す。
// profile が空の場合は設定ファイルの current インスタンスを使う。
func Load(profile string) (*Config, error) {
	cfg := &Config{
		Timeout: DefaultTimeout,
		Output:  DefaultOutput,
	}

	// 設定ファイルを読み込む（存在しない場合はスキップ）
	fc, err := ReadFile()
	if err != nil {
		return nil, err
	}
	if fc != nil {
		if fc.Timeout > 0 {
			cfg.Timeout = fc.Timeout
		}
		if fc.Output != "" {
			cfg.Output = fc.Output
		}

		// アクティブなインスタンスを解決:
		// --profile フラグ > SNIPE_PROFILE 環境変数 > 設定ファイルの current
		name := profile
		if name == "" {
			name = os.Getenv("SNIPE_PROFILE")
		}
		if name == "" {
			name = fc.Current
		}
		if name != "" {
			inst, ok := fc.Instances[name]
			if !ok {
				return nil, fmt.Errorf("instance %q not found in config file (check with snip config list)", name)
			}
			cfg.URL = inst.URL
			cfg.Token = inst.Token
		}
	}

	// 環境変数で上書き（設定ファイルより優先）
	if v := os.Getenv("SNIPEIT_URL"); v != "" {
		cfg.URL = v
	}
	if v := os.Getenv("SNIPEIT_TOKEN"); v != "" {
		cfg.Token = v
	}
	if v := os.Getenv("SNIPEIT_TIMEOUT"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid SNIPEIT_TIMEOUT value (must be an integer): %w", err)
		}
		cfg.Timeout = t
	}
	if v := os.Getenv("SNIPEIT_OUTPUT"); v != "" {
		cfg.Output = v
	}

	return cfg, nil
}

// ReadFile は設定ファイルを読み込んで FileConfig を返す。
// ファイルが存在しない場合は nil を返す（エラーではない）。
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

	// セキュリティ: 設定ファイルには API トークンが含まれる。
	// 0600 以外のパーミッションは他のユーザーがトークンを読める可能性があるため警告する。
	if info, statErr := os.Stat(path); statErr == nil {
		if perm := info.Mode().Perm(); perm != 0o600 {
			slog.Warn("config file has insecure permissions; recommend chmod 0600",
				"path", path, "permissions", fmt.Sprintf("%04o", perm))
		}
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
