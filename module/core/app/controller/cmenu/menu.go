package cmenu

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/menu"{{{ if .HasModule "sandbox" }}}
	"{{{ .Package }}}/app/lib/sandbox"{{{ end }}}
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

func MenuFor(
	ctx context.Context, isAuthed bool, isAdmin bool, profile *user.Profile, params filter.ParamSet, as *app.State, logger util.Logger, //nolint:revive
) (menu.Items, any, error) {
	var ret menu.Items
	var data any
	// $PF_SECTION_START(menu)$ {{{ if .HasExportModels }}}{{{ if .HasAccount }}}
	if isAdmin {
		ret = append(ret, generatedMenu()...)
	}{{{ else }}}
	ret = append(ret, generatedMenu()...){{{ end }}}{{{ end }}}
	// This is your menu, feel free to customize it
	admin := &menu.Item{Key: "admin", Title: "Settings", Description: "System-wide settings and preferences", Icon: "cog", Route: "/admin"}
	ret = append(ret, {{{ if .HasModule "graphql" }}}graphQLMenu(ctx, as.GraphQL), menu.Separator, {{{end}}}{{{ if .HasModule "sandbox" }}}sandbox.Menu(ctx), menu.Separator, {{{ end }}}admin{{{ if .HasModule "docbrowse" }}}, menu.Separator, docMenu(logger){{{ end }}})
	const aboutDesc = "Get assistance and advice for using " + util.AppName
	ret = append(ret, &menu.Item{Key: "about", Title: "About", Description: aboutDesc, Icon: "question", Route: "/about"})
	// $PF_SECTION_END(menu)$
	return ret, data, nil
}
