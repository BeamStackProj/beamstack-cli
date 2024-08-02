/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "'create' sub commands",
	Long:  `create flink or spark clusters`,
}

func init() {

	CreateCmd.AddCommand(FlinkClusterCmd)
}
