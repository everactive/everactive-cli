package lib

import (
	"fmt"
	"github.com/spf13/viper"
)

var Version string

const (
	EVERACTIVE_API_URL = "EVERACTIVE_API_URL"
	EVERACTIVE_CLIENT_ID = "EVERACTIVE_CLIENT_ID"
	EVERACTIVE_CLIENT_SECRET = "EVERACTIVE_CLIENT_SECRET"
)

var DefaultConfigParams = map[string]interface{}{
	EVERACTIVE_API_URL: "https://api.data.everactive.com",
}

func InitConfiguration() {
	Version="set_at_build"
	//bind and setup defaults
	for key, value := range DefaultConfigParams {
		viper.SetDefault(key, value)
	}
	viper.AutomaticEnv()
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.everactive") // call multiple times to add many search paths
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	_ = viper.ReadInConfig()                 // Find and read the config file
}

func SaveConfiguration() {
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("error writing configuration file %w", err))
	}
}

