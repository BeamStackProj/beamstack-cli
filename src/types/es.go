package types

type EsNodeConfig struct {
	NodeStoreAllowMMAP bool `yaml:"node.store.allow_mmap"`
}

type EsNodeSet struct {
	Name   string       `yaml:"name"`
	Count  uint16       `yaml:"count"`
	Config EsNodeConfig `yaml:"config"`
}

type EsSpec struct {
	Version  string      `yaml:"version"`
	NodeSets []EsNodeSet `yaml:"nodeSets"`
}
