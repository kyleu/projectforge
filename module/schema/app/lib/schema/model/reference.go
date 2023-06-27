package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

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
	fields := lo.FilterMap(rel.SourceFields, func(x string, _ int) (string, bool) {
		_, col := m.Fields.Get(x)
		if col != nil {
			return col.Name(), true
		}
		return "", false
	})
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
	return lo.FindOrElse(s, nil, func(x *Reference) bool {
		return x.Key == key
	})
}
