package cproject

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/stats"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectFileRoot(rc *fasthttp.RequestCtx) {
	controller.Act("project.file.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		ps.SetTitleAndData(fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key), prj)
		return controller.Render(rc, as, &vproject.Files{Project: prj}, ps, "projects", prj.Key, "Files")
	})
}

func ProjectFile(rc *fasthttp.RequestCtx) {
	controller.Act("project.file", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		pathS, err := cutil.RCRequiredString(rc, "path", false)
		if err != nil {
			return "", err
		}
		path := util.StringSplitAndTrim(pathS, "/")
		bcAppend := "||/p/" + prj.Key + "/fs"
		bc := []string{"projects", prj.Key, "Files" + bcAppend}
		lo.ForEach(path, func(x string, _ int) {
			bcAppend += "/" + x
			bc = append(bc, x+bcAppend)
		})
		ps.SetTitleAndData(fmt.Sprintf("[%s] /%s", prj.Key, pathS), pathS)
		return controller.Render(rc, as, &vproject.Files{Project: prj, Path: path}, ps, bc...)
	})
}

func ProjectFileStats(rc *fasthttp.RequestCtx) {
	controller.Act("project.file.stats", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		dir := string(rc.URI().QueryArgs().Peek("dir"))
		pth := util.StringSplitAndTrim(dir, "/")
		ext := string(rc.URI().QueryArgs().Peek("ext"))
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		ret, err := stats.GetFileStats(pfs, dir, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] File Stats", prj.Key), ret)
		page := &vproject.FileStats{Project: prj, Path: pth, Ext: ext, Files: ret}
		return controller.Render(rc, as, page, ps, "projects", prj.Key, "File Stats")
	})
}
