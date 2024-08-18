package controller

import (
	"fmt"
	"net/http"
	"strings"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/views/vtest"
)

func Testbed(w http.ResponseWriter, r *http.Request) {
	Act("testbed", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := "Testbed!"
		frm, _ := cutil.ParseForm(r, ps.RequestBody)
		if len(frm) > 0 {
			ret = frm.GetStringOpt("x")
		}
		if ret == "modules" {
			ret = strings.Join(echoModules(as.Services.Modules.ModulesVisible()), "\n")
		}
		ps.SetTitleAndData("Testbed", ret)
		return Render(r, as, &vtest.Testbed{Param: ret}, ps)
	})
}

func echoModules(mods module.Modules) []string {
	var ret []string
	log := func(msg string, args ...any) {
		ret = append(ret, fmt.Sprintf(msg, args...))
	}
	log("|   |   |")
	log("|---|---|")
	for _, m := range mods {
		log("|[%s](%s)|%s", m.Name, fmt.Sprintf("module/%s/doc/module/%s", m.Key, m.Key), m.Description)
	}
	return ret
}
