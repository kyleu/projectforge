package model

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/util"
)

const SQLServer = "sqlserver"

type Args struct {
	Config     util.ValueMap              `json:"config,omitempty"`
	ConfigFile json.RawMessage            `json:"-"`
	Enums      enum.Enums                 `json:"enums,omitempty"`
	EnumFiles  map[string]json.RawMessage `json:"-"`
	Models     Models                     `json:"models,omitempty"`
	ModelFiles map[string]json.RawMessage `json:"-"`
	Groups     Groups                     `json:"groups,omitempty"`
	GroupsFile json.RawMessage            `json:"-"`
	Modules    []string                   `json:"-"`
}

func (a *Args) HasModule(key string) bool {
	return lo.Contains(a.Modules, key)
}

func (a *Args) DBRef() string {
	if a.HasModule("readonlydb") {
		return "dbRead"
	}
	return "db"
}

func (a *Args) Validate() error {
	packages := make(map[string]struct{}, len(a.Models))
	err := a.Models.Validate(a.Modules, a.Groups)
	if err != nil {
		return err
	}
	for _, m := range a.Models {
		for _, rel := range m.Relations {
			relTable := a.Models.Get(rel.Table)
			if relTable == nil {
				return errors.Errorf("relation [%s] refers to missing table [%s]", rel.Name, rel.Table)
			}
			for _, t := range rel.Tgt {
				if relTable.Columns.Get(t) == nil {
					return errors.Errorf("relation [%s] references missing target column [%s]", rel.Name, t)
				}
			}
		}
		if _, ok := packages[m.Package]; ok {
			return errors.Wrap(err, "multiple models are in package ["+m.Package+"]")
		}
		packages[m.Package] = struct{}{}
	}
	return nil
}

func (a *Args) Database() string {
	if a.HasModule("postgres") {
		return "postgres"
	}
	if a.HasModule("sqlite") {
		return "sqlite"
	}
	if a.HasModule(SQLServer) {
		return SQLServer
	}
	return util.KeyUnknown
}

func (a *Args) Audit(m *Model) bool {
	return m.HasTag("audit") && lo.Contains(a.Modules, "audit")
}
