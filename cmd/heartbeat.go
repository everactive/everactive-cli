/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/services"
)

// heartbeatCmd represents the heartbeat command
var heartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Check the connection to the Everactive API",
	Long: `Verifies that the Everactive API is working and that the credentials are correct.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeHeartbeat()
	},
}

func init() {
	rootCmd.AddCommand(heartbeatCmd)
}

func executeHeartbeat() {
	api := services.NewEveractiveAPIService(DebugEnabled, context.Background())
	if api.Health() {
		TUI_Info("the connection to the Everactive API is healthy")
	} else {
		TUI_Error("something is wrong. Try enabling the debug option to get more information")
	}
}
