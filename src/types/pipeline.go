package types

type Pipeline struct {
	Pipeline  PipelineSpec       `yaml:"pipeline,omitempty"`
	Options   *map[string]string `yaml:"options,omitempty"`
	Providers *[]TransformSpecs  `yaml:"providers,omitempty"`
}

type PipelineSpec struct {
	Type       string           `yaml:"type,omitempty"`
	Source     *SourceSinkSpec  `yaml:"source,omitempty"`
	Transforms []TransformSpecs `yaml:"transforms,omitempty"`
	Sink       *SourceSinkSpec  `yaml:"sink,omitempty"`
}

type TransformSpecs struct {
	Type       string                  `yaml:"type,omitempty"`
	Name       *string                 `yaml:"name,omitempty"`
	Input      *string                 `yaml:"input,omitempty"`
	Config     *map[string]interface{} `yaml:"config,omitempty"`
	Transforms *map[string]interface{} `yaml:"transforms,omitempty"`
	Windowing  *map[string]interface{} `yaml:"windowing,omitempty"`
}

type SourceSinkSpec struct {
	Type   string                  `yaml:"type,omitempty"`
	Config *map[string]interface{} `yaml:"config,omitempty"`
}
