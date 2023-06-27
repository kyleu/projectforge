package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Relationship struct {
	Key          string   `json:"key"`
	SourceFields []string `json:"srcFields"`
	TargetPkg    util.Pkg `json:"tgtPackage"`
	TargetModel  string   `json:"tgtModel"`
	TargetFields []string `json:"tgtFields"`
}

func (r *Relationship) String() string {
	return fmt.Sprintf("%s: [%s] -> %s[%s]", r.Key, strings.Join(r.SourceFields, ", "), r.TargetModel, strings.Join(r.TargetFields, ", "))
}

func (r *Relationship) Path() string {
	return strings.Join(append(r.TargetPkg, r.TargetModel), "/")
}

type Relationships []*Relationship

func (s Relationships) Get(key string) *Relationship {
	return lo.FindOrElse(s, nil, func(x *Relationship) bool {
		return x.Key == key
	})
}

func (m *Model) ApplicableRelations(key string) Relationships {
	ret := Relationships{}
	lo.ForEach(m.Relationships, func(r *Relationship, _ int) {
		lo.ForEach(r.SourceFields, func(sf string, _ int) {
			if sf == key {
				ret = append(ret, r)
			}
		})
	})
	return ret
}
