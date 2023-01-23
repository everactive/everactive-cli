package tests

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/everactive/everactive-cli/lib"
	"gitlab.com/everactive/everactive-cli/services"
	"net/http"
	"strings"
	"testing"
)

var apiTestClient services.EveractiveAPIService

func SetupTests() {
	fmt.Println("Preparing Mock for API")
	apiTestClient = services.NewEveractiveAPIService(services.GetApiClient(false), true)
	httpmock.ActivateNonDefault(apiTestClient.Client.GetClient())
	httpmock.Activate()
}

func regexUrlPath(path string) string {
	return fmt.Sprintf("=~.*%s", strings.ReplaceAll(path, "/", "\\/"))
}

func TestHealth(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()

	endpoint := regexUrlPath("/ds/v1/health")
	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, RESPONSE_HEALTH_HAPPY)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
	result := apiTestClient.Health()
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.Equal(t, true, result)
}

func TestHealthBad(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()

	endpoint := regexUrlPath("/ds/v1/health")
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(400, ""))
	result := apiTestClient.Health()
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.Equal(t, false, result)
}

func TestGetSensorListHappy(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/ds/v1/eversensors")
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_LIST_SENSORS_HAPPY))

	result, err := apiTestClient.GetSensorsList()
	assert.Nil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])

	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, TEST_CUSTOMER_ID, result.Data[0].Customer.ID)
	assert.Equal(t, TEST_CUSTOMER_NAME, result.Data[0].Customer.Name)
	assert.Equal(t, TEST_GATEWAY, result.Data[0].LastAssociation.GatewaySerialNumber)
	sensors := make([]string, 0)
	for _, sensor := range result.Data {
		sensors = append(sensors, sensor.MacAddress)
	}
	assert.Contains(t, sensors, TEST_MAC_ADDRESS_1)
	assert.Contains(t, sensors, TEST_MAC_ADDRESS_2)
}

func TestGetSensorListBad(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/ds/v1/eversensors")
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_LIST_SENSORS_BAD))

	result, err := apiTestClient.GetSensorsList()
	assert.NotNil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.Nil(t, result)
}

func TestGetSensorListError(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/ds/v1/eversensors")
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(502, ""))

	result, err := apiTestClient.GetSensorsList()
	assert.NotNil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.Nil(t, result)
}

func TestGetSensorLastReadingHappy(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings/last", TEST_MAC_ADDRESS_1))

	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_DATA_LAST_MAC_1_HAPPY))

	result, err := apiTestClient.GetSensorLastReading(TEST_MAC_ADDRESS_1)
	assert.Nil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.NotNil(t, result)
	assert.NotNil(t, result.Data)
	byteData, _ := json.Marshal(result.Data)
	sensorReading := lib.SensorReading{}
	err = json.Unmarshal(byteData, &sensorReading)
	assert.Nil(t, err)
	assert.Equal(t, TEST_MAC_ADDRESS_1, sensorReading.MacAddress)
	assert.Equal(t, TEST_GATEWAY, sensorReading.GatewaySerialNumber)
}

func TestGetSensorLastReadingError(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings/last", TEST_MAC_ADDRESS_1))

	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(502, ""))

	result, err := apiTestClient.GetSensorLastReading(TEST_MAC_ADDRESS_1)
	assert.NotNil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.Nil(t, result)
}

func TestGetSensorDataRangeHappy(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	start := int64(1674061998)
	end := int64(1674148398)
	endpoint := regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings\\?start-time=%d&end-time=%d", TEST_MAC_ADDRESS_2, start, end))
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_DATA_RANGE_MAC_2_HAPPY))

	result, err := apiTestClient.GetSensorReadings(TEST_MAC_ADDRESS_2, start, end)
	assert.Nil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.NotNil(t, result)
	assert.NotNil(t, result.Data)
	byteData, _ := json.Marshal(result.Data)
	sensorReadings := make([]lib.SensorReading, 0)
	err = json.Unmarshal(byteData, &sensorReadings)
	assert.Nil(t, err)
	assert.Equal(t, RESPONSE_DATA_RANGE_MAC_2_HAPPY_LEN, len(sensorReadings))
	assert.Equal(t, TEST_MAC_ADDRESS_2, sensorReadings[0].MacAddress)
	assert.Equal(t, TEST_GATEWAY, sensorReadings[0].GatewaySerialNumber)
	assert.Equal(t, TEST_MAC_ADDRESS_2, sensorReadings[RESPONSE_DATA_RANGE_MAC_2_HAPPY_LEN-1].MacAddress)
	assert.Equal(t, TEST_GATEWAY, sensorReadings[RESPONSE_DATA_RANGE_MAC_2_HAPPY_LEN-1].GatewaySerialNumber)
}

func TestGetSensorDataRangeError(t *testing.T) {
	SetupTests()
	defer httpmock.DeactivateAndReset()
	start := int64(1674061998)
	end := int64(1674148398)
	endpoint := regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings\\?start-time=%d&end-time=%d", TEST_MAC_ADDRESS_2, start, end))
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(502, ""))

	result, err := apiTestClient.GetSensorReadings(TEST_MAC_ADDRESS_2, start, end)
	assert.NotNil(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
	assert.Nil(t, result)
}
