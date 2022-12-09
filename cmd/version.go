/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/everactive/everactive-cli/lib"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of the application",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lib.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
