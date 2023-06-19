package cproject

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/verror"
)

var changeDirArgs = cutil.Args{{Key: "dir", Title: "Directory", Description: "Filesystem directory to use as the main working directory"}}

func ChangeDir(rc *fasthttp.RequestCtx) {
	controller.Act("change.dir", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.HideMenu = true
		argRes := cutil.CollectArgs(rc, changeDirArgs)
		dir, ok := argRes.Values["dir"]
		if !ok || len(dir) == 0 || len(argRes.Missing) > 0 {
			ps.Data = argRes
			dir, _ := filepath.Abs(".")
			argRes.Values["dir"] = dir
			msg := "Choose the working directory to use for loading the main project"
			return controller.Render(rc, as, &verror.Args{URL: "/welcome/changedir", Directions: msg, ArgRes: argRes}, ps, "Welcome")
		}

		err := os.Chdir(dir)
		if err != nil {
			return "", errors.Wrap(err, "unable to set new directory")
		}

		ps.Data = dir
		return "/welcome", nil
	})
}
