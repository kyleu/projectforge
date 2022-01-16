package schema

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/schema/model"
	"{{{ .Package }}}/app/lib/schema/types"
)

const (
	LevelInfo = iota
	LevelWarn
	LevelError
)

type ValidationMessage struct {
	Category string `json:"category,omitempty"`
	ModelKey string `json:"modelKey,omitempty"`
	Message  string `json:"message,omitempty"`
	Level    int    `json:"level,omitempty"`
}

type ValidationResult struct {
	Schema   string              `json:"schema,omitempty"`
	Messages []ValidationMessage `json:"messages,omitempty"`
	Duration int64               `json:"duration,omitempty"`
}

func (v *ValidationResult) log(category string, modelKey string, msg string, level int) {
	v.Messages = append(v.Messages, ValidationMessage{Category: category, ModelKey: modelKey, Message: msg, Level: level})
}

func validateSchema(n string, s *Schema) *ValidationResult {
	r := &ValidationResult{Schema: n}
	for _, m := range s.Models {
		r = validateModel(r, s, m)
	}
	return r
}

func validateModel(r *ValidationResult, s *Schema, m *model.Model) *ValidationResult {
	encountered := map[string]bool{}
	for _, f := range m.Fields {
		if encountered[f.Key] {
			msg := fmt.Sprintf("%s [%s] field [%s] appears twice", m.Type.String(), m.Key, f.Key)
			r.log(m.Type.Key, m.Key, msg, LevelError)
		}
		encountered[f.Key] = true
	}
	for _, v := range m.Fields {
		validateType(r, s, "model", m.Key, v.Key, v.Type)
	}
	return r
}

func validateType(r *ValidationResult, s *Schema, mType string, mKey string, fKey string, f types.Type) {
	switch t := f.(type) {
	case *types.Wrapped:
		validateType(r, s, mType, mKey, fKey, t.T)
	case *types.Unknown:
		r.log(mType, mKey, fmt.Sprintf("field [%s] has unknown type [%s]", fKey, t.X), LevelWarn)
	case *types.Error:
		r.log(mType, mKey, fmt.Sprintf("field [%s] has error: %s", fKey, t.Message), LevelWarn)
	case *types.Option:
		validateType(r, s, mType, mKey, fKey, t.V)
	case *types.List:
		validateType(r, s, mType, mKey, fKey, t.V)
	case *types.Range:
		validateType(r, s, mType, mKey, fKey, t.V)
	case *types.Map:
		validateType(r, s, mType, mKey, fKey, t.K)
		validateType(r, s, mType, mKey, fKey, t.V)
	case *types.Reference:
		if s.Models.Get(t.Pkg, t.K) == nil && s.Scalars.Get(t.Pkg, t.K) == nil {
			pkg := strings.Join(t.Pkg, ".")
			msg := fmt.Sprintf("field [%s] has reference to unknown type [%s::%s]", fKey, pkg, t.K)
			r.log(mType, mKey, msg, LevelWarn)
		}

	default:
		if fKey == "" {
			r.log(mType, mKey, "field has an empty key", LevelError)
		}
	}
}

func (s *Schema) CreateReferences() error {
	for _, src := range s.Models {
		if src.References != nil {
			return errors.New("double call of CreateReferences")
		}
		for _, tgt := range s.Models {
			for _, rel := range tgt.Relationships {
				if rel.TargetPkg.Equals(src.Pkg) && rel.TargetModel == src.Key {
					err := src.AddReference(model.ReferenceFromRelation(rel, tgt))
					if err != nil {
						return errors.Wrapf(err, "unable to add reference to [%s]", src.String())
					}
				}
			}
		}
	}
	return nil
}
