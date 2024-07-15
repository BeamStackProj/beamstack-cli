/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/BeamStackProj/beamstack-cli/src/cmd/deploy"
	"github.com/BeamStackProj/beamstack-cli/src/cmd/info"
	"github.com/BeamStackProj/beamstack-cli/src/cmd/initialize"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "beamstack",
	Short: "Welcome to the beamstack cli tool",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubCommandPallets() {
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(deploy.DeployCmd)
	rootCmd.AddCommand(initialize.InitCmd)
}

func init() {
	cobra.OnInitialize(utils.InitConfig)
	addSubCommandPallets()

}
