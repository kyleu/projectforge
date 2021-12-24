package action

import (
	"github.com/kyleu/projectforge/app/schema"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func onCodegen(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)

	start := util.TimerStart()
	fs := pm.PSvc.GetFilesystem(pm.Prj)
	schemaJSON, err := fs.ReadFile("schema.json")
	if err != nil {
		return errorResult(errors.Errorf("project [%s] has no schema defined", pm.Prj.Key), pm.Cfg, pm.Logger)
	}

	sch := &schema.Schema{}
	err = util.FromJSON([]byte(schemaJSON), sch)
	if err != nil {
		return errorResult(errors.Wrapf(err, "project [%s] has an invalid schema", pm.Prj.Key), pm.Cfg, pm.Logger)
	}
	ret.AddLog("schema loaded in [%s]", util.MicrosToMillis(util.TimerEnd(start)))

	if model := pm.Cfg.GetStringOpt("model"); model == "" {
		ret.HTML = append(ret.HTML, "<em>Hoooooray!</em>")
	} else {
		err = codegenModel(pm, sch, ret, model)
		if err != nil {
			return errorResult(errors.Wrapf(err, "unable to generate code for project [%s]", pm.Prj.Key), pm.Cfg, pm.Logger)
		}
	}
	return ret
}

func codegenModel(pm *PrjAndMods, sch *schema.Schema, ret *Result, model string) error {
	ret.AddLog("generating code for project [%s], model [%s]", pm.Prj.Key, model)

	start := util.TimerStart()
	ret.AddLog("processing model [%s]", model)

	ret.Duration = util.TimerEnd(start)
	ret.AddLog("code generated for project [%s] in [%s]", pm.Prj.Key, util.MicrosToMillis(ret.Duration))
	return nil
}
