package types

type EsNodeSet struct {
	Name   string      `yaml:"name"`
	Count  uint16      `yaml:"count"`
	Config map[string]interface{} `yaml:"config"`
}

type EsSpec struct {
	Version  string      `yaml:"version"`
	NodeSets []EsNodeSet `yaml:"nodeSets"`
}
