package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"kongtools/internal/pkg/log"
	"kongtools/internal/pkg/paths"
	"kongtools/internal/tui"
)

const (
	appName    = "kongtools"
	configName = "config"
)

var (
	// CfgFile can be set by flags to specify the config file
	CfgFile string
	_config config
	once    sync.Once
)

type config struct {
	Log log.Config `mapstructure:"log"`
	App tui.Config `mapstructure:"app"`
}

func Config() (config, error) {
	var err error
	once.Do(func() {
		err = initConfig()
	})
	return _config, err
}

// initConfig reads in config file and ENV variables if set.
func initConfig() error {
	var configFilePath string

	if CfgFile != "" {
		// Use config file from the flag.
		configFilePath = CfgFile
		// Ensure the directory exists for custom config file
		if dir := filepath.Dir(configFilePath); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create config directory %s: %w", dir, err)
			}
		}
	} else {
		// Get config file path using paths package
		configFilePath = paths.ConfigFile(configName + ".yaml")
		CfgFile = configFilePath
	}

	viper.SetConfigFile(configFilePath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// Check if error is because file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
			// Config file not found; create a new one
			if err := createDefaultConfig(configFilePath); err != nil {
				return err
			}
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Print the config file used
	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// read in environment variables that match
	viper.AutomaticEnv()
	// Unmarshal the config into a struct
	if err := viper.Unmarshal(&_config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证并填充默认值
	if err := _config.App.Validate(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default config file using yaml.Marshal
func createDefaultConfig(configPath string) error {
	// Build default config from all sub-packages
	defaultCfg := config{
		Log: *log.Default(),
		App: *tui.Default(),
	}

	// Marshal to YAML with 2-space indentation
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2) // Use 2-space indentation
	if err := encoder.Encode(defaultCfg); err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close encoder: %w", err)
	}

	// Prepend header
	defaultContent := "# KongTools Configuration\n\n" + buf.String()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(defaultContent), 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Created default config file: %s\n", configPath)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read default config: %w", err)
	}
	return nil
}
