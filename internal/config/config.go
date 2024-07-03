package config

import (
	"fmt"
	"kongtools/internal/pkg/log"
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
	defaultName   = ".kongtoolsrc"
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
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with the default name.
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(defaultName)

	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// Config file not found; create a new one
			createDefaultConfig()
		default:
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

// createDefaultConfig creates a default config file in the user's home directory
func createDefaultConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(home, defaultName)

	// Write the default config to the file
	cobra.CheckErr(os.WriteFile(configPath, []byte(DefaultConfig(log.DefaultConfig, view.DefaultConfig)), 0644))

	fmt.Fprintf(os.Stderr, "Created default config file: %s\n", configPath)
	viper.SetConfigFile(configPath)
	cobra.CheckErr(viper.ReadInConfig())
}
