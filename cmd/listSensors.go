package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/services"
	"os"
)

// listSensorsCmd represents the listSensors command
var listSensorsCmd = &cobra.Command{
	Use:   "list-sensors",
	Short: "Get a list of the Eversensors in your account",
	Long: `Retrieves and prints the list of the eversensors that are visible to this API account.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeListSensors()
	},
}

func init() {
	rootCmd.AddCommand(listSensorsCmd)
}

func executeListSensors() {
	api := services.NewEveractiveAPIService(DebugEnabled)
	sensors, err := api.GetSensorsList()
	if err != nil {
		Tui_error(fmt.Sprintf("Failed to retrieved sensors: %s", err.Error()))
		os.Exit(1)
	}
	Tui_info(fmt.Sprintf("Total count: %d", sensors.PaginationInfo.TotalItems))
	for _, record := range sensors.Data {
		Tui_info(fmt.Sprintf("Mac: %s - Type: %s - FW: %s - Association: %s %s",
			record.MacAddress, record.Type, record.LastInfo.FirmwareVersion,
			record.LastAssociation.GatewaySerialNumber, record.LastAssociation.Timestamp))
	}
}
