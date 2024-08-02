package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func PortForwardPod(clientset *kubernetes.Clientset, req types.PortForwardAPodRequest) error {
	request := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(req.Pod.Namespace).
		Name(req.Pod.Name).
		SubResource("portforward")
	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, request.URL())
	ports := []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}
	fw, err := portforward.New(dialer, ports, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func PortForwardSvc(clientset *kubernetes.Clientset, req types.PortForwardASVCRequest) error {
	service, err := clientset.CoreV1().Services(req.Service.Namespace).Get(context.Background(), req.Service.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods(req.Service.Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labels.Set(service.Spec.Selector).String(),
	})
	if err != nil {
		return err
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for service %s", req.Service.Name)
	}

	podName := pods.Items[0].Name
	return PortForwardPod(clientset, types.PortForwardAPodRequest{
		Pod: v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podName,
				Namespace: req.Service.Namespace,
			},
		},
		PortForward: types.PortForward{
			RestConfig: req.RestConfig,
			LocalPort:  req.LocalPort,
			PodPort:    req.PodPort,
			StopCh:     req.StopCh,
			ReadyCh:    req.ReadyCh,
			Streams:    req.Streams,
		},
	})
}
