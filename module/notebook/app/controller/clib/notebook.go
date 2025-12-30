package clib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vnotebook"
)

func Notebook(w http.ResponseWriter, r *http.Request) {
	controller.Act("notebook.view", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		status := as.Services.Notebook.Status(ps.Context)
		if status == "running" {
			pathS, _ := cutil.PathString(r, "path", false)
			path := util.StringSplitAndTrim(pathS, "/")
			ps.SetTitleAndData("Notebook", path)
			bc := []string{"notebook"}
			if pathS != "" {
				bc = append(bc, pathS)
			}
			return controller.Render(r, as, &vnotebook.Notebook{BaseURL: as.Services.Notebook.BaseURL, Path: pathS}, ps, bc...)
		}
		ps.SetTitleAndData("Notebook Options", status)
		return controller.Render(r, as, &vnotebook.Options{}, ps, "notebook")
	})
}

func NotebookFiles(w http.ResponseWriter, r *http.Request) {
	controller.Act("notebook.files", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		pathS, path, bc := notebookGetText(r)
		fsys := as.Services.Notebook.FS
		if cutil.QueryStringString(ps.URI, "download") == "true" {
			b, err := fsys.ReadFile(util.StringFilePath(path...))
			if err != nil {
				return "", errors.Wrapf(err, "unable to read file [%s] for download", pathS)
			}
			return cutil.RespondDownload(path[len(path)-1], b, ps.W)
		}
		files, err := fsys.ListTree(nil, pathS, []string{"cache"}, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "error listing files")
		}
		ps.SetTitleAndData(fmt.Sprintf("Notebook File /%s", pathS), files)
		return controller.Render(r, as, &vnotebook.Files{Path: path, FS: as.Services.Notebook.FS}, ps, bc...)
	})
}

func NotebookFileEdit(w http.ResponseWriter, r *http.Request) {
	controller.Act("notebook.edit", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		pathS, path, bc := notebookGetText(r)
		bc = append(bc, "Edit**edit")
		b, err := as.Services.Notebook.FS.ReadFile(pathS)
		if err != nil {
			return "", err
		}
		ret := string(b)
		ps.SetTitleAndData(fmt.Sprintf("Notebook File /%s", pathS), path)
		return controller.Render(r, as, &vnotebook.FileEdit{Path: path, Content: ret}, ps, bc...)
	})
}

func NotebookFileSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("notebook.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		pathS, _, _ := notebookGetText(r)
		msg := "Notebook file saved"

		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		return controller.FlashAndRedir(true, msg, "/notebook/files/"+pathS, ps)
	})
}

func NotebookAction(w http.ResponseWriter, r *http.Request) {
	controller.Act("notebook.action", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		act, err := cutil.PathString(r, "act", false)
		if err != nil {
			return "", err
		}
		switch act {
		case util.KeyStart:
			err = as.Services.Notebook.Start(ps.Context)
			if err != nil {
				return "", err
			}
			return controller.FlashAndRedir(true, "Notebook started", "/notebook", ps)
		default:
			return "", errors.Errorf("invalid notebook action [%s]", act)
		}
	})
}

func notebookGetText(r *http.Request) (string, []string, []string) {
	pathS, _ := cutil.PathString(r, "path", false)
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
