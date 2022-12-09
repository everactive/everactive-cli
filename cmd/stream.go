/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/lib"
	"gitlab.com/everactive/everactive-cli/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var lastTimestamp = int64(0)

const (
	loopPeriodMinutes = 10
	loopDelaySeconds  = 5
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sensorFlag := cmd.Flag("sensor")
		//outputFileParam := cmd.Flag("output-file")
		if sensorFlag != nil {
			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
			TUI_Debug("press ctrl+c to stop...")
			go streamLoop(sensorFlag.Value.String())
			<-done // Will block here until user hits ctrl+c
			TUI_Debug("good bye")
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.Flags().String("sensor", "", "Required. Mac address of the sensor")
	_ = streamCmd.MarkFlagRequired("sensor")
}

func streamLoop(sensorFilter string) {
	api := services.NewEveractiveAPIService(false, context.Background())
	for {
		endTime := time.Now()
		start := endTime.Add(time.Minute * -loopPeriodMinutes).Unix()
		end := endTime.Unix()
		TUI_Debug(fmt.Sprintf("readings from start %d end %d", start, end))
		time.Sleep(time.Second * loopDelaySeconds)
		response, err := api.GetSensorReadings(sensorFilter, start, end)
		if err != nil {
			TUI_Error(fmt.Sprintf("failed to retrieved sensors data: %s", err.Error()))
			return
		}
		var reading lib.SensorReading
		for _, record := range response.Data {
			jsonRecord, err := json.Marshal(record)
			if err != nil {
				TUI_Error(fmt.Sprintf("error processing response %s", err.Error()))
				os.Exit(1)
			}
			err = json.Unmarshal(jsonRecord, &reading)
			if err != nil {
				TUI_Error(fmt.Sprintf("error processing response %s", err.Error()))
				os.Exit(1)
			}
			if reading.Timestamp > lastTimestamp {
				lastTimestamp = reading.Timestamp
				TUI_Debug(fmt.Sprintf("pkt %d - timestamp: %s", reading.PacketNumberGateway, reading.ReadingDate))
				TUI_Info(fmt.Sprintf("%s", jsonRecord))
			}
		}
	}
}
