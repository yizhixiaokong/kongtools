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
	slog.Debug("run app start ...")
	app := tui.NewApp(slog.Default(), config.Config().App)

	if err := app.Run(); err != nil {
		slog.Error("run app error", slog.String("error", err.Error()))
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

func init() {
	cobra.OnInitialize(func() {
		cfg := config.Config() // init config
		log.InitLogger(cfg.Log)
		sysinfo.PrintSystemInfo() // print system info after all initialization
	})

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.kongtoolsrc)")
}
