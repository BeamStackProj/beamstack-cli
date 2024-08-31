package objects

import (
	"context"
	"fmt"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/BeamStackProj/beamstack-cli/src/utils"
)

type resourceStruct struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              interface{} `json:"spec,omitempty"`
}

func CreateObject(path string) error {

	config := utils.GetKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Example: Apply a YAML file from a URL
	err = applyYAML(dynamicClient, path)
	if err != nil {
		return fmt.Errorf("error applying yaml from url %v", err)
	}
	return nil
}

func CreateNamespace(name string) error {
	config := utils.GetKubeConfig()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreateSecret(name string, namespace string, data map[string][]byte) error {
	config := utils.GetKubeConfig()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreateDynamicResource(typeMeta metav1.TypeMeta, metaData metav1.ObjectMeta, specs interface{}, resourcetype string) error {
	config := utils.GetKubeConfig()

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	group, version, err := SplitAPIVersion(typeMeta.APIVersion)
	if err != nil {
		return err
	}
	resourceInterface := client.Resource(
		schema.GroupVersionResource{
			Group:    group,
			Version:  version,
			Resource: resourcetype,
		},
	).Namespace(metaData.Namespace)

	_, err = resourceInterface.Get(context.TODO(), metaData.Name, metav1.GetOptions{})
	if !errors.IsNotFound(err) {
		return fmt.Errorf("%s %s already exists", typeMeta.Kind, metaData.Name)
	}

	unstructuredObj, err := toUnstructured(&resourceStruct{
		TypeMeta:   typeMeta,
		ObjectMeta: metaData,
		Spec:       specs,
	})
	if err != nil {
		return err
	}

	_, err = resourceInterface.Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreatePVC(clientset *kubernetes.Clientset, name string, namespace string, size string) error {

	_, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if !errors.IsNotFound(err) {
		var userInput string
		fmt.Printf("persisten volume claim %s already exists. do you want to recreate it? y/n: ", name)
		fmt.Scanln(&userInput)

		if strings.ToLower(userInput) == "n" {
			return nil
		} else if strings.ToLower(userInput) == "y" {
			err = clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("error deleting pvc: %s", err)
			}
			time.Sleep(time.Second * 5)
		}
	}

	PVCSpec := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteMany,
			},
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(size),
				},
			},
		},
	}

	_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), &PVCSpec, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreateJob(clientset *kubernetes.Clientset, job batchv1.Job) (jobInterface *batchv1.Job, err error) {

	_, err = clientset.BatchV1().Jobs(job.Namespace).Get(context.TODO(), job.Name, metav1.GetOptions{})
	if !errors.IsNotFound(err) {
		return jobInterface, errors.NewAlreadyExists(schema.GroupResource{Group: "Batchv1", Resource: "Job"}, job.Name)
	}

	jobInterface, err = clientset.BatchV1().Jobs(job.Namespace).Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		return
	}

	return
}

func CreatePod(clientset *kubernetes.Clientset, podspec v1.Pod) (pod *v1.Pod, err error) {

	pod, err = clientset.CoreV1().Pods(podspec.Namespace).Get(context.TODO(), podspec.Name, metav1.GetOptions{})
	if !errors.IsNotFound(err) {
		return pod, nil
	}

	pod, err = clientset.CoreV1().Pods(podspec.Namespace).Create(context.TODO(), &podspec, metav1.CreateOptions{})
	if err != nil {
		return
	}
	return
}
