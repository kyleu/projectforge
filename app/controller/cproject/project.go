package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		prj.ExportArgs, _ = prj.ModuleArgExport(as.Services.Projects, ps.Logger)
		mods := as.Services.Modules.Modules()
		execs := as.Services.Exec.Execs.GetByKey(prj.Key)
		fs, _ := as.Services.Projects.GetFilesystem(prj)
		validation := project.Validate(prj, fs, as.Services.Modules.Deps(), as.Services.Modules.Dangerous())
		ps.SetTitleAndData(fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key), prj)
		page := &vproject.Detail{Project: prj, Modules: mods, Execs: execs, Validation: validation}
		return controller.Render(w, r, as, page, ps, "projects", prj.Key)
	})
}

func ProjectForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj := project.NewProject("", ".")
		ps.SetTitleAndData("New Project", prj)
		return controller.Render(w, r, as, &vproject.Edit{Project: prj}, ps, "projects", "New")
	})
}

func ProjectCreate(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.create", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		prj := project.NewProject("", "")
		err = projectFromForm(frm, prj)
		if err != nil {
			return "", err
		}
		key := frm.GetStringOpt("key")
		if key == "" {
			return "", errors.New("Must provide a non-empty value for [key]")
		}
		prj.Key = key

		curr, _ := as.Services.Projects.Get(prj.Key)
		if curr != nil {
			if prj.Path == "" {
				prj.Path = curr.Path
			}
			if prj.Path != curr.Path {
				return "", errors.Errorf("Path cannot change to [%s], original project is in [%s]", prj.Path, curr.Path)
			}
			if prj.Path == "" || prj.Path == "." {
				return "", errors.New("You can't recreate the default project; instead, edit it through the UI")
			}
		}

		err = as.Services.Projects.Save(prj, ps.Logger)
		if err != nil {
			return controller.ERsp("unable to save project: %+v", err)
		}
		return controller.FlashAndRedir(true, "Created project ["+prj.Title()+"]", prj.WebPath(), w, ps)
	})
}

func ProjectEdit(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.edit", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("Edit %s (project %s)", prj.Title(), prj.Key), prj)
		return controller.Render(w, r, as, &vproject.Edit{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		err = projectFromForm(frm, prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.Save(prj, ps.Logger)
		if err != nil {
			return controller.ERsp("unable to save project: %+v", err)
		}

		return controller.FlashAndRedir(true, "Saved changes", prj.WebPath(), w, ps)
	})
}
