// Package controller $PF_IGNORE$
package controller

import (
	"context"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/menu"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/sandbox"
	"github.com/kyleu/projectforge/app/util"
)

func MenuFor(ctx context.Context, isAuthed bool, isAdmin bool, as *app.State) (menu.Items, error) {
	var ret menu.Items
	if isAuthed {
		ret = append(ret,
			projectMenu(as.Services.Projects.Projects()),
			menu.Separator,
			moduleMenu(as.Services.Modules.Modules()),
			menu.Separator,
		)
	}
	if isAdmin {
		ret = append(ret,
			sandbox.Menu(),
			menu.Separator,
			&menu.Item{Key: "settings", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin/settings"},
			itemFor(action.TypeDoctor, "first-aid", "/doctor"),
		)
	}
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"})
	return ret, nil
}

func itemFor(t action.Type, i string, r string) *menu.Item {
	return &menu.Item{Key: t.Key, Title: t.Title, Description: t.Description, Icon: i, Route: r}
}

func projectMenu(prjs project.Projects) *menu.Item {
	ret := &menu.Item{Key: "projects", Title: "Projects", Description: "View all of the projects managed by this application", Icon: "code", Route: "/p"}
	for _, prj := range prjs {
		key := prj.Key
		i := &menu.Item{Key: key, Title: prj.Title(), Icon: prj.SafeIcon(), Route: "/p/" + prj.Key}
		ret.Children = append(ret.Children, i)
	}
	return ret
}

func moduleMenu(mods module.Modules) *menu.Item {
	ret := &menu.Item{Key: "modules", Title: "Modules", Description: "View all of the modules managed by this application", Icon: "archive", Route: "/m"}
	for _, mod := range mods {
		key := mod.Key
		i := &menu.Item{Key: key, Title: mod.Name, Icon: mod.SafeIcon(), Route: "/m/" + mod.Key}
		ret.Children = append(ret.Children, i)
	}
	return ret
}
