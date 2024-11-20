package typescript

import (
	"fmt"
	"github.com/pkg/errors"
	"slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func All(models model.Models, enums enum.Enums, linebreak string) (file.Files, error) {
	var ret file.Files
	groups := []string{}
	groupedModels := map[string]model.Models{}
	groupedEnums := map[string]enum.Enums{}
	for _, m := range models {
		pkg := m.PackageWithGroup("")
		groups = append(groups, pkg)
		groupedModels[pkg] = append(groupedModels[pkg], m)
	}
	for _, e := range enums {
		pkg := e.PackageWithGroup("")
		groups = append(groups, pkg)
		groupedEnums[pkg] = append(groupedEnums[pkg], e)
	}

	for _, group := range groups {
		m, e := groupedModels[group], groupedEnums[group]
		imps, err := tsImports(enums, models, e, m)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing imports for group [%s]", group)
		}
		x, err := Group(group, m, e, imps, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing group [%s]", group)
		}
		ret = append(ret, x)
	}
	return ret, nil
}

func tsImports(allEnums enum.Enums, allModels model.Models, es enum.Enums, ms model.Models) ([]string, error) {
	var ret []string
	add := func(s string, args ...any) {
		s = fmt.Sprintf(s, args...)
		if !slices.Contains(ret, s) {
			ret = append(ret, s)
		}
	}
	for _, m := range ms {
		for _, col := range m.Columns {
			r, rm, _ := helper.LoadRef(col, allModels)
			if rm == nil {
				if col.Metadata != nil {
					if tsImport := col.Metadata.GetStringOpt("tsImport"); tsImport != "" {
						add(`import type {%s} from "%s";`, r.K, tsImport)
					}
				}
			} else if rm.PackageWithGroup("") != m.PackageWithGroup("") {
				add(`import type {%s} from "../%s/%s";`, r.K, rm.Camel(), rm.Camel())
			}
		}
	}
	return ret, nil
}

func Group(group string, models model.Models, enums enum.Enums, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	if group != "" {
		dir = append(dir, group)
	}
	filename := "models"
	switch {
	case len(models) == 1 && len(enums) == 0:
		filename = models[0].CamelLower()
	case len(models) == 0 && len(enums) == 1:
		filename = enums[0].Camel()
	case len(models) == 0:
		filename = "enums"
	}
	g := golang.NewGoTemplate(dir, filename+".ts")
	g.AddBlocks(tsContent(imps, enums, models)...)
	return g.Render(linebreak)
}
