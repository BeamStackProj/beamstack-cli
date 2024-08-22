package utils

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeConfig     *rest.Config
	kubeConfigOnce sync.Once
)

func GetKubeConfig() *rest.Config {
	kubeConfigOnce.Do(func() {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

		flag.Parse()
		var err error
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	})

	return kubeConfig
}

func GetCurrentContext() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %v", err)
	}

	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return "", fmt.Errorf("kubeconfig file does not exist at %s", kubeconfigPath)
	} else if err != nil {
		return "", fmt.Errorf("error checking kubeconfig file: %v", err)
	}

	// Load the kubeconfig file
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return "", fmt.Errorf("error loading kubeconfig file: %v", err)
	}

	// Get the current context
	currentContext := config.CurrentContext
	if currentContext == "" {
		return "", fmt.Errorf("Current Context not found, please initialize this cluster")
	}

	return currentContext, nil
}

func GetNodeEndpoints() ([]interface{}, error) {
	config := GetKubeConfig()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	var endpoints []interface{}
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeInternalIP || address.Type == v1.NodeExternalIP {
				endpoints = append(endpoints, address.Address)
			}
		}
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no node endpoints found")
	}

	return endpoints, nil
}
