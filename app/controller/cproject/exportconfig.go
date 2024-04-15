package cproject

import (
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportConfigForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.config.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		args, err := prj.ModuleArgExport(as.Services.Projects, ps.Logger)
		if err != nil {
			return "", err
		}

		ps.Data = args

		bc := []string{"projects", prj.Key, "Export"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Export", prj.Key), args.Config)
		return controller.Render(r, as, &vexport.ConfigForm{Project: prj, Cfg: args.Config}, ps, bc...)
	})
}

func ProjectExportConfigSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.config.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		args, err := prj.ModuleArgExport(as.Services.Projects, ps.Logger)
		if err != nil {
			return "", err
		}

		cfgJSON := frm.GetStringOpt("cfg")
		cfg := util.ValueMap{}
		err = util.FromJSON([]byte(cfgJSON), &cfg)
		if err != nil {
			return "", err
		}

		args.Config = cfg

		err = as.Services.Projects.Save(prj, ps.Logger)
		if err != nil {
			return "", err
		}

		return controller.FlashAndRedir(true, "configuration saved", fmt.Sprintf("/p/%s/export", prj.Key), ps)
	})
}
