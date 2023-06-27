package cproject

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vpage"
)

var changeDirArgs = cutil.Args{{Key: "dir", Title: "Directory", Description: "Filesystem directory to use as the main working directory"}}

func ChangeDir(rc *fasthttp.RequestCtx) {
	controller.Act("change.dir", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.HideMenu = true
		argRes := cutil.CollectArgs(rc, changeDirArgs)
		dir, ok := argRes.Values["dir"]
		if !ok || dir == "" || len(argRes.Missing) > 0 {
			ps.Data = argRes
			d, _ := filepath.Abs(".")
			argRes.Values["dir"] = d
			msg := "Choose the working directory to use for loading the main project"
			return controller.Render(rc, as, &vpage.Args{URL: "/welcome/changedir", Directions: msg, ArgRes: argRes}, ps, "Welcome")
		}

		err := os.Chdir(dir)
		if err != nil {
			err = os.MkdirAll(dir, 0o755)
			if err != nil {
				return "", errors.Wrap(err, "unable to find or create new directory")
			}
		}

		ps.Data = dir
		return "/welcome", nil
	})
}
