package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/views/vproject"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ProjectList(rc *fasthttp.RequestCtx) {
	act("project.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ps.Title = "Project Listing"
		ps.Data = prjs
		return render(rc, as, &vproject.List{Projects: prjs}, ps, "projects")
	})
}

func ProjectDetail(rc *fasthttp.RequestCtx) {
	act("project.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(rc, as, &vproject.Detail{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectForm(rc *fasthttp.RequestCtx) {
	act("project.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj := project.NewProject("", "")
		ps.Title = "New Project"
		ps.Data = prj
		return render(rc, as, &vproject.Edit{Project: prj}, ps, "projects", "New")
	})
}

func ProjectCreate(rc *fasthttp.RequestCtx) {
	act("project.create", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		prj := project.NewProject("", "")
		err = projectFromForm(frm, prj)
		if err != nil {
			return "", err
		}
		key := frm.GetStringOpt("key")
		prj.Key = key
		err = as.Services.Projects.Save(prj)
		if err != nil {
			return ersp("unable to save project: %+v", err)
		}

		msg := "Saved changes"
		return flashAndRedir(true, msg, "/p/"+prj.Key, rc, ps)
	})
}

func ProjectEdit(rc *fasthttp.RequestCtx) {
	act("project.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(rc, as, &vproject.Edit{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectSave(rc *fasthttp.RequestCtx) {
	act("project.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		err = projectFromForm(frm, prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.Save(prj)
		if err != nil {
			return ersp("unable to save project: %+v", err)
		}

		msg := "Saved changes"
		return flashAndRedir(true, msg, "/p/"+prj.Key, rc, ps)
	})
}
