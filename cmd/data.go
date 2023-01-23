package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Retrieve readings data from a sensor and a given time range.",
	Long: `The sensor MAC address is required. The time range can be set as:
duration: get the readings from a period of time until now. Example:  1h, 10m, 100s
last: get the last reading regardless of its date and time.
start-end: Unix timestamps range. Example: 1670507054-1670533006
The maximum range is 24 hours`,
	Run: func(cmd *cobra.Command, args []string) {
		sensorFlag := cmd.Flag("sensor")
		rangeParam := cmd.Flag("range")
		if sensorFlag != nil {
			if len(rangeParam.Value.String()) == 0 || rangeParam.Value.String() == "last" {
				executeDataLast(sensorFlag.Value.String())
			} else {
				start, end := CalculateRage(rangeParam.Value.String())
				if start == 0 || end == 0 {
					Tui_error("Invalid range. The maximum is 24 hours.")
					os.Exit(1)
				}
				executeDataWithRange(sensorFlag.Value.String(), start, end)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initAPIClient()
	},
}

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.Flags().String("sensor", "", "Required. Mac address of the sensor")
	_ = dataCmd.MarkFlagRequired("sensor")
	dataCmd.Flags().StringP("range", "r", "last", "Time range for the data: last (default), 1h, 10m, 100s, start-end timestamps. Max is 24h or 86400s")
}

func executeDataWithRange(sensorFilter string, start, end int64) {
	//DebugEnabled
	Tui_debug(fmt.Sprintf("get readings from MAC: %s - Range:  %d - %d", sensorFilter, start, end))

	response, err := ApiClient.GetSensorReadings(sensorFilter, start, end)
	if err != nil {
		Tui_error(fmt.Sprintf("Failed to retrieved sensors data: %s", err.Error()))
		os.Exit(1)
	}
	records := make([]string, 0)
	for _, record := range response.Data {
		jsonRecord, err := json.Marshal(record)
		if err != nil {
			Tui_error(fmt.Sprintf("error processing response %s", err.Error()))
			os.Exit(1)
		}
		records = append(records, string(jsonRecord))
	}
	for _, jsonRecord := range records {
		Tui_info(fmt.Sprintf("%s", jsonRecord))
	}
}

func executeDataLast(sensorFilter string) {
	Tui_debug(fmt.Sprintf("get last reading from MAC: %s ", sensorFilter))
	response, err := ApiClient.GetSensorLastReading(sensorFilter)
	if err != nil {
		Tui_error(fmt.Sprintf("Failed to retrieved sensors data: %s", err.Error()))
		os.Exit(1)
	}
	jsonRecord, err := json.Marshal(response.Data)
	Tui_info(fmt.Sprintf("%s", jsonRecord))
}

func CalculateRage(rangeParam string) (int64, int64) {
	start := int64(0)
	end := int64(0)
	limit := time.Hour * 24
	matched, err := regexp.Match(`\d+-\d+`, []byte(rangeParam))
	if err == nil && matched {
		toks := strings.Split(rangeParam, "-")
		start, _ = strconv.ParseInt(toks[0], 10, 64)
		end, _ = strconv.ParseInt(toks[1], 10, 64)
		return start, end
	}
	matched, err = regexp.Match(`\d+[hms]`, []byte(rangeParam))
	if err == nil && matched {
		offset, _ := time.ParseDuration(rangeParam)
		if offset.Seconds() > limit.Seconds() {
			return 0,0
		}
		endTime := time.Now().UTC()
		start = endTime.Add(-offset).Unix()
		Tui_debug(fmt.Sprintf("range %s - %s", endTime.Add(-offset), endTime))
		end = endTime.Unix()
	}

	return start, end
}
