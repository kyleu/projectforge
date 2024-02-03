package cmenu

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
)

func projectMenu(prjs project.Projects) *menu.Item {
	ret := &menu.Item{Key: "projects", Title: "Projects", Description: "View all of the projects managed by this application", Icon: "code", Route: "/p"}
	lo.ForEach(prjs, func(prj *project.Project, _ int) {
		key := prj.Key
		i := &menu.Item{Key: key, Title: prj.Title(), Description: prj.DescriptionSafe(), Icon: prj.IconSafe(), Route: prj.WebPath()}
		ret.Children = append(ret.Children, i)
	})
	return ret
}

func moduleMenu(mods module.Modules) *menu.Item {
	ret := &menu.Item{Key: "modules", Title: "Modules", Description: "View all of the modules managed by this application", Icon: "archive", Route: "/m"}
	lo.ForEach(mods, func(mod *module.Module, _ int) {
		key := mod.Key
		i := &menu.Item{Key: key, Title: mod.Name, Description: mod.Description, Icon: mod.IconSafe(), Route: "/m/" + mod.Key}
		ret.Children = append(ret.Children, i)
	})
	return ret
}
