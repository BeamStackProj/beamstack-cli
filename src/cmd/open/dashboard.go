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

// infoCmd represents the info command
var DashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "opens up grafana dashboard",
	Long:  `opens up grafana dashboard`,
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
				Name:      "grafana",
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
						PodPort:   TargetPort,
						LocalPort: LocalPort,
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

		println("insert 'q' to exit")

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
	DashboardCmd.Flags().Uint16("targetport", 3000, "target container port")
	DashboardCmd.Flags().Uint16("localport", 3000, "This port will be forwarded to the target port on the cluster")

}
