package cmenu

import (
	"context"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
)

func projectMenu(_ context.Context, prjs project.Projects) *menu.Item {
	ret := &menu.Item{Key: "projects", Title: "Projects", Description: "View all of the projects managed by this application", Icon: "code", Route: "/p"}
	lo.ForEach(prjs, func(prj *project.Project, _ int) {
		if prj == nil {
			return
		}
		key := prj.Key
		i := &menu.Item{Key: key, Title: prj.Title(), Description: prj.DescriptionSafe(), Icon: prj.IconSafe(), Route: prj.WebPath()}
		ret.Children = append(ret.Children, i)
	})
	return ret
}

func moduleMenu(_ context.Context, mods module.Modules) *menu.Item {
	ret := &menu.Item{Key: "modules", Title: "Modules", Description: "View all of the modules managed by this application", Icon: "archive", Route: "/m"}
	lo.ForEach(mods, func(mod *module.Module, _ int) {
		key := mod.Key
		i := &menu.Item{Key: key, Title: mod.Name, Description: mod.Description, Icon: mod.IconSafe(), Route: "/m/" + mod.Key}
		ret.Children = append(ret.Children, i)
	})
	return ret
}
