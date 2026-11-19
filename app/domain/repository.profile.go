package domain

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// ProfileRepositoryInterface reads and writes the configuration file: the
// saved profiles and the current one.
type ProfileRepositoryInterface interface {
	// Load returns the configuration; a missing file is an empty one, not
	// an error.
	Load() (*ConfigModel, error)
	// Save writes config, creating the directory, readable by the owner only.
	Save(config *ConfigModel) error
}

// ProfileRepository keeps the configuration as JSON in the file at Path.
type ProfileRepository struct {
	// Path is the configuration file; see [DefaultConfigPath].
	Path string
}

// DefaultConfigPath returns where the configuration lives: the
// WARP_ENGINE_CONFIG environment variable when set, otherwise
// warp-engine/config.json under the user's configuration directory
// (~/.config on Linux, ~/Library/Application Support on macOS, %AppData% on
// Windows).
func DefaultConfigPath() string {
	if path := os.Getenv("WARP_ENGINE_CONFIG"); path != "" {
		return path
	}
	directory, err := os.UserConfigDir()
	if err != nil {
		directory = "."
	}
	return filepath.Join(directory, "warp-engine", "config.json")
}

// Load reads the file at Path. A missing file returns an empty configuration;
// a file that is no valid JSON returns an error naming the file.
func (r ProfileRepository) Load() (*ConfigModel, error) {
	config := &ConfigModel{Profiles: map[string]ProfileModel{}}
	data, err := os.ReadFile(r.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return config, nil
	}
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, config); err != nil {
		slog.Debug("config file is not valid JSON", "path", r.Path, "error", err)
		return nil, &ConfigError{Path: r.Path, Err: err}
	}
	if config.Profiles == nil {
		config.Profiles = map[string]ProfileModel{}
	}
	return config, nil
}

// Save writes config to Path through a temporary file, so a crash leaves
// the old file whole. The file is created with mode 0600: it holds tokens.
func (r ProfileRepository) Save(config *ConfigModel) error {
	var err error
	if err = os.MkdirAll(filepath.Dir(r.Path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	temporary := r.Path + ".tmp"
	if err = os.WriteFile(temporary, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, r.Path)
}

// ConfigError is a configuration file that could not be read.
type ConfigError struct {
	// Path is the file.
	Path string
	// Err is why it could not be read.
	Err error
}

// Error names the file and the reason.
func (e *ConfigError) Error() string {
	return "the configuration file " + e.Path + " is broken: " + e.Err.Error()
}

// Unwrap returns the reason, for errors.Is and errors.As.
func (e *ConfigError) Unwrap() error {
	return e.Err
}
