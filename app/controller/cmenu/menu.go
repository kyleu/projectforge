// Content managed by Project Forge, see [projectforge.md] for details.
package cmenu

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

func MenuFor(
	ctx context.Context, isAuthed bool, isAdmin bool, profile *user.Profile, params filter.ParamSet, as *app.State, logger util.Logger, //nolint:revive
) (menu.Items, any, error) {
	var ret menu.Items
	var data any
	// $PF_SECTION_START(routes_start)$
	ret = append(ret,
		projectMenu(as.Services.Projects.Projects()),
		menu.Separator,
		moduleMenu(as.Services.Modules.Modules()),
		menu.Separator,
	)
	// $PF_SECTION_END(routes_start)$
	// $PF_SECTION_START(routes_end)$
	adm := &menu.Item{Key: "admin", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin"}
	ret = append(ret, docMenu(logger), menu.Separator, adm)
	if len(as.Services.Exec.Execs) > 0 {
		ret = append(ret, processMenu(as.Services.Exec.Execs))
	}
	ret = append(ret, DoctorMenu("first-aid", "/doctor"))
	const desc = "Get assistance and advice for using " + util.AppName
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: desc, Icon: "question", Route: "/about"})
	// $PF_SECTION_END(routes_end)$
	return ret, data, nil
}
