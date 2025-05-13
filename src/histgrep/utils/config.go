package utils

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DefaultLogs struct {
		Directory   string `toml:"directory"`
		FilePattern string `toml:"file_pattern"`
	} `toml:"default_logs"`
	Search struct {
		CaseSensitive bool   `toml:"case_sensitive"`
		DefaultName   string `toml:"default_name"`
	} `toml:"search"`
	Display struct {
		ColorEnabled bool `toml:"color_enabled"`
		PagerEnabled bool `toml:"pager_enabled"`
		VimExit      bool `toml:"vim_exit"`
	} `toml:"display"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{
		DefaultLogs: struct {
			Directory   string `toml:"directory"`
			FilePattern string `toml:"file_pattern"`
		}{
			Directory:   "~/.logs/",
			FilePattern: "EMPTY",
		},
		Search: struct {
			CaseSensitive bool   `toml:"case_sensitive"`
			DefaultName   string `toml:"default_name"`
		}{
			CaseSensitive: false,
			DefaultName:   "EMPTY",
		},
		Display: struct {
			ColorEnabled bool `toml:"color_enabled"`
			PagerEnabled bool `toml:"pager_enabled"`
			VimExit      bool `toml:"vim_exit"`
		}{
			ColorEnabled: true,
			PagerEnabled: false,
			VimExit:      false,
		},
	}

	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}

	// Expand ~ to home directory if present
	if config.DefaultLogs.Directory[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		config.DefaultLogs.Directory = filepath.Join(home, config.DefaultLogs.Directory[2:])
	}
	return config, nil
}

func GetMatchingLogFiles(config *Config) ([]string, error) {
	pattern := filepath.Join(config.DefaultLogs.Directory, config.DefaultLogs.FilePattern)
	files, err := filepath.Glob(pattern)
	return files, err
}
