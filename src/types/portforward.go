package types

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

type PortForward struct {
	RestConfig *rest.Config
	PodPort    int
	Streams    genericclioptions.IOStreams
	StopCh     <-chan struct{}
	ReadyCh    chan struct{}
	LocalPort  int
}

type PortForwardAPodRequest struct {
	PortForward
	Pod v1.Pod
}

type PortForwardASVCRequest struct {
	PortForward
	Service v1.Service
}
