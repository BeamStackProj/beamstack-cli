package open

import (
	"os"
	"sync"

	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

var (
	config     *rest.Config = utils.GetKubeConfig()
	LocalPort  uint16
	TargetPort uint16
	wg         sync.WaitGroup
	readyCh    chan struct{}               = make(chan struct{})
	stopCh     chan struct{}               = make(chan struct{}, 1)
	stream     genericclioptions.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
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
