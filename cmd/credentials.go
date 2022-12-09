/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var credentialsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientID, clientSecret, err := paramsPrompt()
		viper.Set(lib.EVERACTIVE_CLIENT_ID, clientID)
		viper.Set(lib.EVERACTIVE_CLIENT_SECRET, clientSecret)
		err = viper.WriteConfigAs(lib.ConfigurationFile)
		if err != nil {
			fmt.Println(fmt.Sprintf("error writing configuration file %s", err.Error()))
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf("Saved configuration in %s", lib.ConfigurationFile))
	},
}

var credentialsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(lib.ConfigurationFile)
		if err != nil {
			fmt.Println(fmt.Sprintf("configuration file could not be found at %s", lib.ConfigurationPath))
		}
		fmt.Println(fmt.Sprintf("configuration file is located at %s", lib.ConfigurationFile))
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
