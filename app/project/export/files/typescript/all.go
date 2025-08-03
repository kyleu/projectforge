package typescript

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func All(models model.Models, enums enum.Enums, extraTypes model.Models, linebreak string) (file.Files, error) {
	var ret file.Files

	for _, e := range enums {
		x, err := EnumContent(e, nil, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing enum [%s]", e.Name)
		}
		ret = append(ret, x)
	}

	for _, m := range models {
		imps, err := tsModelImports(enums, models, extraTypes, m)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing imports for model [%s]", m.Name)
		}
		x, err := ModelContent(m, enums, imps, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, x)
	}

	return ret, nil
}

func tsModelImports(allEnums enum.Enums, allModels model.Models, extraTypes model.Models, m *model.Model) ([]string, error) {
	var ret []string
	add := func(s string, args ...any) {
		s = fmt.Sprintf(s, args...)
		if !slices.Contains(ret, s) {
			ret = append(ret, s)
		}
	}
	for _, col := range m.Columns {
		if e, _ := model.AsEnum(col.Type); e != nil {
			if en := allEnums.Get(e.Ref); en != nil {
				if en.PackageWithGroup("") != m.PackageWithGroup("") {
					add(`import type { %s } from "../%s/%s";`, en.Proper(), en.Camel(), en.Camel())
				} else {
					add(`import type { %s } from "./%s";`, en.Proper(), en.Camel())
				}
			}
		}
		r, rm, _ := helper.LoadRef(col, allModels, extraTypes)
		if rm == nil {
			if col.Metadata != nil {
				if tsImport := col.Metadata.GetStringOpt("tsImport"); tsImport != "" {
					add(`import type { %s } from "%s";`, r.K, tsImport)
				}
			}
		} else {
			if tsImport := rm.Config.GetStringOpt("tsImport"); tsImport != "" {
				add(`import type { %s } from "%s";`, r.K, tsImport)
			} else if rm.PackageWithGroup("") != m.PackageWithGroup("") {
				add(`import type { %s } from "../%s/%s";`, r.K, rm.Camel(), rm.Camel())
			} else {
				add(`import type { %s } from "./%s";`, r.K, rm.Camel())
			}
		}
		if col.Type.Key() == types.KeyNumeric {
			relPath := util.StringRepeat("../", len(m.Group)+1)
			add(`import type { Numeric } from "%snumeric/numeric";`, relPath)
		}
	}
	return ret, nil
}

func ModelContent(m *model.Model, allEnums enum.Enums, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, m.PackageWithGroup(""))
	filename := m.CamelLower()
	g := golang.NewGoTemplate(dir, filename+".ts")
	g.AddBlocks(tsModelContent(imps, m, allEnums)...)
	return g.Render(linebreak)
}

func EnumContent(e *enum.Enum, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, e.PackageWithGroup(""))
	filename := e.Camel()
	g := golang.NewGoTemplate(dir, filename+".ts")
	g.AddBlocks(tsEnumContent(imps, e)...)
	return g.Render(linebreak)
}
