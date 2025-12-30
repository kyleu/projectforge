package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/stats"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func FileRoot(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.file.root", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := GetProject(r, as)
		if err != nil {
			return "", err
		}
		fsys, _ := as.Services.Projects.GetFilesystem(prj)
		ps.SetTitleAndData(fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key), prj)
		return controller.Render(r, as, &vproject.Files{Project: prj, FS: fsys}, ps, "projects", prj.Key, "Files")
	})
}

func File(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.file", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := GetProject(r, as)
		if err != nil {
			return "", err
		}

		pathS, err := cutil.PathString(r, "path", false)
		if err != nil {
			return "", err
		}
		path := util.StringSplitAndTrim(pathS, "/")
		fsys, _ := as.Services.Projects.GetFilesystem(prj)
		if cutil.QueryStringString(ps.URI, "download") == "true" {
			b, err := fsys.ReadFile(util.StringFilePath(path...))
			if err != nil {
				return "", errors.Wrapf(err, "unable to read project file [%s] for download", pathS)
			}
			return cutil.RespondDownload(path[len(path)-1], b, ps.W)
		}

		bcAppend := dblpipe + prj.WebPath() + "/fs"
		bc := []string{"projects", prj.Key, "Files" + bcAppend}
		lo.ForEach(path, func(x string, _ int) {
			bcAppend += "/" + x
			bc = append(bc, x+bcAppend)
		})
		ps.SetTitleAndData(fmt.Sprintf("[%s] /%s", prj.Key, pathS), pathS)
		return controller.Render(r, as, &vproject.Files{Project: prj, Path: path, FS: fsys}, ps, bc...)
	})
}

func FileStats(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.file.stats", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := GetProject(r, as)
		if err != nil {
			return "", err
		}
		dir := cutil.QueryStringString(ps.URI, "dir")
		pth := util.StringSplitAndTrim(dir, "/")
		ext := cutil.QueryStringString(ps.URI, "ext")
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
		return controller.Render(r, as, page, ps, "projects", prj.Key, "File Stats")
	})
}
