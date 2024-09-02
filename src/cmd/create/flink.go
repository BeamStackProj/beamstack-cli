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
	cpu         string = "500m"
	memory      string = "1024Mi"
	cpuLimit    string = "1"
	memoryLimit string = "2048Mi"
	volumeSize  string = "1Gi"
	taskslots   uint8  = 1
	replicas    uint8  = 1
	Previledged bool   = false
)

// Description and Examples for creating flink clsuters
var (
	flinkLongDesc = utils.LongDesc(`
		Create a flink cluster with specified requirments.
		`)
)

// infoCmd represents the info command
var FlinkClusterCmd = &cobra.Command{
	Use:   "flink [NAME]",
	Short: "create a flink cluster",
	Long:  flinkLongDesc,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("flink command requires exactly one argument: cluster Name. Provided %d arguments", len(args))
		}
		return nil
	},
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

		var spec types.FlinkDeploymentSpec
		taskmanagerImage := fmt.Sprintf("beamstackproj/beam-harness-%s:latest", flinkVersion)

		if Previledged {
			flinkImage := fmt.Sprintf("beamstackproj/flink-%s-docker:latest", flinkVersion)

			spec = types.FlinkDeploymentSpec{
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
								Name:         "flink-main-container",
								Image:        flinkImage,
								VolumeMounts: []v1.VolumeMount{},
								SecurityContext: &v1.SecurityContext{
									Privileged: func(b bool) *bool { return &b }(true),
								},
							},
						},
						Volumes: []v1.Volume{},
					},
				},
				JobManager: types.JobManagerSpec{
					Replicas: 1,
					Resource: types.Resource{
						Memory:      memory,
						CPU:         cpu,
						CPULimit:    cpuLimit,
						MemoryLimit: memoryLimit,
					},
				},
				TaskManager: types.TaskManagerSpec{
					Replicas: replicas,
					Resource: types.Resource{
						Memory:      memory,
						CPU:         cpu,
						CPULimit:    cpuLimit,
						MemoryLimit: memoryLimit,
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
								{
									Name: "flink-main-container",
									SecurityContext: &v1.SecurityContext{
										Privileged: func(b bool) *bool { return &b }(true),
									},
									VolumeMounts: []v1.VolumeMount{
										{
											MountPath: "/var/run/docker.sock",
											Name:      "docker-socket",
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
								{
									Name: "docker-socket",
									VolumeSource: v1.VolumeSource{
										HostPath: &v1.HostPathVolumeSource{
											Path: "/var/run/docker.sock",
											Type: func() *v1.HostPathType {
												t := v1.HostPathSocket
												return &t
											}(),
										},
									},
								},
							},
						},
					},
				},
			}
		} else {
			flinkImage := fmt.Sprintf("beamstackproj/flink-%s:latest", flinkVersion)
			spec = types.FlinkDeploymentSpec{
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
						Memory:      memory,
						CPU:         cpu,
						CPULimit:    cpuLimit,
						MemoryLimit: memoryLimit,
					},
				},
				TaskManager: types.TaskManagerSpec{
					Replicas: replicas,
					Resource: types.Resource{
						Memory:      memory,
						CPU:         cpu,
						CPULimit:    cpuLimit,
						MemoryLimit: memoryLimit,
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
	FlinkClusterCmd.Flags().BoolVar(&Previledged, "Previledged", Previledged, "")
}
