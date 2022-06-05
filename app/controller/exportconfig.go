package controller

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportConfigForm(rc *fasthttp.RequestCtx) {
	act("project.export.config.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		args, err := prj.Info.ModuleArgExport()
		if err != nil {
			return "", err
		}

		ps.Data = args

		bc := []string{"projects", prj.Key, "Export"}
		ps.Title = fmt.Sprintf("[%s] Export", prj.Key)
		return render(rc, as, &vexport.ConfigForm{Project: prj, Cfg: args.Config}, ps, bc...)
	})
}

func ProjectExportConfigSave(rc *fasthttp.RequestCtx) {
	act("project.export.config.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		args, err := prj.Info.ModuleArgExport()
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

		// TODO save

		return flashAndRedir(true, "configuration saved", fmt.Sprintf("/p/%s/export", prj.Key), rc, ps)
	})
}
