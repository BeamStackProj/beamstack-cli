/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package info

import (
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Pallete that contains information based commands",
	Long:  `Pallete that contains information based commands`,
}

func init() {

	InfoCmd.AddCommand(ClusterCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
