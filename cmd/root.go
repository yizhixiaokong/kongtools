/*
Copyright © 2023 yizhixiaokong
*/
package cmd

import (
	"fmt"
	"kongtools/internal/view"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kongtools",
	Short: "kongtools is a command line tool for kong",
	Long:  `kongtools is a command line tool for kong`,
	Run:   rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	InitLogger()

	slog.Debug("run app start ...")
	app := view.NewApp(slog.Default())
	if err := app.Init(); err != nil {
		slog.Error("init app error", err)
		return
	}

	if err := app.Run(); err != nil {
		slog.Error("run app error", err)
		return
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// InitLogger 初始化日志
func InitLogger() { // TODO: new package
	ll := &lumberjack.Logger{
		Filename:   "log/log.log", // filename // TODO: set from config
		MaxSize:    10,            // megabytes // TODO: set from config
		MaxAge:     30,            // days // TODO: set from config
		MaxBackups: 15,            // max backups // TODO: set from config
	}
	err := ll.Rotate() // 每次启动程序都会归档之前的日志
	if err != nil {
		panic(err)
	}

	// w := io.MultiWriter(os.Stdout, ll) // 同时写文件和屏幕 // !不需要

	logHandler := slog.NewJSONHandler(ll, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // TODO: set level from config
		AddSource: true,            // TODO: set source from config
	})
	slog.SetDefault(slog.New(logHandler))
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kongtoolsrc)")
}

type Config struct {
	Key1 string
	Key2 string
}

var config Config

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
	if err := viper.Unmarshal(&config); err != nil {
		cobra.CheckErr(err)
	}

	fmt.Printf("Config: %+v\n", config) // TODO: for test, delete
}

// createDefaultConfig creates a default config file
func createDefaultConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(home, ".kongtoolsrc")

	// Define default configuration
	// TODO: set default configuration
	defaultConfig := []byte(`# Default configuration
key1: value1
key2: value2
`)

	// Write the default config to the file
	err = os.WriteFile(configPath, defaultConfig, 0644)
	cobra.CheckErr(err)

	fmt.Fprintf(os.Stderr, "Created default config file: %s\n", configPath)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}
}
