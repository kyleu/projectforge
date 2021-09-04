// Package controller $PF_IGNORE$
package controller

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/menu"
	"{{{ .Package }}}/app/sandbox"
	"{{{ .Package }}}/app/util"
)

func MenuFor(ctx context.Context, isAuthed bool, isAdmin bool, as *app.State) (menu.Items, error) {
	ret := menu.Items{
		&menu.Item{Key: "quickstart", Title: "Quickstart", Description: "Check out your fancy app!", Icon: "star", Route: "/quickstart"},
		menu.Separator,
	}
	if (isAdmin) {
		sandbox.Menu(),
		menu.Separator,
		ret = append(ret, &menu.Item{Key: "settings", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin/settings"})
	}
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"})
	return ret, nil
}
