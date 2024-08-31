/*
Copyright © 2024 MavenCode <opensource-dev@mavencode.com>
*/
package info

import (
	"fmt"

	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Long Description and Example initizalization
var (
	infoLongDesc = utils.LongDesc(`
		This command displays detailed information about specified resources such as clusters.
		`)
)

// infoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "get cluster info",
	Long:  infoLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := utils.GetCurrentContext()
		if err != nil {
			fmt.Println("could not retrieve current context")
		}
		fmt.Printf("current kube context: %s\n", ctx)
		contextsStringMap := viper.GetStringMapString("contexts")

		var profileName string
		if p, ok := contextsStringMap[ctx]; ok {
			fmt.Printf("Current profile: %s\n", p)
			profileName = p
		}

		if profileName == "" {
			return
		}
		profile, err := utils.GetProfile(profileName)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("installed packages:")
		for _, p := range profile.Packages {
			fmt.Println("\t", p.Name)
			for _, d := range p.Dependencies {
				fmt.Println("\t\t", d.Name)
			}
		}
	},
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
