package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"gitlab.com/everactive/everactive-cli/lib"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type EveractiveAPI interface {
	Health() bool
	GetSensorsList() (*lib.GetEversensorsShortResponse, error)
	GetReadings(sensorMac string)
	GetSensorLastReading(mac string)
}

type EveractiveAPIService struct {
	debug  bool
	client *resty.Client
}

func NewEveractiveAPIService(debug bool) EveractiveAPIService {
	service := EveractiveAPIService{
		debug:  debug,
		client: GetApiClient(debug)}
	return service
}

func (api EveractiveAPIService) Health() bool {
	endpoint := fmt.Sprintf("%s/ds/v1/health", viper.GetString(lib.EVERACTIVE_API_URL))
	resp, err := api.client.R().Get(endpoint)
	if err != nil {
		return false
	}
	return resp.IsSuccess()
}

func (api EveractiveAPIService) GetSensorsList() (*lib.GetEversensorsShortResponse, error) {
	endpoint := fmt.Sprintf("%s/ds/v1/eversensors", viper.GetString(lib.EVERACTIVE_API_URL))
	resp, err := api.client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	var response lib.GetEversensorsShortResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, errors.New("error parsing API response: " + err.Error())
	}
	//TODO: Pagination
	return &response, nil
}

func (api EveractiveAPIService) GetSensorReadings(mac string, start, end int64) (*lib.GetSensorReadingsRawResponse, error) {
	endpoint := fmt.Sprintf("%s/ds/v1/eversensors/%s/readings?start-time=%d&end-time=%d",
		viper.GetString(lib.EVERACTIVE_API_URL), mac, start, end)
	resp, err := api.client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	var response lib.GetSensorReadingsRawResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, errors.New("error parsing API response: " + err.Error())
	}
	return &response, nil
}

func (api EveractiveAPIService) GetSensorLastReading(mac string) (*lib.GetSensorLastReadingRawResponse, error) {
	endpoint := fmt.Sprintf("%s/ds/v1/eversensors/%s/readings/last",
		viper.GetString(lib.EVERACTIVE_API_URL), mac)
	resp, err := api.client.R().Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	var response lib.GetSensorLastReadingRawResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, errors.New("error parsing API response: " + err.Error())
	}
	return &response, nil
}

func GetApiClient(debug bool) *resty.Client {
	oauth2Config := &clientcredentials.Config{
		ClientID:     viper.GetString(lib.EVERACTIVE_CLIENT_ID),
		ClientSecret: viper.GetString(lib.EVERACTIVE_CLIENT_SECRET),
		TokenURL:     fmt.Sprintf("%s/auth/token", viper.GetString(lib.EVERACTIVE_API_URL)),
		Scopes:       []string{"user", "developer"},
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	oauth2Client := oauth2Config.Client(context.Background())
	apiClient := resty.NewWithClient(oauth2Client)
	apiClient.SetRetryCount(3)
	apiClient.SetDebug(debug)
	return apiClient
}
