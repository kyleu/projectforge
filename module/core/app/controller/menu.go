package controller

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/menu"{{{ if .HasModule "sandbox" }}}
	"{{{ .Package }}}/app/lib/sandbox"{{{ end }}}
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

func MenuFor(ctx context.Context, isAuthed bool, isAdmin bool, as *app.State) (menu.Items, error) {
	ctx, span, logger := telemetry.StartSpan(ctx, "menu:generate", nil)
	defer span.Complete()
	_ = logger

	var ret menu.Items
	// $PF_SECTION_START(routes_start)$
	// $PF_SECTION_END(routes_start)${{{ if.HasModule "export" }}}
	// $PF_INJECT_START(codegen)$
	// $PF_INJECT_END(codegen)${{{ end }}}
	// $PF_SECTION_START(routes_end)$
	if isAdmin {
		admin := &menu.Item{Key: "admin", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin"}
		ret = append(ret, {{{ if .HasModule "sandbox" }}}sandbox.Menu(ctx), {{{ end }}}menu.Separator, admin{{{ if .HasModule "docbrowse" }}}, menu.Separator, docMenu(ctx, as){{{ end }}})
	}
	const aboutDesc = "Get assistance and advice for using " + util.AppName
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: aboutDesc, Icon: "question", Route: "/about"})
	// $PF_SECTION_END(routes_end)$
	return ret, nil
}
