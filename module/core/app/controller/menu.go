package controller

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/lib/sandbox"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

func MenuFor(ctx context.Context, isAuthed bool, isAdmin bool, as *app.State) (menu.Items, error) {
	ctx, span, _ := telemetry.StartSpan(ctx, "menu:generate", nil)
	defer span.Complete()

	var ret menu.Items
	// $PF_SECTION_START(routes_start)$
	ret = append(ret,
		&menu.Item{Key: "quickstart", Title: "Quickstart", Description: "Check out your fancy app!", Icon: "star", Route: "/quickstart"},
		menu.Separator,
	)
	// $PF_SECTION_END(routes_start)${{{ if.HasModule "export" }}}
	// $PF_INJECT_START(codegen)$
	// $PF_INJECT_END(codegen)${{{ end }}}
	// $PF_SECTION_START(routes_end)$
	if isAdmin {
		admin := &menu.Item{Key: "admin", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin"}
		ret = append(ret, sandbox.Menu(ctx), menu.Separator, admin)
	}
	const aboutDesc = "Get assistance and advice for using " + util.AppName
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: aboutDesc, Icon: "question", Route: "/about"})
	// $PF_SECTION_END(routes_end)$
	return ret, nil
}
