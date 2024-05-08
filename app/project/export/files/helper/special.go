package helper

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func SimpleJSON(database string) bool {
	return database == util.DatabaseSQLite || database == util.DatabaseSQLServer
}

func SpecialImports(cols model.Columns, pkg string, enums enum.Enums) (golang.Imports, error) {
	var ret golang.Imports
	for _, col := range cols {
		switch col.Type.Key() {
		case types.KeyReference:
			ref, err := model.AsRef(col.Type)
			if err != nil {
				return nil, err
			}
			split := strings.Split(pkg, "/")
			if ref.Pkg.Last() != split[len(split)-1] {
				ret = append(ret, golang.NewImport(golang.ImportTypeApp, ref.Pkg.ToPath()))
			}
		case types.KeyEnum:
			e, err := model.AsEnumInstance(col.Type, enums)
			if err != nil {
				return nil, err
			}
			if e.PackageWithGroup("") != pkg {
				ret = append(ret, AppImport(e.PackageWithGroup("")))
			}
		}
	}
	return ret, nil
}

func EnumImports(ts types.Types, pkg string, enums enum.Enums) (golang.Imports, error) {
	var ret golang.Imports
	for _, t := range ts {
		switch t.Key() {
		case types.KeyEnum:
			e, err := model.AsEnumInstance(t, enums)
			if err != nil {
				return nil, err
			}
			ep := e.PackageWithGroup("")
			if ep != pkg {
				ret = append(ret, AppImport(ep))
			}
		case types.KeyList:
			w, _ := t.(*types.Wrapped)
			if x := w.ListType(); t != nil && x.Key() == types.KeyEnum {
				e, err := model.AsEnumInstance(x, enums)
				if err != nil {
					return nil, errors.Wrapf(err, "unable to find list's enum [%s]", x.EnumKey())
				}
				ep := e.PackageWithGroup("")
				if ep != pkg {
					ret = append(ret, AppImport(ep))
				}
			}
		}
	}
	return ret, nil
}
