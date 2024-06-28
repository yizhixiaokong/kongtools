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

type config struct {
	Log log.Config
	App view.Config
}

var (
	CfgFile string
	_config config
	once    sync.Once
)

func Config() config {
	once.Do(func() {
		initConfig()
	})
	return _config
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
		viper.SetConfigType("yaml")
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kongtoolsrc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kongtoolsrc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; create a new one
		createDefaultConfig()
	} else {
		// Config file was found but another error was produced
		cobra.CheckErr(err)
	}

	// Unmarshal the config into a struct
	cobra.CheckErr(viper.Unmarshal(&_config))

	fmt.Printf("Config: %+v", _config)
}

// Define default configuration
const defaultConfig = `# Default configuration`

func DefaultConfig(configs ...string) (config string) {
	config = defaultConfig
	for _, cfg := range configs {
		config += "\n" + cfg
	}
	return
}

// createDefaultConfig creates a default config file
func createDefaultConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(home, ".kongtoolsrc")

	// Write the default config to the file
	cobra.CheckErr(os.WriteFile(configPath, []byte(DefaultConfig(log.DefaultConfig, view.DefaultConfig)), 0644))

	fmt.Fprintf(os.Stderr, "Created default config file: %s\n", configPath)
	viper.SetConfigFile(configPath)
	cobra.CheckErr(viper.ReadInConfig())
	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}
