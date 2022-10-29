package helper

import (
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func SpecialImports(g *golang.File, cols model.Columns, pkg string, enums enum.Enums) error {
	for _, col := range cols {
		switch col.Type.Key() {
		case types.KeyReference:
			ref, err := model.AsRef(col.Type)
			if err != nil {
				return err
			}
			if ref.Pkg.Last() != pkg {
				g.AddImport(golang.NewImport(golang.ImportTypeApp, ref.Pkg.ToPath()))
			}
		case types.KeyEnum:
			e, err := model.AsEnumInstance(col.Type, enums)
			if err != nil {
				return err
			}
			if e.PackageWithGroup("") != pkg {
				g.AddImport(AppImport(e.PackageWithGroup("app/")))
			}
		}
	}
	return nil
}
