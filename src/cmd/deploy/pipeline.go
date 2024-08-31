/*
Copyright © 2024 MavenCode <opensource-dev@mavencode.com>
*/
package deploy

import (
	"fmt"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/BeamStackProj/beamstack-cli/src/objects"
	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/BeamStackProj/beamstack-cli/src/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	deployLongDesc = utils.LongDesc(`
		Deploy an Apache Beam pipeline on a specified Operator.
		`)
)

var (
	flinkCluster          string       = ""
	PVCMountPath          string       = "/pvc"
	JobName               string       = "beamjob-asc"
	Parallelism           uint8        = 1
	Wait                  bool         = false
	Migrate               bool         = true
	config                *rest.Config = utils.GetKubeConfig()
	pipelineFilename      string
	CleanPipelineFilename string
)

type FileInfo struct {
	Src  string
	Dest string
}

// infoCmd represents the info command
var PipelineCmd = &cobra.Command{
	Use:   "pipeline [FILE]",
	Short: "Deploy an Apache Beam pipeline",
	Long:  deployLongDesc,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("pipeline command requires exactly one argument: the FILE to deploy. Provided %d arguments", len(args))
		}
		return nil
	},
	Run: DeployPipeline,
}

func init() {
	PipelineCmd.Flags().StringVar(&flinkCluster, "flink", flinkCluster, "Specify the Flink cluster to deploy the Apache Beam pipeline.")
	PipelineCmd.Flags().StringVar(&PVCMountPath, "pvcMountPath", PVCMountPath, "Mount path for the Persistent Volume Claim. Note: The mount path is set to 'pvc' during cluster creation, so changing this may cause issues.")
	PipelineCmd.Flags().StringVar(&JobName, "jobname", JobName, "Specify the name of the pipeline job.")
	PipelineCmd.Flags().Uint8Var(&Parallelism, "parallelism", Parallelism, "Set the pipeline parallelism.")
	PipelineCmd.Flags().BoolVarP(&Wait, "wait", "w", Wait, "Wait for the pipeline to complete.")
	PipelineCmd.Flags().BoolVarP(&Migrate, "migrate", "m", Migrate, "Migrate data to the Kubernetes cluster. This is necessary if the pipeline is to be run on local data. Pipeline Results will also be migrated to local system if wait is true.")

	PipelineCmd.MarkFlagRequired("flink")
}

