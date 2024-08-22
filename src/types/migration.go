package types

import v1 "k8s.io/api/core/v1"

type MigrationParams struct {
	Pod           v1.Pod
	SrcPath       string
	DestPath      string
	ContainerName *string
}
