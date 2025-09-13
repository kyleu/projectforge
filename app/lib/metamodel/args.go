package metamodel

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

type Args struct {
	Config         util.ValueMap              `json:"config,omitempty"`
	ConfigFile     json.RawMessage            `json:"-"`
	Enums          enum.Enums                 `json:"enums,omitempty"`
	EnumFiles      map[string]json.RawMessage `json:"-"`
	Events         model.Events               `json:"events,omitempty"`
	EventFiles     map[string]json.RawMessage `json:"-"`
	Models         model.Models               `json:"models,omitempty"`
	ModelFiles     map[string]json.RawMessage `json:"-"`
	Groups         model.Groups               `json:"groups,omitempty"`
	GroupsFile     json.RawMessage            `json:"-"`
	Acronyms       []string                   `json:"acronyms,omitempty"`
	ExtraTypes     model.Models               `json:"extraTypes,omitempty"`
	ExtraTypesFile json.RawMessage            `json:"-"`
	Modules        []string                   `json:"-"`
	Database       string                     `json:"-"`
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

func (a *Args) DatabaseNow() string {
	switch a.Database {
	case util.DatabaseSQLite:
		return "current_timestamp"
	case util.DatabaseSQLServer:
		return "getdate()"
	default:
		return "now()"
	}
}

func (a *Args) Validate() error {
	packages := make(map[string]struct{}, len(a.Enums)+len(a.Events)+len(a.Models))
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
		packages[m.Package] = util.EmptyStruct
	}
	return nil
}

func (a *Args) Audit(m *model.Model) bool {
	return m.HasTag("audit") && lo.Contains(a.Modules, "audit")
}

func (a *Args) ApplyAcronyms(acronyms ...string) {
	a.Acronyms = acronyms
	for _, x := range a.Enums {
		x.Acronyms = acronyms
	}
	for _, x := range a.Events {
		x.SetAcronyms(acronyms...)
	}
	for _, x := range a.Models {
		x.SetAcronyms(acronyms...)
	}
}

func (a *Args) Empty() bool {
	return a == nil || (len(a.Enums) == 0 && len(a.Events) == 0 && len(a.Models) == 0)
}
