package create

import (
	"fmt"

	"github.com/BeamStackProj/beamstack-cli/src/objects"
	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	cpu         string = "2"
	memory      string = "2048Mi"
	cpuLimit    string = "2"
	memoryLimit string = "2048Mi"
	volumeSize  string = "1Gi"
	taskslots   uint8  = 10
	replicas    uint8  = 1
)

// Description and Examples for creating flink clsuters
var (
	flinkLongDesc = utils.LongDesc(`
		Create a flink cluster with specified requirments.
		`)
)

// infoCmd represents the info command
var FlinkClusterCmd = &cobra.Command{
	Use:   "flink-cluster",
	Short: "create a flink cluster",
	Long:  flinkLongDesc,
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

		namespace := "flink"
		flinkVersion := "v1_16"
		flinkImage := fmt.Sprintf("beamstackproj/flink-%s:latest", flinkVersion)
		// flinkImage := "localhost:5000/flink-v1_16"
		taskmanagerImage := fmt.Sprintf("beamstackproj/flink-harness-%s:latest", flinkVersion)
		// taskmanagerImage := "localhost:5000/harness"
		ClaimName := fmt.Sprintf("%s-pvc", args[0])
		fmt.Printf("creating flink cluster %s\n", args[0])

		config := utils.GetKubeConfig()

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := objects.CreatePVC(clientset, ClaimName, namespace, volumeSize); err != nil {
			fmt.Println(err)
			return
		}

		spec := types.FlinkDeploymentSpec{
			Image:           &flinkImage,
			ImagePullPolicy: "IfNotPresent",
			FlinkVersion:    flinkVersion,
			FlinkConfiguration: map[string]string{
				"taskmanager.numberOfTaskSlots": fmt.Sprintf("%d", taskslots),
			},
			ServiceAccount: "flink",
			PodTemplate: &v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "flink-main-container",
							VolumeMounts: []v1.VolumeMount{
								{
									MountPath: "/opt/flink/log",
									Name:      "flink-logs",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "flink-logs",
						},
					},
				},
			},
			JobManager: types.JobManagerSpec{
				Replicas: 1,
				Resource: types.Resource{
					Memory: memory,
					CPU:    cpu,
				},
			},
			TaskManager: types.TaskManagerSpec{
				Replicas: replicas,
				Resource: types.Resource{
					Memory: memory,
					CPU:    cpu,
				},

				PodTemplate: &v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "worker",
								Image: taskmanagerImage,
								Args:  []string{"-worker_pool"},
								Ports: []v1.ContainerPort{
									{
										Name:          "harness-port",
										ContainerPort: 50000,
									},
								},
								VolumeMounts: []v1.VolumeMount{
									{
										MountPath: "/pvc",
										Name:      "flink-cluster-pvc",
									},
								},
							},
						},
						Volumes: []v1.Volume{
							{
								Name: "flink-cluster-pvc",
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: ClaimName,
									},
								},
							},
						},
					},
				},
			},
		}

		err = objects.CreateDynamicResource(
			metav1.TypeMeta{
				APIVersion: "flink.apache.org/v1beta1",
				Kind:       "FlinkDeployment",
			},
			metav1.ObjectMeta{
				Name:      args[0],
				Namespace: "flink",
			},
			spec,
			"flinkdeployments",
		)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Flink cluster %s created\n", args[0])
	},
}

func init() {
	FlinkClusterCmd.Flags().StringVar(&cpu, "cpu", cpu, "Cpu request for task manager")
	FlinkClusterCmd.Flags().StringVar(&cpuLimit, "cpuLimit", cpu, "Cpu request for task manager")
	FlinkClusterCmd.Flags().StringVar(&memory, "memory", memory, "Cpu request for task manager")
	FlinkClusterCmd.Flags().StringVar(&memoryLimit, "memoryLimit", memoryLimit, "Cpu request for task manager")
	FlinkClusterCmd.Flags().Uint8Var(&replicas, "replicas", replicas, "numbers of replicas sets for task manager")
	FlinkClusterCmd.Flags().Uint8Var(&taskslots, "taskslots", taskslots, "numbers of taskslots to be created for the task manager")
	FlinkClusterCmd.Flags().StringVar(&volumeSize, "volumeSize", volumeSize, "size of persistent volume to be attached to flink cluster")
}
