/*
Copyright © 2024 MavenCode <opensource-dev@mavencode.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get current beamstack version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.Get("version"))
	},
}
