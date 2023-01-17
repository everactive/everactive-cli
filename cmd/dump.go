package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/services"
	"os"
	"time"
)

var dumpDays int8

const dumpDelaySeconds = 1

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dumps eversensor data from the given number of days.",
	Long: `Retrieves the eversensor data in 24H intervals and up to 30 days. 
To avoid the hitting the rate limits it pauses between queries.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dumpDays < 1 || dumpDays > 30 {
			Tui_error("Invalid --days argument. The range of valid --days values is 1-30 ")
			os.Exit(1)
		}
		sensorFlag := cmd.Flag("sensor")
		if sensorFlag != nil {
			dumpData(sensorFlag.Value.String(), dumpDays)
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.Flags().String("sensor", "", "Required. Mac address of the sensor")
	_ = dumpCmd.MarkFlagRequired("sensor")
	dumpCmd.Flags().Int8Var(&dumpDays, "days", 1, "The number of days to dump. Defaults to 1")
}

func dumpData(sensorFilter string, days int8) {
	api := services.NewEveractiveAPIService(false)
	end := time.Now().UTC()
	start := end
	var d int8
	for d = 0; d < days; d++ {
		start = end.Add(-24 * time.Hour)
		Tui_debug(fmt.Sprintf("get readings from start %s end %s", start, end))

		response, err := api.GetSensorReadings(sensorFilter, start.Unix(), end.Unix())
		if err != nil {
			Tui_error(fmt.Sprintf("Failed to retrieved sensors data: %s", err.Error()))
			os.Exit(1)
		}
		for _, record := range response.Data {
			jsonRecord, err := json.Marshal(record)
			if err != nil {
				Tui_error(fmt.Sprintf("error processing response %s", err.Error()))
				os.Exit(1)
			}
			Tui_info(fmt.Sprintf("%s", jsonRecord))
		}
		end = start
		time.Sleep(time.Second * dumpDelaySeconds)
	}
}
