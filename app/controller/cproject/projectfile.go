package cproject

import (
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/stats"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectFileRoot(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.file.root", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}

		ps.SetTitleAndData(fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key), prj)
		return controller.Render(w, r, as, &vproject.Files{Project: prj}, ps, "projects", prj.Key, "Files")
	})
}

func ProjectFile(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.file", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}

		pathS, err := cutil.RCRequiredString(r, "path", false)
		if err != nil {
			return "", err
		}
		path := util.StringSplitAndTrim(pathS, "/")
		bcAppend := "||" + prj.WebPath() + "/fs"
		bc := []string{"projects", prj.Key, "Files" + bcAppend}
		lo.ForEach(path, func(x string, _ int) {
			bcAppend += "/" + x
			bc = append(bc, x+bcAppend)
		})
		ps.SetTitleAndData(fmt.Sprintf("[%s] /%s", prj.Key, pathS), pathS)
		return controller.Render(w, r, as, &vproject.Files{Project: prj, Path: path}, ps, bc...)
	})
}

func ProjectFileStats(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.file.stats", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		dir := r.URL.Query().Get("dir")
		pth := util.StringSplitAndTrim(dir, "/")
		ext := r.URL.Query().Get("ext")
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
		return controller.Render(w, r, as, page, ps, "projects", prj.Key, "File Stats")
	})
}
