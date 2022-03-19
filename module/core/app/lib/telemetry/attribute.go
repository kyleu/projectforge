package telemetry

import (
	"go.opentelemetry.io/otel/attribute"

	"{{{ .Package }}}/app/util"
)

type Attribute struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (a *Attribute) ToOT() attribute.KeyValue {
	switch t := a.Value.(type) {
	case string:
		return attribute.String(a.Key, t)
	case int:
		return attribute.Int(a.Key, t)
	case bool:
		return attribute.Bool(a.Key, t)
	default:
		return attribute.String(a.Key, util.ToJSON(t))
	}
}
