package cmenu

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

func MenuFor(
	ctx context.Context, as *app.State, isAuthed bool, isAdmin bool, profile *user.Profile, params filter.ParamSet, logger util.Logger,
) (menu.Items, any, error) {
	ctx, sp, _ := telemetry.StartSpan(ctx, "menu", logger)
	defer sp.Complete()
	var ret menu.Items
	var data any
	// $PF_SECTION_START(menu)$
	ret = append(ret, projectMenu(ctx, as.Services.Projects.Projects())...)
	ret = append(ret, menu.Separator, moduleMenu(ctx, as.Services.Modules.ModulesVisible()), menu.Separator)
	if len(as.Services.Exec.Execs) > 0 {
		ret = append(ret, processMenu(as.Services.Exec.Execs))
	}
	adm := &menu.Item{Key: "admin", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin"}
	ret = append(ret, DoctorMenu("first-aid", "/doctor"), DocsMenu(logger), adm)
	const desc = "Get assistance and advice for using " + util.AppName
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: desc, Icon: "question", Route: "/about"})
	// $PF_SECTION_END(menu)$
	return ret, data, nil
}
