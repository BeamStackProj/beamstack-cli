/*
Copyright © 2024 MavenCode <opensource-dev@mavencode.com>
*/
package deploy

import (
	"fmt"
	"os"
	"strings"

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
	flinkCluster     string       = ""
	PVCMountPath     string       = "/pvc"
	JobName          string       = "beamjob"
	Parallelism      uint8        = 1
	Wait             bool         = false
	Migrate          bool         = true
	config           *rest.Config = utils.GetKubeConfig()
	pipelineFilename string
)

type FileInfo struct {
	Src  string
	Dest string
}

// infoCmd represents the info command
var PipelineCmd = &cobra.Command{
	Use:   "pipeline [FILE]",
	Short: "Deploy an Apache Beam pipeline",
	Long:  "Deploy an Apache Beam pipeline on a specified Operator.",
	Args:  cobra.ExactArgs(1),
	Run:   DeployPipeline,
}

func init() {
	PipelineCmd.Flags().StringVar(&flinkCluster, "flinkCluster", flinkCluster, "Specify the Flink cluster to deploy the Apache Beam pipeline.")
	PipelineCmd.Flags().StringVar(&PVCMountPath, "pvcMountPath", PVCMountPath, "Mount path for the Persistent Volume Claim. Note: The mount path is set to 'pvc' during cluster creation, so changing this may cause issues.")
	PipelineCmd.Flags().StringVar(&JobName, "jobname", JobName, "Specify the name of the pipeline job.")
	PipelineCmd.Flags().Uint8Var(&Parallelism, "parallelism", Parallelism, "Set the pipeline parallelism.")
	PipelineCmd.Flags().BoolVarP(&Wait, "wait", "w", Wait, "Wait for the pipeline to complete.")
	PipelineCmd.Flags().BoolVarP(&Migrate, "migrate", "m", Migrate, "Migrate data to the Kubernetes cluster. This is necessary if the pipeline is to be run on local data. Pipeline Results will also be migrated to local system if wait is true.")
}

func DeployPipeline(cmd *cobra.Command, args []string) {
	pipelineFilename = args[0]
	profile, err := utils.ValidateCluster()

	if err != nil {
		fmt.Println(err)
		return
	}
	if profile.Operators.Flink == nil {
		fmt.Println("Flink Operator not initialized on this cluster")
		return
	}
	flinkVersion := strings.Replace(profile.Operators.Flink.Version, ".", "_", 1)
	pipeline := &types.Pipeline{}
	err = utils.ParseYAML(pipelineFilename, pipeline)
	if err != nil {
		fmt.Println(err)
		return
	}

	uploadList := []FileInfo{}
	downloadList := []FileInfo{}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	MigrationPod, err := objects.CreatePod(clientset,
		v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "migration",
				Namespace: "flink",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "busybox",
						Image: "busybox",
						Command: []string{
							"sh", "-c", "'while true; do echo \"Running Migration!\"; sleep 3600; done'",
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
								ClaimName: fmt.Sprintf("%s-pvs", flinkCluster),
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
					resultPath = filepath.Join(splits[len(splits)-2], splits[len(splits)-1])
				} else if len(splits) == 1 {
					resultPath = splits[len(splits)-1]
				}

				(*sink.Config)["path"] = filepath.Join(PVCMountPath, resultPath)

				downloadList = append(downloadList, FileInfo{Src: filepath.Join(PVCMountPath, resultPath), Dest: path})
			}
		}

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
						resultPath = filepath.Join(splits[len(splits)-2], splits[len(splits)-1])
					} else if len(splits) == 1 {
						resultPath = splits[len(splits)-1]
					}

					(*tf.Config)["path"] = filepath.Join(PVCMountPath, resultPath)

					downloadList = append(downloadList, FileInfo{Src: filepath.Join(PVCMountPath, resultPath), Dest: path})
				}
			}
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
	}

	if err := utils.MigrateFilesToContainer(
		clientset,
		types.MigrationParams{
			Pod:      *MigrationPod,
			SrcPath:  pipelineFilename,
			DestPath: filepath.Join(PVCMountPath, pipelineFilename),
		},
	); err != nil {
		fmt.Println(err)
		return
	}

	pipelineJob, err := objects.CreateJob(
		clientset,
		batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      JobName,
				Namespace: "flink",
			},
			Spec: batchv1.JobSpec{
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": JobName},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:    "beam-pipeline",
								Image:   fmt.Sprintf("beamstackproj/flink-harness%s:latest", flinkVersion),
								Command: []string{"python"},
								Args: []string{
									"-m",
									"apache_beam.yaml.main",
									fmt.Sprintf("--pipeline_spec_file=%s", filepath.Join(PVCMountPath, pipelineFilename)),
									"--runner=FlinkRunner",
									fmt.Sprintf("--flink_master=%s-rest.flink.svc.cluster.local:8081", flinkCluster),
									fmt.Sprintf("--job_name=%s", JobName),
									fmt.Sprintf("--parallelism=%s", string(Parallelism)),
									"--environment_type=EXTERNAL",
									"--environment_config=localhost:50000",
									"--flink_submit_uber_jar",
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
										ClaimName: fmt.Sprintf("%s-pvs", flinkCluster),
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
	}

	err = os.Remove(pipelineFilename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Pipeline deployed!")

	if Wait {
		donChan := make(chan string, 1)
		go objects.HandleSpecificResource(schema.GroupVersionResource{
			Group:    "batch",
			Version:  "v1",
			Resource: "jobs",
		}, pipelineJob.Name, "flink", "Complete", &donChan)

		for i := range donChan {
			fmt.Println(i)
			fmt.Printf("pipeline conditions %s", pipelineJob.Status.Conditions)
		}

		if Migrate {
			fmt.Println("migrating pipeline results!")
			for _, path := range downloadList {
				utils.MigrateFilesFromContainer(clientset,
					types.MigrationParams{
						Pod:      *MigrationPod,
						SrcPath:  path.Src,
						DestPath: path.Dest,
					},
				)
			}
		}

		fmt.Println("Pipeline is done!")
		clientset.BatchV1().Jobs("flink").Delete(cmd.Context(), pipelineJob.Name, metav1.DeleteOptions{})
	}
	clientset.CoreV1().Pods("flink").Delete(cmd.Context(), MigrationPod.Name, metav1.DeleteOptions{})

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
