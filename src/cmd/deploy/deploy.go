/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deploy

import (
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy resources to k8s cluster",
	Long:  `deploy resources to k8s cluster`,
}

func init() {
	DeployCmd.AddCommand(PipelineCmd)
}
