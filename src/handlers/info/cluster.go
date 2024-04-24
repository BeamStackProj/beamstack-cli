package info_handler

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Beamflow/beamflow-cli/src/types"
	"github.com/Beamflow/beamflow-cli/src/utils"
)

func Health() (types.ClusterHealth, error) {
	kubeconfig := utils.GetKubeConfig()

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return types.ClusterHealth{}, err
	}

	clusterHealth := types.ClusterHealth{
		Nodes: nodes.Items,
	}

	return clusterHealth, nil
}
