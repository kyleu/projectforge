package typescript

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
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
		x, err := Group(group, m, e, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing group [%s]", group)
		}
		ret = append(ret, x)
	}
	return ret, nil
}

func Group(group string, models model.Models, enums enum.Enums, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	if group != "" {
		dir = append(dir, group)
	}
	g := golang.NewGoTemplate(dir, "models.ts")
	g.AddBlocks(tsContent(enums, models)...)
	return g.Render(linebreak)
}
