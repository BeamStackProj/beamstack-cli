package types

import (
	"errors"
	"time"
)

type Timestamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OperatorDetails struct {
	Version   string `json:"version"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

type Operator struct {
	Flink *OperatorDetails `json:"flink,omitempty"` // Use pointers and omitempty to make them optional
	Spark *OperatorDetails `json:"spark,omitempty"`
}

type Monitoring struct {
	Name string `json:"name"`
}

type Package struct {
	Name         string     `json:"name"`
	Type         string     `json:"type,omitempty"`
	Version      string     `json:"version,omitempty"`
	Dependencies []*Package `json:"dependencies,omitempty"`
	Url          string     `json:"url,omitempty"`
	*Timestamp
}

// type PackageDependencies struct {
// 	Name    string `json:"Name"`
// 	Type    string `json:"type,omitempty"`
// 	Version string `json:"version,omitempty"`
// 	Url     string `json:"url,omitempty"`
// }

type Profiles struct {
	Name       string      `json:"name"`
	Operators  Operator    `json:"operators"`
	Monitoring *Monitoring `json:"monitoring,omitempty"`
	Packages   []Package   `json:"packages"`
}

// Validate method to ensure only one operator is default, or none if Operators is nil
func (c Profiles) Validate() error {
	flinkDefault := c.Operators.Flink != nil && c.Operators.Flink.IsDefault
	sparkDefault := c.Operators.Spark != nil && c.Operators.Spark.IsDefault

	if flinkDefault && sparkDefault {
		return errors.New("only one of Flink or Spark can be default, not both")
	}

	if !flinkDefault && !sparkDefault && (c.Operators.Flink != nil || c.Operators.Spark != nil) {
		return errors.New("one of Flink or Spark must be default if either is specified")
	}

	return nil
}
