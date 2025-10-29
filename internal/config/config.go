package config

import (
	"fmt"
	"kongtools/internal/pkg/log"
	"kongtools/internal/pkg/paths"
	"kongtools/internal/view"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// Define default configuration
	defaultCfgStr = `# Default configuration`
	appName       = "kongtools"
	configName    = "config"
)

var (
	// CfgFile can be set by flags to specify the config file
	CfgFile string
	_config config
	once    sync.Once
)

type config struct {
	Log log.Config
	App view.Config
}

func Config() config {
	once.Do(func() {
		initConfig()
	})
	return _config
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configFilePath string

	if CfgFile != "" {
		// Use config file from the flag.
		configFilePath = CfgFile
		// Ensure the directory exists for custom config file
		if dir := filepath.Dir(configFilePath); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				cobra.CheckErr(fmt.Errorf("failed to create config directory %s: %w", dir, err))
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
			createDefaultConfig(configFilePath)
		} else {
			// Config file was found but another error was produced
			cobra.CheckErr(err)
		}
	}

	// Print the config file used
	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// read in environment variables that match
	viper.AutomaticEnv()
	// Unmarshal the config into a struct
	cobra.CheckErr(viper.Unmarshal(&_config))

	// fmt.Printf("Config: %+v", _config) // Debug
}

// DefaultConfig returns the default configuration as a string
func DefaultConfig(configs ...string) (config string) {
	config = defaultCfgStr
	for _, cfg := range configs {
		config += "\n" + cfg
	}
	return
}

// createDefaultConfig creates a default config file in the config directory
func createDefaultConfig(configPath string) {
	// Directory is already ensured by paths.ConfigFile
	// Write the default config to the file
	cobra.CheckErr(os.WriteFile(configPath, []byte(DefaultConfig(log.DefaultConfig, view.DefaultConfig)), 0644))

	fmt.Fprintf(os.Stderr, "Created default config file: %s\n", configPath)
	viper.SetConfigFile(configPath)
	cobra.CheckErr(viper.ReadInConfig())
}
