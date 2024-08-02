package types

type Pipeline struct {
	Pipeline  PipelineSpec       `json:"pipeline,omitempty"`
	Options   *map[string]string `json:"options,omitempty"`
	Providers *[]TransformSpecs  `json:"providers,omitempty"`
}

type PipelineSpec struct {
	Type       string           `json:"type,omitempty"`
	Source     *SourceSinkSpec  `json:"source,omitempty"`
	Transforms []TransformSpecs `json:"transforms,omitempty"`
	Sink       *SourceSinkSpec  `json:"sink,omitempty"`
}

type TransformSpecs struct {
	Type       string                  `json:"type,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Input      *string                 `json:"input,omitempty"`
	Config     *map[string]interface{} `json:"config,omitempty"`
	Transforms *map[string]interface{} `json:"transforms,omitempty"`
	Windowing  *map[string]interface{} `json:"windowing,omitempty"`
}

type SourceSinkSpec struct {
	Type   string                  `json:"type,omitempty"`
	Config *map[string]interface{} `json:"config,omitempty"`
}
