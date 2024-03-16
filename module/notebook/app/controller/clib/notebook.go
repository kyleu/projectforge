package clib

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vnotebook"
)

func Notebook(rc *fasthttp.RequestCtx) {
	controller.Act("notebook.view", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		status := as.Services.Notebook.Status()
		if status == "running" {
			pathS, _ := cutil.RCRequiredString(rc, "path", false)
			path := util.StringSplitAndTrim(pathS, "/")
			ps.SetTitleAndData("Notebook", path)
			bc := []string{"notebook"}
			if pathS != "" {
				bc = append(bc, pathS)
			}
			return controller.Render(rc, as, &vnotebook.Notebook{BaseURL: as.Services.Notebook.BaseURL, Path: pathS}, ps, bc...)
		}
		ps.SetTitleAndData("Notebook Options", status)
		return controller.Render(rc, as, &vnotebook.Options{}, ps, "notebook")
	})
}

func NotebookFiles(rc *fasthttp.RequestCtx) {
	controller.Act("notebook.files", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pathS, path, bc := notebookGetText(rc)
		fs := as.Services.Notebook.FS
		files, err := fs.ListTree(nil, pathS, []string{"cache"}, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "error listing files")
		}
		ps.SetTitleAndData(fmt.Sprintf("Notebook File /%s", pathS), files)
		return controller.Render(rc, as, &vnotebook.Files{Path: path, FS: as.Services.Notebook.FS}, ps, bc...)
	})
}

func NotebookFileEdit(rc *fasthttp.RequestCtx) {
	controller.Act("notebook.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pathS, path, bc := notebookGetText(rc)
		bc = append(bc, "Edit**edit")
		b, err := as.Services.Notebook.FS.ReadFile(pathS)
		if err != nil {
			return "", err
		}
		ret := string(b)
		ps.SetTitleAndData(fmt.Sprintf("Notebook File /%s", pathS), path)
		return controller.Render(rc, as, &vnotebook.FileEdit{Path: path, Content: ret}, ps, bc...)
	})
}

func NotebookFileSave(rc *fasthttp.RequestCtx) {
	controller.Act("notebook.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pathS, _, _ := notebookGetText(rc)
		msg := "Notebook file saved"

		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		content := frm.GetStringOpt("content")
		if strings.TrimSpace(content) == "" {
			return "", errors.Errorf("file [%s] contents may not be empty", pathS)
		}
		content = strings.ReplaceAll(content, "\r\n", "\n")

		b, err := as.Services.Notebook.FS.ReadFile(pathS)
		if err != nil {
			return "", err
		}
		current := string(b)

		if current == content {
			msg = "No changes required"
		} else {
			err = as.Services.Notebook.FS.WriteFile(pathS, []byte(content), filesystem.DefaultMode, true)
			if err != nil {
				return "", err
			}
		}
		return controller.FlashAndRedir(true, msg, "/notebook/files/"+pathS, rc, ps)
	})
}

func NotebookAction(rc *fasthttp.RequestCtx) {
	controller.Act("notebook.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		act, err := cutil.RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}
		switch act {
		case "start":
			err = as.Services.Notebook.Start()
			return controller.FlashAndRedir(true, "Notebook started", "/notebook", rc, ps)
		default:
			return "", errors.Errorf("invalid notebook action [%s]", act)
		}
	})
}

func notebookGetText(rc *fasthttp.RequestCtx) (string, []string, []string) {
	pathS, _ := cutil.RCRequiredString(rc, "path", false)
	path := util.StringSplitAndTrim(pathS, "/")
	bcAppend := "||/notebook/files"
	bc := []string{"notebook", "files"}
	lo.ForEach(path, func(x string, _ int) {
		bcAppend += "/" + x
		b := x + bcAppend
		bc = append(bc, b)
	})
	return pathS, path, bc
}
