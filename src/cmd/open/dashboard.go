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

var (
	DashboardLocalPort  uint16 = 8081
	DashboardTargetPort uint16 = 8081
)

// infoCmd represents the info command
var DashboardCmd = &cobra.Command{
	Use:   "show",
	Short: "opens up grafana dashboard",
	Long:  `opens up grafana dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := utils.ValidateCluster()
		if err != nil {
			fmt.Println(err)
			return
		}
		if profile.Monitoring == nil {
			fmt.Println("Monitoring not enabled on this cluster")
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
				Namespace: "monitoring",
			},
		}

		GrafanaSvc, err := clientset.CoreV1().Services("monitoring").Get(context.TODO(), svc.Name, metav1.GetOptions{})

		if errors.IsNotFound(err) {
			fmt.Printf("could not find service/monitoring/%s\n", svc.Name)
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
						PodPort:   DashboardTargetPort,
						LocalPort: DashboardLocalPort,
						Streams:   stream,
						StopCh:    stopCh,
						ReadyCh:   readyCh,
					},
					Service: *GrafanaSvc,
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
	FlinkClusterCmd.Flags().Uint16Var(&DashboardTargetPort, "targetport", DashboardTargetPort, "target container port")
	FlinkClusterCmd.Flags().Uint16Var(&DashboardLocalPort, "localport", DashboardLocalPort, "This port will be forwarded to the target port on the cluster")
}
