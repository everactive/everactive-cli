/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/services"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sensorFlag := cmd.Flag("sensor")
		rangeParam := cmd.Flag("range")
		outputFileParam := cmd.Flag("output-file")
		if sensorFlag != nil {
			if len(rangeParam.Value.String()) == 0 || rangeParam.Value.String() == "last" {
				executeDataLast(sensorFlag.Value.String(), outputFileParam.Value.String())
			} else {
				start, end := calculateRage(rangeParam.Value.String())
				if start == 0 || end == 0 {
					TUI_Error("invalid range")
					os.Exit(1)
				}
				executeDataWithRange(sensorFlag.Value.String(), start, end, outputFileParam.Value.String())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.Flags().String("sensor", "", "Required. Mac address of the sensor")
	_ = dataCmd.MarkFlagRequired("sensor")
	dataCmd.Flags().StringP("range", "r", "last", "Time range for the data: last (default), 1h, 10m, 100s, start-end timestamps. Max is 24h or 86400s")
	dataCmd.Flags().StringP("output-file", "o", "", "Save data to this file")
	dataCmd.Flags().StringP("format", "f", "json", "Output format: [json|csv]. Default is json")
}

func executeDataWithRange(sensorFilter string, start, end int64, outputFile string) {
	//DebugEnabled
	TUI_Debug(fmt.Sprintf("get readings from MAC: %s - Range:  %d - %d", sensorFilter, start, end))
	api := services.NewEveractiveAPIService(DebugEnabled, context.Background())
	response, err := api.GetSensorReadings(sensorFilter, start, end)
	if err != nil {
		TUI_Error(fmt.Sprintf("Failed to retrieved sensors: %s", err.Error()))
		return
	}
	records := make([]string, 0)
	for _, record := range response.Data {
		jsonRecord, err := json.Marshal(record)
		if err != nil {
			TUI_Error(fmt.Sprintf("error processing response %s", err.Error()))
			os.Exit(1)
		}
		records = append(records, string(jsonRecord))
	}
	if len(outputFile) == 0 {
		for _, jsonRecord := range records {
			TUI_Info(fmt.Sprintf("%s", jsonRecord))
		}
	} else {
		outputToFile(records, outputFile)
	}
}

func executeDataLast(sensorFilter string, outputFile string) {
	TUI_Debug(fmt.Sprintf("get last reading from MAC: %s ", sensorFilter))
	api := services.NewEveractiveAPIService(DebugEnabled, context.Background())
	response, err := api.GetSensorLastReading(sensorFilter)
	if err != nil {
		TUI_Error(fmt.Sprintf("Failed to retrieved sensors: %s", err.Error()))
		return
	}
	jsonRecord, err := json.Marshal(response.Data)
	if len(outputFile) == 0 {
		TUI_Info(fmt.Sprintf("%s", jsonRecord))
	} else {
		records := []string{string(jsonRecord)}
		outputToFile(records, outputFile)
	}
}

func calculateRage(rangeParam string) (int64, int64) {
	start := int64(0)
	end := int64(0)

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
		endTime := time.Now()
		start = endTime.Add(-offset).Unix()
		end = endTime.Unix()
	}

	return start, end
}

func outputToFile(records []string, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		TUI_Error(fmt.Sprintf("error creating output file: %s", err.Error()))
		os.Exit(1)
	}
	defer f.Close()
	for _, r := range records {
		_, err = f.WriteString(fmt.Sprintf("%s\n", r))
		if err != nil {
			TUI_Error(fmt.Sprintf("error writing to file: %s", err.Error()))
			os.Exit(1)
		}
	}

}
