package cproject

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vpage"
)

var changeDirArgs = util.FieldDescs{{Key: "dir", Title: "Directory", Description: "Filesystem directory to use as the main working directory"}}

func ChangeDir(w http.ResponseWriter, r *http.Request) {
	controller.Act("change.dir", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		if runtime.GOOS == "js" {
			return controller.FlashAndRedir(true, "Change directory not available on WASM", "/welcome?override=true", ps)
		}
		ps.HideMenu = true
		argRes := util.FieldDescsCollect(r, changeDirArgs)
		dir, err := argRes.Values.GetString("dir", false)
		if err != nil || dir == "" || argRes.HasMissing() {
			ps.SetTitleAndData("Directory Change", argRes)
			d, _ := filepath.Abs(".")
			argRes.Values["dir"] = d
			msg := "Choose the working directory to use for loading the main project"
			return controller.Render(r, as, &vpage.Args{URL: "/welcome/changedir", Directions: msg, Results: argRes}, ps, "Welcome")
		}

		err = os.Chdir(dir)
		if err != nil {
			fs, err := filesystem.NewFileSystem(dir, false, "")
			if err != nil {
				return "", errors.Wrapf(err, "unable to create filesystem for new directory [%s]", dir)
			}
			err = fs.CreateDirectory(dir)
			if err != nil {
				return "", errors.Wrapf(err, "unable to find or create new directory [%s]", dir)
			}
			err = os.Chdir(dir)
			if err != nil {
				return "", errors.Wrapf(err, "unable to change to new directory [%s]", dir)
			}
		}

		ps.SetTitleAndData("Directory Change", dir)
		return "/welcome", nil
	})
}
