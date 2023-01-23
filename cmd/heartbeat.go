package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

const (
	MSG_HEARTBEAT_SUCCESS = "the connection to the Everactive API is healthy"
	MSG_HEARTBEAT_FAILURE = "something is wrong. Try enabling the debug option to get more information"
)

// heartbeatCmd represents the heartbeat command
var heartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Check the connection to the Everactive API",
	Long: `Verifies that the Everactive API is working and that the credentials are correct.
If there is a problem, the program will finish with and error code.`,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteHeartbeat()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initAPIClient()
	},
}

func init() {
	rootCmd.AddCommand(heartbeatCmd)
}

func ExecuteHeartbeat() {
	if ApiClient.Health() {
		Tui_info(MSG_HEARTBEAT_SUCCESS)
	} else {
		Tui_error(MSG_HEARTBEAT_FAILURE)
		os.Exit(1)
	}
}
