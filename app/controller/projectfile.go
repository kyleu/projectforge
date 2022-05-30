package controller

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/stats"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectFileRoot(rc *fasthttp.RequestCtx) {
	act("project.file.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(rc, as, &vproject.Files{Project: prj}, ps, "projects", prj.Key)
	})
}

func ProjectFile(rc *fasthttp.RequestCtx) {
	act("project.file", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
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
		for _, x := range path {
			bcAppend += "/" + x
			b := x + bcAppend
			bc = append(bc, b)
		}
		ps.Title = fmt.Sprintf("[%s] /%s", prj.Key, strings.Join(path, "/"))
		return render(rc, as, &vproject.Files{Project: prj, Path: path}, ps, bc...)
	})
}

func ProjectFileStats(rc *fasthttp.RequestCtx) {
	act("project.file.stats", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		dir := string(rc.URI().QueryArgs().Peek("dir"))
		pth := util.StringSplitAndTrim(dir, "/")
		ext := string(rc.URI().QueryArgs().Peek("ext"))
		ret, err := stats.GetFileStats(as.Services.Projects.GetFilesystem(prj), dir, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Data = ret
		ps.Title = fmt.Sprintf("[%s] File Stats", prj.Key)
		page := &vproject.FileStats{Project: prj, Path: pth, Ext: ext, Files: ret}
		return render(rc, as, page, ps, "projects", prj.Key, "File Stats")
	})
}
