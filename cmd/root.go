package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var DebugEnabled bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "everactive-cli",
	Short: "A tool to interact with the Everactive Edge Platform API",
	Long: `Set up your API credentials via the "credentials" command or as environment variables.
After that, you can list sensors, get readings, or stream data from the API.`,
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
	rootCmd.PersistentFlags().BoolVarP(&DebugEnabled, "debug", "d", false, "enable debug mode")
}


