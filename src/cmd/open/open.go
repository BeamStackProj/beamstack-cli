package open

import (
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "opens up resource ui",
	Long:  `opens up resource ui`,
	Args:  cobra.ExactArgs(1),
}

func init() {
	OpenCmd.AddCommand(FlinkClusterCmd)
	OpenCmd.AddCommand(DashboardCmd)
}
