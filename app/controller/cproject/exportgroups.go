package cproject

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportGroupsEdit(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.groups.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		args, err := prj.ModuleArgExport(as.Services.Projects, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Data = args.Groups

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Groups"}
		ps.Title = fmt.Sprintf("[%s] Groups", prj.Key)
		ex := model.Groups{{Key: "foo", Title: "Foo", Description: "The foos!", Icon: "star", Children: model.Groups{{Key: "bar"}}}}
		return controller.Render(rc, as, &vexport.GroupForm{Project: prj, Groups: args.Groups, Example: ex}, ps, bc...)
	})
}

func ProjectExportGroupsSave(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.groups.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		j := frm.GetStringOpt("groups")
		g := model.Groups{}
		err = util.FromJSON([]byte(j), &g)
		if err != nil {
			return "", err
		}

		err = as.Services.Projects.SaveExportGroups(as.Services.Projects.GetFilesystem(prj), g, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := "groups saved successfully"
		u := fmt.Sprintf("/p/%s/export", prj.Key)
		return controller.FlashAndRedir(true, msg, u, rc, ps)
	})
}
