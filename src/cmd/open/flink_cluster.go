package open

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Description and Examples for opening grafana dashboard
var (
	openFlinkClusterLongDesc = utils.LongDesc(`
		This command opens up the GUI for the flink clusters and forward it to a specified local port.
		`)
)

// infoCmd represents the info command
var FlinkClusterCmd = &cobra.Command{
	Use:   "flink [Name]",
	Short: "opens up flink cluster ui",
	Long:  openFlinkClusterLongDesc,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("command requires exactly one argument: Flink Cluster Name. Provided %d arguments", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		LocalPort, err := cmd.Flags().GetUint16("localport")
		if err != nil {
			fmt.Println(err)
			return
		}

		TargetPort, err := cmd.Flags().GetUint16("targetport")

		if err != nil {
			fmt.Println(err)
			return
		}

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
						PodPort:   TargetPort,
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
	FlinkClusterCmd.Flags().Uint16("targetport", 8081, "target container port")
	FlinkClusterCmd.Flags().Uint16("localport", 8081, "This port will be forwarded to the target port on the cluster")

}
