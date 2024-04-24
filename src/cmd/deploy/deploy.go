/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Pallete that contains information based commands",
	Long:  `Pallete that contains information based commands`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploying resources!!")
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
