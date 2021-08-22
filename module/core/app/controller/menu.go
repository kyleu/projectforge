// Package controller $PF_IGNORE$
package controller

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/menu"
	"{{{ .Package }}}/app/sandbox"
	"{{{ .Package }}}/app/util"
)

func MenuFor(as *app.State) (menu.Items, error) {
	return menu.Items{
		&menu.Item{Key: "quickstart", Title: "Quickstart", Description: "Check out your fancy app!", Icon: "star", Route: "/quickstart"},
		menu.Separator,
		sandbox.Menu(),
		menu.Separator,
		&menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"},
	}, nil
}
