/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/services"
)

// listSensorsCmd represents the listSensors command
var listSensorsCmd = &cobra.Command{
	Use:   "list-sensors",
	Short: "Get a list of active Eversensors",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeListSensors()
	},
}

func init() {
	rootCmd.AddCommand(listSensorsCmd)
}

func executeListSensors() {
	api := services.NewEveractiveAPIService(DebugEnabled, context.Background())
	sensors,err:=api.GetSensorsList()
	if err != nil {
		TUI_Error(fmt.Sprintf("Failed to retrieved sensors: %s", err.Error()))
		return
	}
	TUI_Info(fmt.Sprintf("Total count: %d", sensors.PaginationInfo.TotalItems))
	for _, record := range sensors.Data {
		TUI_Info(fmt.Sprintf("%s - Type: %s - FW: %s - Association: %s %s",
			record.MacAddress, record.Type, record.LastInfo.FirmwareVersion,
			record.LastAssociation.GatewaySerialNumber, record.LastAssociation.Timestamp))
	}
}

