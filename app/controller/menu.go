package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/menu"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/sandbox"
	"github.com/kyleu/projectforge/app/util"
)

func MenuFor(as *app.State) (menu.Items, error) {
	return menu.Items{
		&menu.Item{Key: "root", Title: "Main Project", Description: "Details for the default project", Icon: "star", Route: "/p"},
		menu.Separator,
		projectMenu(as.Services.Projects.Projects()),
		menu.Separator,
		toolsMenu(),
		menu.Separator,
		&menu.Item{Key: "sandbox", Title: "Sandboxes", Description: "Playgrounds for testing new features", Icon: "star", Route: "/sandbox", Children: sandboxItems()},
		menu.Separator,
		&menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"},
	}, nil
}

func projectMenu(prjs project.Projects) *menu.Item {
	ret := &menu.Item{Key: "projects", Title: "Projects", Description: "View all of the projects managed by this application", Icon: "star", Route: "/p"}
	for _, prj := range prjs {
		key := prj.Key
		i := &menu.Item{Key: key, Title: prj.Name, Icon: "star", Route: "/p/" + prj.Key}
		ret.Children = append(ret.Children, i)
	}
	return ret
}

func toolsMenu() *menu.Item {
	desc := "Standalone tools for configuring projects and generating code"
	ret := &menu.Item{Key: "tools", Title: "Tools", Description: desc, Icon: "star", Route: "/tools", Children: menu.Items{
		&menu.Item{Key: "svg", Title: "SVG Icons", Description: "Configures SVG image assets", Icon: "question", Route: "/tools/svg"},
	}}
	return ret
}

func sandboxItems() menu.Items {
	ret := make(menu.Items, 0, len(sandbox.AllSandboxes))
	for _, s := range sandbox.AllSandboxes {
		ret = append(ret, &menu.Item{
			Key:         s.Key,
			Title:       s.Title,
			Icon:        s.Icon,
			Description: fmt.Sprintf("Sandbox [%s]", s.Key),
			Route:       fmt.Sprintf("/sandbox/%s", s.Key),
		})
	}
	return ret
}
