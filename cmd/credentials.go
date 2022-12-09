package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/everactive/everactive-cli/lib"
	"os"
)

// credentialsCmd represents the credentials command
var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Credentials configuration",
	Long: `Use this command to create or find the configuration file. 
The configuration can also be set via environment variables.`,
}

var credentialsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the configuration file with the API credentials",
	Long: `A file is created in the default configuration path with the Client ID
and the Client Secret. This command needs to be run only once.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientID, clientSecret, err := paramsPrompt()
		viper.Set(lib.EVERACTIVE_CLIENT_ID, clientID)
		viper.Set(lib.EVERACTIVE_CLIENT_SECRET, clientSecret)
		err = viper.WriteConfigAs(lib.ConfigurationFile)
		if err != nil {
			TUI_Error(fmt.Sprintf("error writing configuration file %s", err.Error()))
			os.Exit(1)
		}
		TUI_Info(fmt.Sprintf("saved configuration in %s", lib.ConfigurationFile))
	},
}

var credentialsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(lib.ConfigurationFile)
		if err != nil {
			TUI_Error(fmt.Sprintf("configuration file could not be found at %s", lib.ConfigurationPath))
		}
		TUI_Info(fmt.Sprintf("configuration file is located at %s", lib.ConfigurationFile))
	},
}

func init() {
	rootCmd.AddCommand(credentialsCmd)
	credentialsCmd.AddCommand(credentialsInitCmd)
	credentialsCmd.AddCommand(credentialsFindCmd)
}

func paramsPrompt() (id, key string, err error) {
	prompt := promptui.Prompt{
		Label: "Please enter the Client ID",
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("client-id not valid")
			}
			return nil
		},
	}
	id, err = prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("client prompt fail: %w", err)
	}

	prompt = promptui.Prompt{
		Mask:  '*',
		Label: "Please enter the Client Secret",
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("client secret not valid")
			}
			return nil
		},
	}
	key, err = prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("client secret prompt fail: %w", err)
	}

	return id, key, nil
}
