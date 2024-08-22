package info

import (
	"fmt"
	"time"

	info_handler "github.com/BeamStackProj/beamstack-cli/src/handlers/info"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

var (
	clusterLongDesc = utils.LongDesc(`
		This command fetches detailed information about the Kubernetes cluster,
  	including the status of each node and the overall health of the cluster.
		`)

	clusterExample = utils.Examples(`
  # Get cluster information
  beamstack info cluster
	`)
)

// infoCmd represents the info command
var ClusterCmd = &cobra.Command{
	Use:     "cluster",
	Short:   "get information the kubernetes cluster",
	Long:    clusterLongDesc,
	Example: clusterExample,
	Run: func(cmd *cobra.Command, args []string) {
		clusterHealth, err := info_handler.Health()

		if err != nil {
			fmt.Printf("error connecting to kubernetes cluser: %v\n", err)
			return
		}

		fmt.Printf("%-50s %-10s %s\n", "NAME", "STATUS", "AGE")
		for _, node := range clusterHealth.Nodes {
			status := "Not Ready"

			for _, condition := range node.Status.Conditions {
				if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
					status = "Ready"
					break
				}
			}
			age := time.Since(node.CreationTimestamp.Time).Round(time.Hour)

			var ageString string
			if age.Hours() > 50 {
				ageString = fmt.Sprintf("%.0fd", age.Hours()/24)
			} else {
				ageString = fmt.Sprintf("%.0fh", age.Hours())
			}

			fmt.Printf("%-50s %-10s %s\n", node.Name, status, ageString)
		}

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
