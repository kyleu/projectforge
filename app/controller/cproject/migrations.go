package cproject

import (
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vproject"
)

func Migrations(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.migrations", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := GetProject(r, as)
		if err != nil {
			return "", err
		}
		ret, err := as.Services.Projects.MigrationList(ps.Context, prj, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("Migrations for [%s]", prj.Title()), ret)
		return controller.Render(r, as, &vproject.Migrations{Project: prj, Migrations: ret}, ps, "projects", prj.Key, "Migrations")
	})
}
