package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vproject"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ProjectList(ctx *fasthttp.RequestCtx) {
	act("project.root", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ps.Title = "Project Listing"
		ps.Data = prjs
		return render(ctx, as, &vproject.List{Projects: prjs}, ps, "projects")
	})
}

func ProjectDetail(ctx *fasthttp.RequestCtx) {
	act("project.detail", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(ctx, as, &vproject.Detail{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectEdit(ctx *fasthttp.RequestCtx) {
	act("project.edit", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(ctx, as, &vproject.Edit{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectSave(ctx *fasthttp.RequestCtx) {
	act("project.save", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}

		msg := "Saved changes"
		return flashAndRedir(true, msg, "/p/" + prj.Key, ctx, ps)
	})
}

func ProjectFileRoot(ctx *fasthttp.RequestCtx) {
	act("project.file.root", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(ctx, as, &vproject.Files{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectFile(ctx *fasthttp.RequestCtx) {
	act("project.file", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}

		pathS, err := ctxRequiredString(ctx, "path", false)
		if err != nil {
			return "", err
		}
		path := util.SplitAndTrim(pathS, "/")
		bcAppend := "||/p/" + prj.Key + "/fs"
		bc := []string{"projects", prj.Key, "Files" + bcAppend}
		for _, x := range path {
			bcAppend += "/" + x
			b := x + bcAppend
			bc = append(bc, b)
		}
		return render(ctx, as, &vproject.Files{Project: prj, Path: path}, ps, bc...)
	})
}

func getProject(ctx *fasthttp.RequestCtx, as *app.State) (*project.Project, error) {
	key, err := ctxRequiredString(ctx, "key", true)
	if err != nil {
		return nil, err
	}

	mod, err := as.Services.Projects.Get(key)
	if err != nil {
		return nil, err
	}
	return mod, nil
}
