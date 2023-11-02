package helper

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

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
				ret = append(ret, AppImport(e.PackageWithGroup("app")))
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
			e, _ := model.AsEnumInstance(t, enums)
			ep := e.PackageWithGroup("")
			if ep != pkg {
				ret = append(ret, AppImport("app/"+ep))
			}
		case types.KeyList:
			if x := (t.(*types.Wrapped)).ListType(); t != nil && x.Key() == types.KeyEnum {
				e, _ := model.AsEnumInstance(x, enums)
				ep := e.PackageWithGroup("")
				if ep != pkg {
					ret = append(ret, AppImport("app/"+ep))
				}
			}
		}
	}
	return ret, nil
}
