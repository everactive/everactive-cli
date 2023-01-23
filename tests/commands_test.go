package tests

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/everactive/everactive-cli/cmd"
	"gitlab.com/everactive/everactive-cli/services"
	"os"
	"strings"
	"testing"
	"time"
)

func SetupCommandTests() {
	client := services.GetApiClient(false)
	httpmock.ActivateNonDefault(client.GetClient())
	httpmock.Activate()
	cmd.ApiClient = services.NewEveractiveAPIService(client, false)
}

func TestCmdHeartbeat(t *testing.T) {
	SetupCommandTests()
	defer httpmock.DeactivateAndReset()

	endpoint := regexUrlPath("/auth/token")
	httpmock.RegisterResponder("POST", endpoint,
		httpmock.NewStringResponder(201, RESPONSE_TOKEN))

	endpoint = regexUrlPath("/ds/v1/health")
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_HEALTH_HAPPY))

	//capture stdout for assertions.
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "heartbeat"}
	cmd.Execute()
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	cleanOutput := string(buf[:n-1])
	os.Stdout = rescueStdout
	w.Close()
	r.Close()
	assert.Equal(t, cmd.MSG_HEARTBEAT_SUCCESS, cleanOutput)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])

}

func TestCmdListSensors(t *testing.T) {
	SetupCommandTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/auth/token")
	httpmock.RegisterResponder("POST", endpoint,
		httpmock.NewStringResponder(201, RESPONSE_TOKEN))

	endpoint = regexUrlPath("/ds/v1/eversensors")
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_LIST_SENSORS_HAPPY))

	//capture stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "list-sensors"}
	cmd.Execute()

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	cleanOutput := string(buf[:n-1])
	os.Stdout = rescueStdout
	w.Close()
	r.Close()

	assert.True(t, strings.Contains(cleanOutput, TEST_MAC_ADDRESS_1))
	assert.True(t, strings.Contains(cleanOutput, TEST_MAC_ADDRESS_2))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
}

func TestCmdDataLast(t *testing.T) {
	SetupCommandTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/auth/token")
	httpmock.RegisterResponder("POST", endpoint,
		httpmock.NewStringResponder(201, RESPONSE_TOKEN))

	endpoint = regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings/last", TEST_MAC_ADDRESS_1))
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_DATA_LAST_MAC_1_HAPPY))

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "data", "--range", "last", "--sensor", TEST_MAC_ADDRESS_1}
	cmd.Execute()

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	cleanOutput := string(buf[:n-1])
	os.Stdout = rescueStdout
	w.Close()
	r.Close()
	assert.True(t, strings.Contains(cleanOutput, TEST_MAC_ADDRESS_1))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
}

func TestCmdData1hr(t *testing.T) {
	SetupCommandTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/auth/token")
	httpmock.RegisterResponder("POST", endpoint,
		httpmock.NewStringResponder(201, RESPONSE_TOKEN))

	endpoint = regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings\\?start-time=\\d{10}&end-time=\\d{10}", TEST_MAC_ADDRESS_2))
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_DATA_RANGE_MAC_2_HAPPY))

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "data", "--range", "1h", "--sensor", TEST_MAC_ADDRESS_2}
	cmd.Execute()

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	cleanOutput := string(buf[:n-1])
	os.Stdout = rescueStdout
	w.Close()
	r.Close()
	assert.True(t, strings.Contains(cleanOutput, TEST_MAC_ADDRESS_2))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
}

func TestCmdDataStartEnd(t *testing.T) {
	SetupCommandTests()
	defer httpmock.DeactivateAndReset()
	endpoint := regexUrlPath("/auth/token")
	httpmock.RegisterResponder("POST", endpoint,
		httpmock.NewStringResponder(201, RESPONSE_TOKEN))

	endpoint = regexUrlPath(fmt.Sprintf("/ds/v1/eversensors/%s/readings\\?start-time=1670507054&end-time=1670533006", TEST_MAC_ADDRESS_2))
	httpmock.RegisterResponder("GET", endpoint,
		httpmock.NewStringResponder(200, RESPONSE_DATA_RANGE_MAC_2_HAPPY))

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "data", "--range", "1670507054-1670533006", "--sensor", TEST_MAC_ADDRESS_2}
	cmd.Execute()

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	cleanOutput := string(buf[:n-1])
	os.Stdout = rescueStdout
	w.Close()
	r.Close()
	assert.True(t, strings.Contains(cleanOutput, TEST_MAC_ADDRESS_2))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info[fmt.Sprintf("GET %s", endpoint)])
}

func TestCalculateRageStartEnd(t *testing.T) {
	start, end := cmd.CalculateRage("1670507054-1670533006")
	assert.Equal(t, int64(1670507054), start)
	assert.Equal(t, int64(1670533006), end)
}

func TestCalculateRageDuration(t *testing.T) {
	now := time.Now().UTC().Unix()
	start, end := cmd.CalculateRage("1h")
	assert.True(t, start > 0)
	assert.True(t, end > 0)
	assert.Equal(t, int64(3600), end-start)
	assert.True(t, start < end)
	assert.True(t, start < now)
	assert.True(t, end >= now)

	start, end = cmd.CalculateRage("60m")
	assert.True(t, start > 0)
	assert.True(t, end > 0)
	assert.Equal(t, int64(3600), end-start)
	assert.True(t, start < end)
	assert.True(t, start < now)
	assert.True(t, end >= now)

	start, end = cmd.CalculateRage("3600s")
	assert.True(t, start > 0)
	assert.True(t, end > 0)
	assert.Equal(t, int64(3600), end-start)
	assert.True(t, start < end)
	assert.True(t, start < now)
	assert.True(t, end >= now)

	//invalid range
	start, end = cmd.CalculateRage("36h")
	assert.Equal(t, int64(0), start)
	assert.Equal(t, int64(0), end)

}


