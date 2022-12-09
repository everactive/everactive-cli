package lib

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	EVERACTIVE_API_URL = "EVERACTIVE_API_URL"
	EVERACTIVE_CLIENT_ID = "EVERACTIVE_CLIENT_ID"
	EVERACTIVE_CLIENT_SECRET = "EVERACTIVE_CLIENT_SECRET"
	EVERACTIVE_CONFIGURATION_FILE = "config.yaml"
)
var Version string
var HomeDirectory string
var ConfigurationPath string
var ConfigurationFile string
var DefaultConfigParams = map[string]interface{}{
	EVERACTIVE_API_URL: "https://api.data.everactive.com",
}

func InitConfiguration() {
	HomeDirectory, _ = os.UserHomeDir()
	ConfigurationPath = fmt.Sprintf("%s/.everactive", HomeDirectory)
	ConfigurationFile= fmt.Sprintf("%s/%s", ConfigurationPath, EVERACTIVE_CONFIGURATION_FILE)
	_=os.Mkdir(ConfigurationPath, 0755)
	//bind and setup defaults
	for key, value := range DefaultConfigParams {
		viper.SetDefault(key, value)
	}
	viper.AutomaticEnv()
	viper.SetConfigFile(ConfigurationFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigurationPath)
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
}

