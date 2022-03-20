package controller

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/views/vproject"

	"projectforge.dev/projectforge/app/controller/cutil"

	"projectforge.dev/projectforge/app"
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

		mods := as.Services.Modules.Modules()
		gitStatus, _ := as.Services.Git.Status(ps.Context, prj, ps.Logger)

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(rc, as, &vproject.Detail{Project: prj, Modules: mods, GitResult: gitStatus}, ps, "projects", prj.Key)
	})
}

func ProjectForm(rc *fasthttp.RequestCtx) {
	act("project.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj := project.NewProject("", ".")
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

		return flashAndRedir(true, "Created project ["+prj.Title()+"]", "/p/"+prj.Key, rc, ps)
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

		return flashAndRedir(true, "Saved changes", "/p/"+prj.Key, rc, ps)
	})
}
