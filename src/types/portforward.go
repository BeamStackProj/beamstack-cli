package types

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type PortForward struct {
	PodPort   uint16
	Streams   genericclioptions.IOStreams
	StopCh    <-chan struct{}
	ReadyCh   chan struct{}
	LocalPort uint16
}

type PortForwardAPodRequest struct {
	PortForward
	Pod v1.Pod
}

type PortForwardASVCRequest struct {
	PortForward
	Service v1.Service
}
