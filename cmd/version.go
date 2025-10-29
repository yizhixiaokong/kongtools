package cmd

import (
	"fmt"
	"kongtools/internal/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information of kongtools`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetFullVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
