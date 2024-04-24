package types

import (
	v1 "k8s.io/api/core/v1"
)

type ClusterHealth struct {
	Nodes []v1.Node
}
