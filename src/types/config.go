package types

type Config struct {
	Version  string `json:"Version,omitempty"`
	Contexts map[string]string
}
