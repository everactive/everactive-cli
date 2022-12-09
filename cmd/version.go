package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/lib"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the application",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		TUI_Info(lib.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
