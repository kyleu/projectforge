package clib

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/scripting"
	"{{{ .Package }}}/views/vscripting"
)

func ScriptingList(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, sizes := as.Services.Script.ListScriptSizes(ps.Logger)
		ps.SetTitleAndData("Scripting", ret)
		return controller.Render(w, r, as, &vscripting.List{Scripts: ret, Sizes: sizes}, ps, "scripting")
	})
}

func ScriptingDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", true)
		if err != nil {
			return "", err
		}
		key, src, err := as.Services.Script.LoadScript(key, ps.Logger)
		if err != nil {
			return "", err
		}
		loadResult, vm, err := scripting.LoadVM(key, src, ps.Logger)
		if err != nil {
			return "", err
		}
		res, err := scripting.RunExamples(vm)
		if err != nil {
			return "", err
		}
		ps.Title = key
		ps.Data = map[string]any{"script": src, "results": res}
		page := &vscripting.Detail{Path: key, Script: src, LoadResult: loadResult, Results: res}
		return controller.Render(w, r, as, page, ps, "scripting", key)
	})
}

func ScriptingNew(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.new", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "New Script"
		return controller.Render(w, r, as, &vscripting.Form{}, ps, "scripting", "New")
	})
}

func ScriptingCreate(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.create", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		pth := frm.GetStringOpt("path")
		if pth == "" {
			return "", errors.New("must provide path")
		}
		content := frm.GetStringOpt("content")
		if content == "" {
			return "", errors.New("content is required")
		}
		err = as.Services.Script.SaveScript(pth, content, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "Scripting created", "/admin/scripting/"+url.QueryEscape(pth), w, ps)
	})
}

func ScriptingForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", true)
		if err != nil {
			return "", err
		}
		key, sc, err := as.Services.Script.LoadScript(key, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Edit ["+key+"]", sc)
		return controller.Render(w, r, as, &vscripting.Form{Path: key, Content: sc}, ps, "scripting", key, "Edit")
	})
}

func ScriptingSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", true)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		content := frm.GetStringOpt("content")
		if content == "" {
			return "", errors.New("content is required")
		}
		err = as.Services.Script.SaveScript(key, content, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "Scripting saved", "/admin/scripting/"+url.QueryEscape(key), w, ps)
	})
}

func ScriptingDelete(w http.ResponseWriter, r *http.Request) {
	controller.Act("scripting.delete", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", true)
		if err != nil {
			return "", err
		}
		err = as.Services.Script.DeleteScript(key, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "Script deleted", "/admin/scripting", w, ps)
	})
}
