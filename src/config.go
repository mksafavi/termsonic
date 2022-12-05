package src

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BaseURL  string
	Username string
	Password string
}

func LoadConfigFromFile(filename string) (*Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(filename, &cfg)

	return &cfg, err
}

func LoadDefaultConfig() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if os.IsNotExist(err) {
		return &Config{}, nil
	}
	f.Close()

	return LoadConfigFromFile(path)
}

func getConfigFilePath() (string, error) {
	path := ""
	if runtime.GOOS == "linux" {
		configDir := os.Getenv("XDG_CONFIG_DIR")
		if configDir == "" {
			home := os.Getenv("HOME")
			if home == "" {
				return "", fmt.Errorf("could not determine where to store configuration")
			}

			path = filepath.Join(home, ".config")
			os.MkdirAll(path, os.ModeDir.Perm())

			path = filepath.Join(path, "termsonic.toml")
		} else {
			path = filepath.Join(configDir, "termsonic.toml")
		}
	} else if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return "", fmt.Errorf("could not find %%APPDATA%%")
		}

		path = filepath.Join(appdata, "Termsonic")
		os.MkdirAll(path, os.ModeDir.Perm())

		path = filepath.Join(path, "termsonic.toml")
	} else {
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return path, nil
}

func (c *Config) Save() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	return enc.Encode(*c)
}