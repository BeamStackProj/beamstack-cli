/*
Copyright © 2024 MavenCode <opensource-dev@mavencode.com>
*/
package cmd

import (
	"os"

	"github.com/BeamStackProj/beamstack-cli/src/cmd/create"
	"github.com/BeamStackProj/beamstack-cli/src/cmd/deploy"
	"github.com/BeamStackProj/beamstack-cli/src/cmd/info"
	"github.com/BeamStackProj/beamstack-cli/src/cmd/initialize"
	"github.com/BeamStackProj/beamstack-cli/src/cmd/open"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "beamstack",
	Short: "Welcome to the beamstack cli tool",
	Long:  ``,
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
	rootCmd.AddCommand(initialize.InitCmd)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(deploy.DeployCmd)
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(open.OpenCmd)
}

func init() {
	cobra.OnInitialize(utils.InitConfig)
	addSubCommandPallets()

}
