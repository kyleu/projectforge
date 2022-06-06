package data

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/lib/types"
)

func exportCol(f *Field) (*model.Column, error) {
	typ, tags, err := typeFor(f.Type)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse type for field [%s]", f.Name)
	}

	ret := &model.Column{Name: f.Name, Type: typ, PK: f.Unique, Nullable: false, Search: f.Unique, Tags: tags, HelpString: f.Description}
	return ret, nil
}

func typeFor(t string) (*types.Wrapped, []string, error) {
	switch t {
	default:
		return types.NewString(), nil, nil
	}
}
