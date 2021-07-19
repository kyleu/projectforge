// Package controller $PF_IGNORE$
package controller

import (
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/menu"
	"$PF_PACKAGE$/app/util"
)

func MenuFor(as *app.State) (menu.Items, error) {
	return menu.Items{
		&menu.Item{Key: "quickstart", Title: "Quickstart", Description: "Check out your fancy app!", Icon: "star", Route: "/quickstart"},
		menu.Separator,
		&menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"},
	}, nil
}
