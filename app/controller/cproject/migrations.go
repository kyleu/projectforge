package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views"
)

func Migrations(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.migrations", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		ret, err := as.Services.Projects.MigrationList(ps.Context, prj, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("Migrations for [%s]", prj.Title()), ret)
		return controller.Render(r, as, &views.Debug{}, ps, "projects", prj.Key, "Migrations")
	})
}

func Migration(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.migration", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		if !prj.HasModule("migration") {
			return "", errors.Errorf("project [%s] does not include the migration module", prj.Key)
		}
		idx, err := cutil.PathInt(r, "idx")
		if err != nil {
			return "", err
		}
		migrations, err := as.Services.Projects.MigrationList(ps.Context, prj, ps.Logger)
		if err != nil {
			return "", err
		}
		if idx > len(migrations) {
			return "", errors.Errorf("project [%s] doesn't have a migration at index [%d]", prj.Key, idx)
		}
		ret := migrations[idx]
		title := "???"
		ps.SetTitleAndData(title, ret)
		return controller.Render(r, as, &views.Debug{}, ps, "Migrations||"+prj.WebPath()+"/migrations", title)
	})
}
