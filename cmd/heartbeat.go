package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/services"
	"os"
)

// heartbeatCmd represents the heartbeat command
var heartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Check the connection to the Everactive API",
	Long: `Verifies that the Everactive API is working and that the credentials are correct.
If there is a problem, the program will finish with and error code.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeHeartbeat()
	},
}

func init() {
	rootCmd.AddCommand(heartbeatCmd)
}

func executeHeartbeat() {
	api := services.NewEveractiveAPIService(DebugEnabled)
	if api.Health() {
		Tui_info("the connection to the Everactive API is healthy")
	} else {
		Tui_error("something is wrong. Try enabling the debug option to get more information")
		os.Exit(1)
	}
}
