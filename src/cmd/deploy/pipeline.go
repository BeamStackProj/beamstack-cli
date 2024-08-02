/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	flinkCluster string = ""
	PVCMountPath string = "/pvc"
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
	PipelineCmd.Flags().StringVar(&flinkCluster, "flink-cluster", flinkCluster, "Specify the Flink cluster to deploy the Apache Beam pipeline")

}

func DeployPipeline(cmd *cobra.Command, args []string) {
	pipeline := &types.Pipeline{}
	err := utils.ParseYAML(args[0], pipeline)
	if err != nil {
		fmt.Println(err)
		return
	}

	uploadList := []FileInfo{}
	downloadList := []FileInfo{}

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
	config := utils.GetKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	pod, err := objects.CreatePod(clientset,
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

	fmt.Println("Performing data migration!")
	for _, file := range uploadList {
		if err := utils.MigrateFilesToContainer(
			clientset,
			types.MigrationParams{
				Pod:      *pod,
				SrcPath:  file.Src,
				DestPath: file.Dest,
			},
		); err != nil {
			fmt.Println(err)
		}
	}

	// TODO
	// move pipeline yaml file to pvc
	// create job
	// monitory job
	// move file from pvc to local system
	// delete migration pod
}

// func SaveStructToYAML(filename string, data interface{}) error {
// 	yamlData, err := yaml.Marshal(data)
// 	if err != nil {
// 		return fmt.Errorf("error marshalling to YAML: %w", err)
// 	}

// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("error creating file: %w", err)
// 	}
// 	defer file.Close()

// 	_, err = file.Write(yamlData)
// 	if err != nil {
// 		return fmt.Errorf("error writing to file: %w", err)
// 	}

// 	return nil
// }
