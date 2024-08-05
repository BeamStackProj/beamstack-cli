package open

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	config    *rest.Config = utils.GetKubeConfig()
	LocalPort uint16       = 8081
	TagetPort uint16       = 8081
	wg        sync.WaitGroup
	readyCh   chan struct{}               = make(chan struct{})
	stopCh    chan struct{}               = make(chan struct{}, 1)
	stream    genericclioptions.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
)

// infoCmd represents the info command
var FlinkClusterCmd = &cobra.Command{
	Use:   "show",
	Short: "opens up flink cluster ui",
	Long:  `opens up flink cluster ui`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := utils.ValidateCluster()
		if err != nil {
			fmt.Println(err)
			return
		}
		if profile.Operators.Flink == nil {
			fmt.Println("Flink Operator not initialized on this cluster")
			return
		}

		clientset, err := kubernetes.NewForConfig(config)

		if err != nil {
			fmt.Println(err)
			return
		}

		svc := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-rest", args[0]),
				Namespace: "flink",
			},
		}

		FlinkClusterSvc, err := clientset.CoreV1().Services("flink").Get(context.TODO(), svc.Name, metav1.GetOptions{})

		if errors.IsNotFound(err) {
			fmt.Printf("could not find service/flink/%s\n", svc.Name)
			return
		} else if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			close(stopCh)
			wg.Done()
		}()

		go func() {
			err := utils.PortForwardSvc(
				clientset,
				types.PortForwardASVCRequest{
					PortForward: types.PortForward{
						PodPort:   TagetPort,
						LocalPort: LocalPort,
						Streams:   stream,
						StopCh:    stopCh,
						ReadyCh:   readyCh,
					},
					Service: *FlinkClusterSvc,
				},
			)

			if err != nil {
				panic(err)
			}
		}()

		println("Port forwarding is ready to get traffic. insert 'q' to stop")

		go func() {
			var quitC string
			for {
				fmt.Scanln(&quitC)
				if quitC == "q" {
					break
				} else {
					fmt.Printf("Unknown command %s\n", quitC)
				}
			}
			wg.Done()
		}()

		wg.Wait()
	},
}

func init() {
	FlinkClusterCmd.Flags().Uint16Var(&TagetPort, "targetport", TagetPort, "target container port")
	FlinkClusterCmd.Flags().Uint16Var(&LocalPort, "localport", LocalPort, "This port will be forwarded to the target port on the cluster")
}
