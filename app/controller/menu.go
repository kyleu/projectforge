package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/menu"
	"github.com/kyleu/projectforge/app/sandbox"
	"github.com/kyleu/projectforge/app/util"
)

func MenuFor(as *app.State) menu.Items {
	return menu.Items{
		actionMenu(as),
		menu.Separator,
		toolsMenu(),
		menu.Separator,
		&menu.Item{Key: "sandbox", Title: "Sandboxes", Description: "Playgrounds for testing new features", Icon: "star", Route: "/sandbox", Children: sandboxItems()},
		menu.Separator,
		&menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"},
	}
}

func toolsMenu() *menu.Item {
	desc := "Standalone tools for configuring projects and generating code"
	ret := &menu.Item{Key: "tools", Title: "Tools", Description: desc, Icon: "star", Route: "/tools", Children: menu.Items{
		&menu.Item{Key: "svg", Title: "SVG Icons", Description: "Configures SVG image assets", Icon: "question", Route: "/tools/svg"},
	}}
	return ret
}

func actionMenu(as *app.State) *menu.Item {
	ret := &menu.Item{Key: "actions", Title: "Actions", Description: "Run the actions of " + util.AppName, Icon: "star", Route: "/run"}
	for _, tgt := range []string{"self", "admini"} {
		for _, at := range []action.Type{action.TypePreview, action.TypeDebug, action.TypeMerge, action.TypeSlam} {
			key := tgt + "-" + at.String()
			route := fmt.Sprintf("/run/%s/%s", tgt, at.String())
			i := &menu.Item{Key: key, Title: util.ToTitle(at.String()) + " " + util.ToTitle(tgt), Icon: "star", Route: route}
			ret.Children = append(ret.Children, i)
		}
	}
	bs := &menu.Item{Key: "test-bootstrap", Title: "Bootstrap", Icon: "star", Route: "/run/test/bootstrap"}
	ret.Children = append(ret.Children, bs)
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
