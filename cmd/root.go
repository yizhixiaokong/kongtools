/*
Copyright © 2023 yizhixiaokong
*/
package cmd

import (
	"kongtools/internal/config"
	"kongtools/internal/tui"
	"os"

	"log/slog"

	"kongtools/internal/pkg/log"
	"kongtools/internal/pkg/sysinfo"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kongtools",
	Short: "kongtools is a command line tool for kong",
	Long:  `kongtools is a command line tool for kong`,
	Run:   rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	cfg, err := config.Config()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Debug("run app start ...")
	app := tui.NewApp(slog.Default(), cfg.App)

	if err := app.Run(); err != nil {
		slog.Error("run app error", slog.String("error", err.Error()))
		os.Exit(1)
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

func init() {
	cobra.OnInitialize(func() {
		cfg, err := config.Config() // init config
		if err != nil {
			slog.Error("failed to load config", slog.String("error", err.Error()))
			os.Exit(1)
		}
		log.InitLogger(cfg.Log)
		sysinfo.PrintSystemInfo() // print system info after all initialization
	})

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.kongtoolsrc)")
}