func DeployPipeline(cmd *cobra.Command, args []string) {
	pipelineFilename = args[0]
	CleanPipelineFilename = args[0]

	profile, err := utils.ValidateCluster()

	if err != nil {
		fmt.Println(err)
		return
	}
	if profile.Operators.Flink == nil {
		fmt.Println("Flink Operator not initialized on this cluster")
		return
	}
	pipeline := &types.Pipeline{}
	err = utils.ParseYAML(pipelineFilename, pipeline)
	if err != nil {
		fmt.Println(err)
		return
	}

	uploadList := []FileInfo{}
	downloadList := []FileInfo{}
	resultsFolder := fmt.Sprintf("%s-pipeline", JobName)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	MigrationPod, err := objects.CreatePod(clientset,
		v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-migration", JobName),
				Namespace: "flink",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "busybox",
						Image: "busybox",
						Command: []string{
							"sh",
						},
						Args: []string{
							"-c",
							`while true; do echo \"Running Migration!\"; sleep 3600; done`,
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "migration-volume",
								MountPath: PVCMountPath,
							},
						},
					},
				},
				Volumes: []v1.Volume{
					{
						Name: "migration-volume",
						VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: fmt.Sprintf("%s-pvc", flinkCluster),
							},
						},
					},
				},
			},
		})

	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 2)

	if Migrate {
		if src := pipeline.Pipeline.Source; src != nil {
			if path, ok := (*src.Config)["path"].(string); ok {
				fileInfo, err := os.Stat(path)
				if err != nil {
					fmt.Printf("error loading file in config.path for source %s: %s\n", src.Type, err)
					return
				}

				if !fileInfo.IsDir() {
					splits := strings.Split(path, "/")
					(*src.Config)["path"] = filepath.Join(PVCMountPath, "data", splits[len(splits)-1])
				} else {
					(*src.Config)["path"] = filepath.Join(PVCMountPath, "data")
				}

				uploadList = append(uploadList, FileInfo{Src: path, Dest: filepath.Join(PVCMountPath, "data")})
			}
		}

		if sink := pipeline.Pipeline.Sink; sink != nil {
			if path, ok := (*sink.Config)["path"].(string); ok {
				splits := strings.Split(path, "/")
				var resultPath string

				if len(splits) > 1 {
					resultPath = filepath.Join(resultsFolder, splits[len(splits)-2], splits[len(splits)-1])
				} else if len(splits) == 1 {
					resultPath = filepath.Join(resultsFolder, splits[len(splits)-1])
				}

				(*sink.Config)["path"] = filepath.Join(PVCMountPath, resultPath)

				downloadList = append(downloadList, FileInfo{Src: filepath.Join(PVCMountPath, resultsFolder), Dest: path})
			}
		}

		hasResuls := false
		for _, tf := range pipeline.Pipeline.Transforms {
			if strings.HasPrefix(strings.ToLower(tf.Type), "readfrom") {
				if path, ok := (*tf.Config)["path"].(string); ok {
					fileInfo, err := os.Stat(path)
					if err != nil {
						fmt.Printf("error loading file in config.path for transform %s: %s\n", tf.Type, err)
						return
					}

					if !fileInfo.IsDir() {
						splits := strings.Split(path, "/")
						(*tf.Config)["path"] = filepath.Join(PVCMountPath, "data", splits[len(splits)-1])
					} else {
						(*tf.Config)["path"] = filepath.Join(PVCMountPath, "data")
					}

					uploadList = append(uploadList, FileInfo{Src: path, Dest: filepath.Join(PVCMountPath, "data")})
				}

			} else if strings.HasPrefix(strings.ToLower(tf.Type), "writeto") {
				if path, ok := (*tf.Config)["path"].(string); ok {
					splits := strings.Split(path, "/")
					var resultPath string

					if len(splits) > 1 {
						resultPath = filepath.Join(resultsFolder, splits[len(splits)-2], splits[len(splits)-1])
					} else if len(splits) == 1 {
						resultPath = filepath.Join(resultsFolder, splits[len(splits)-1])
					}

					(*tf.Config)["path"] = filepath.Join(PVCMountPath, resultPath)
					hasResuls = true
				}
			}
		}
		if hasResuls {
			homeDir, _ := os.UserHomeDir()
			outDir := filepath.Join(homeDir, "beamstack-pipelines", resultsFolder)
			err = os.MkdirAll(outDir, 0777)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
			downloadList = append(downloadList, FileInfo{Src: filepath.Join(PVCMountPath, resultsFolder), Dest: outDir})
		}

		fmt.Println("Performing data migration!")
		for _, file := range uploadList {
			if err := utils.MigrateFilesToContainer(
				clientset,
				types.MigrationParams{
					Pod:      *MigrationPod,
					SrcPath:  file.Src,
					DestPath: file.Dest,
				},
			); err != nil {
				fmt.Println(err)
			}
		}

		pipelineFilename, err = savePipeline(pipeline)

		if err != nil {
			fmt.Println(err)
			return
		}

		splits := strings.Split(pipelineFilename, "/")

		CleanPipelineFilename = splits[len(splits)-1]

	}

	if err := utils.MigrateFilesToContainer(
		clientset,
		types.MigrationParams{
			Pod:      *MigrationPod,
			SrcPath:  pipelineFilename,
			DestPath: PVCMountPath,
		},
	); err != nil {
		fmt.Println(err)
		return
	}
	BackOffLimit := int32(1)
	pipelineJob, err := objects.CreateJob(
		clientset,
		batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      JobName,
				Namespace: "flink",
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: &BackOffLimit,
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": JobName},
					},
					Spec: v1.PodSpec{
						RestartPolicy: "Never",
						Containers: []v1.Container{
							{
								Name:  "beam-pipeline",
								Image: "beamstackproj/beam-harness-v1_16:latest",
								// Image:   "localhost:5000/docker-ext-v1:latest",
								Command: []string{"python"},
								Args: []string{
									"-m",
									"apache_beam.yaml.main",
									fmt.Sprintf("--pipeline_spec_file=%s", filepath.Join(PVCMountPath, CleanPipelineFilename)),
									"--runner=FlinkRunner",
									fmt.Sprintf("--flink_master=%s-rest.flink.svc.cluster.local:8081", flinkCluster),
									fmt.Sprintf("--job_name=%s", JobName),
									fmt.Sprintf("--parallelism=%s", fmt.Sprintf("%d", Parallelism)),
									"--environment_type=EXTERNAL",
									"--environment_config=localhost:50000",
									"--flink_submit_uber_jar",
									"--checkpointing_interval=10000",
								},
								VolumeMounts: []v1.VolumeMount{
									{
										Name:      "migration-volume",
										MountPath: PVCMountPath,
									},
								},
							},
						},
						Volumes: []v1.Volume{
							{
								Name: "migration-volume",
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: fmt.Sprintf("%s-pvc", flinkCluster),
									},
								},
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("could not create pipeline job %s\n", err)
		return
	}

	err = os.Remove(pipelineFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Pipeline deployed!")

	if Wait {
		donChan := make(chan string)
		go objects.HandleSpecificResource(schema.GroupVersionResource{
			Group:    "batch",
			Version:  "v1",
			Resource: "jobs",
		}, pipelineJob.Name, "flink", "Complete", donChan)

		for i := range donChan {
			fmt.Println(i)
		}

		if Migrate && downloadList != nil {
			fmt.Println("migrating pipeline results!")
			for _, path := range downloadList {
				utils.MigrateFilesFromContainer(clientset,
					types.MigrationParams{
						Pod:      *MigrationPod,
						SrcPath:  path.Src,
						DestPath: path.Dest,
					},
				)
				fmt.Printf("Copied pipeline results to path:  %s", path.Dest)
			}
		}

		fmt.Println("Pipeline is done!")
		fg := metav1.DeletePropagationBackground

		clientset.BatchV1().Jobs("flink").Delete(cmd.Context(), pipelineJob.Name, metav1.DeleteOptions{PropagationPolicy: &fg})
	}

	fg := metav1.DeletePropagationBackground
	clientset.CoreV1().Pods("flink").Delete(cmd.Context(), MigrationPod.Name, metav1.DeleteOptions{PropagationPolicy: &fg})

}

func savePipeline(data interface{}) (string, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling to YAML: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(yamlData)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	return tmpFile.Name(), nil
}
