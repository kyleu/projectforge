// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"context"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/lib/menu"
	"github.com/kyleu/projectforge/app/lib/sandbox"
	"github.com/kyleu/projectforge/app/lib/telemetry"
	"github.com/kyleu/projectforge/app/util"
)

func MenuFor(ctx context.Context, isAuthed bool, isAdmin bool, as *app.State) (menu.Items, error) {
	ctx, span, _ := telemetry.StartSpan(ctx, "menu:generate", nil)
	defer span.Complete()

	var ret menu.Items
	// $PF_SECTION_START(routes_start)$
	if isAuthed {
		ret = append(ret,
			projectMenu(ctx, as.Services.Projects.Projects()),
			menu.Separator,
			moduleMenu(as.Services.Modules.Modules()),
			menu.Separator,
		)
	}
	// $PF_SECTION_END(routes_start)$
	// $PF_SECTION_START(routes_end)$
	if isAdmin {
		ret = append(ret,
			sandbox.Menu(ctx),
			menu.Separator,
			&menu.Item{Key: "admin", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin"},
			DoctorMenu("first-aid", "/doctor"),
		)
	}
	const desc = "Get assistance and advice for using " + util.AppName
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: desc, Icon: "question", Route: "/about"})
	// $PF_SECTION_END(routes_end)$
	return ret, nil
}
