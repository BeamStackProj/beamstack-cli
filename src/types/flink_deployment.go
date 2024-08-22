package types

import (
	v1 "k8s.io/api/core/v1"
)

type FlinkDeploymentSpec struct {
	Image              *string             `yaml:"image"`
	ImagePullPolicy    string              `yaml:"imagePullPolicy"`
	FlinkVersion       string              `yaml:"flinkVersion"`
	FlinkConfiguration map[string]string   `yaml:"flinkConfiguration"`
	ServiceAccount     string              `yaml:"serviceAccount"`
	PodTemplate        *v1.PodTemplateSpec `yaml:"podTemplate,omitempty"`
	JobManager         JobManagerSpec      `yaml:"jobManager"`
	TaskManager        TaskManagerSpec     `yaml:"taskManager"`
	Job                *JobSpec            `yaml:"job,omitempty"`
	Mode               *string             `yaml:"mode,omitempty"`
}

type JobManagerSpec struct {
	Replicas    uint8               `yaml:"replicas,omitempty"`
	Resource    Resource            `yaml:"resource"`
	PodTemplate *v1.PodTemplateSpec `yaml:"podTemplate,omitempty"`
}

type TaskManagerSpec struct {
	Replicas    uint8               `yaml:"replicas"`
	Resource    Resource            `yaml:"resource"`
	PodTemplate *v1.PodTemplateSpec `yaml:"podTemplate,omitempty"`
}

type Resource struct {
	Memory string `yaml:"memory"`
	CPU    string `yaml:"cpu"`
}

type JobSpec struct {
	JarFile       string           `yaml:"jarFile"`
	ClassName     string           `yaml:"className,omitempty"`
	Args          []string         `yaml:"args,omitempty"`
	Parallelism   uint8            `yaml:"parallelism,omitempty"`
	RestartPolicy v1.RestartPolicy `yaml:"restartPolicy,omitempty"`
}
