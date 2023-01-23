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
	"log"
)

type EveractiveAPI interface {
	Health() bool
	GetSensorsList() (*lib.GetEversensorsShortResponse, error)
	GetSensorReadings(mac string, start, end int64) (*lib.GetSensorReadingsRawResponse, error)
	GetSensorLastReading(mac string) (*lib.GetSensorLastReadingRawResponse, error)
}

type EveractiveAPIService struct {
	Debug  bool
	Client *resty.Client
}

func NewEveractiveAPIService(client *resty.Client, debug bool) EveractiveAPIService {
	service := EveractiveAPIService{
		Debug: debug,
	}
	if client == nil {
		client = GetApiClient(debug)
	}
	service.Client = client
	return service
}

func (api EveractiveAPIService) Health() bool {
	endpoint := fmt.Sprintf("%s/ds/v1/health", viper.GetString(lib.EVERACTIVE_API_URL))
	if api.Debug {
		log.Println(fmt.Sprintf("calling api Health at %s", endpoint))
	}
	resp, err := api.Client.R().Get(endpoint)
	if err != nil {
		return false
	}
	return resp.IsSuccess()
}

func (api EveractiveAPIService) GetSensorsList() (*lib.GetEversensorsShortResponse, error) {
	endpoint := fmt.Sprintf("%s/ds/v1/eversensors", viper.GetString(lib.EVERACTIVE_API_URL))
	if api.Debug {
		log.Println(fmt.Sprintf("calling api GetSensorsList at %s", endpoint))
	}
	resp, err := api.Client.R().Get(endpoint)
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
	if api.Debug {
		log.Println(fmt.Sprintf("calling api GetSensorReadings at %s", endpoint))
	}
	resp, err := api.Client.R().Get(endpoint)
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
	if api.Debug {
		log.Println(fmt.Sprintf("calling api GetSensorLastReading at %s", endpoint))
	}
	resp, err := api.Client.R().Get(endpoint)
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
