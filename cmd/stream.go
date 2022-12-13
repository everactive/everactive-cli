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
var ringBell = false
const (
	loopPeriodMinutes = 10
	loopDelaySeconds  = 5
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: fmt.Sprintf("Retrieve the eversensor readings in streaming mode, starting now and polling for new data every %d seconds", loopDelaySeconds),
	Long: `The sensor MAC address is required. To stop the stream type ctrl+c`,
	Run: func(cmd *cobra.Command, args []string) {
		sensorFlag := cmd.Flag("sensor")
		if sensorFlag != nil {
			done := make(chan os.Signal, 1)
			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
			Tui_debug("press ctrl+c to stop...")
			go streamLoop(sensorFlag.Value.String())
			<-done // Will block here until user hits ctrl+c
			Tui_debug("good bye")
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.Flags().String("sensor", "", "Required. Mac address of the sensor")
	_ = streamCmd.MarkFlagRequired("sensor")
	streamCmd.Flags().BoolVarP(&ringBell, "bell", "b", false, "Ring the bell on every reading")

}

func streamLoop(sensorFilter string) {
	api := services.NewEveractiveAPIService(false, context.Background())
	for {
		endTime := time.Now()
		start := endTime.Add(time.Minute * -loopPeriodMinutes).Unix()
		end := endTime.Unix()
		Tui_debug(fmt.Sprintf("readings from start %d end %d", start, end))
		time.Sleep(time.Second * loopDelaySeconds)
		response, err := api.GetSensorReadings(sensorFilter, start, end)
		if err != nil {
			Tui_error(fmt.Sprintf("failed to retrieved sensors data: %s", err.Error()))
			os.Exit(1)
		}
		var reading lib.SensorReading
		for _, record := range response.Data {
			jsonRecord, err := json.Marshal(record)
			if err != nil {
				Tui_error(fmt.Sprintf("error processing response %s", err.Error()))
				os.Exit(1)
			}
			err = json.Unmarshal(jsonRecord, &reading)
			if err != nil {
				Tui_error(fmt.Sprintf("error processing response %s", err.Error()))
				os.Exit(1)
			}
			if reading.Timestamp > lastTimestamp {
				lastTimestamp = reading.Timestamp
				Tui_debug(fmt.Sprintf("pkt %d - timestamp: %s", reading.PacketNumberGateway, reading.ReadingDate))
				Tui_info(fmt.Sprintf("%s", jsonRecord))
				if ringBell {
					Tui_bell()
				}
			}
		}
	}
}
