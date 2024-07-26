package objects

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/BeamStackProj/beamstack-cli/src/utils"
)

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
