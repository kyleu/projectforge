package cproject

import (
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportGroupsEdit(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.groups.edit", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Groups"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Groups", prj.Key), prj.ExportArgs.Groups)
		ex := model.Groups{{Key: "foo", Title: "Foo", Description: "The foos!", Icon: "star", Children: model.Groups{{Key: "bar"}}}}
		return controller.Render(r, as, &vexport.GroupForm{Project: prj, Groups: prj.ExportArgs.Groups, Example: ex}, ps, bc...)
	})
}

func ProjectExportGroupsSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.groups.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		j := frm.GetStringOpt("groups")
		g := model.Groups{}
		err = util.FromJSON([]byte(j), &g)
		if err != nil {
			return "", err
		}

		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportGroups(pfs, g)
		if err != nil {
			return "", err
		}

		msg := "groups saved successfully"
		u := fmt.Sprintf("/p/%s/export", prj.Key)
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}
