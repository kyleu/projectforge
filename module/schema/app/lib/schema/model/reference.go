package model

import (
	"fmt"
	"strings"

	"{{{ .Package }}}/app/util"
)

type Reference struct {
	Key          string   `json:"key"`
	TargetFields []string `json:"tgtFields"`
	SourcePkg    util.Pkg `json:"sourcePackage"`
	SourceModel  string   `json:"sourceModel"`
	SourceFields []string `json:"srcFields"`
	str          string
}

func ReferenceFromRelation(rel *Relationship, m *Model) *Reference {
	fields := make([]string, 0, len(rel.SourceFields))
	for _, x := range rel.SourceFields {
		_, col := m.Fields.Get(x)
		if col != nil {
			fields = append(fields, col.Name())
		}
	}
	str := fmt.Sprintf("%s by %s", m.PluralName(), strings.Join(fields, ", "))
	return &Reference{Key: rel.Key, TargetFields: rel.TargetFields, SourcePkg: m.Pkg, SourceModel: m.Key, SourceFields: rel.SourceFields, str: str}
}

func (r *Reference) String() string {
	return r.str
}

func (r *Reference) Debug() string {
	return fmt.Sprintf("%s: [%s] -> %s[%s]", r.Key, strings.Join(r.TargetFields, ", "), r.SourceModel, strings.Join(r.SourceFields, ", "))
}

type References []*Reference

func (s References) Get(key string) *Reference {
	for _, x := range s {
		if x.Key == key {
			return x
		}
	}
	return nil
}
