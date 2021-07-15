package controller

import (
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/menu"
	"$PF_PACKAGE$/app/util"
)

func MenuFor(as *app.State) menu.Items {
	return menu.Items{
		menu.Separator,
		&menu.Item{Key: "about", Title: "About", Description: "Get assistance and advice for using " + util.AppName, Icon: "question", Route: "/about"},
	}
}
