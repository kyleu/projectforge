package controller

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/markdown"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vtest"
)

var (
	testbedKeys   = []string{"index", "markdown", "modules", "openapi"}
	testbedTitles = []string{"", "Markdown", "Modules", "OpenAPI"}
)

func Testbed(w http.ResponseWriter, r *http.Request) {
	Act("testbed", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key := ""
		var param any
		var ret any
		frm, _ := cutil.ParseForm(r, ps.RequestBody)
		if len(frm) > 0 {
			key = frm.GetStringOpt("key")
			param = frm.GetMapOpt("param")
		}
		var err error
		switch key {
		case "":
			// noop
		case "modules":
			ret = strings.Join(echoModules(as.Services.Modules.ModulesVisible()), "\n")
		case "openapi":
			ret = openAPI(as, ps)
		case "markdown":
			ret = md(as, ps)
		default:
			return "", errors.Errorf("unhandled key [%s]", key)
		}
		if err != nil {
			return "", errors.Wrapf(err, "unable to run testbed [%s]", key)
		}

		title := "Testbed"
		if idx := slices.Index(testbedKeys, key); idx > -1 {
			title = testbedTitles[idx]
		}
		ps.SetTitleAndData(title, util.ValueMap{"key": key, "param": param, "result": ret})
		return Render(r, as, &vtest.Testbed{Key: key, Param: param, Result: ret, Keys: testbedKeys, Titles: testbedTitles}, ps, "Testbed||/testbed", title)
	})
}

func openAPI(as *app.State, ps *cutil.PageState) any {
	return as.Services.Projects.Projects().Get(util.AppKey)
}

func md(as *app.State, ps *cutil.PageState) any {
	svc := markdown.NewService()
	data := []byte("# Hello!\n\nTest Markdown")
	ret := svc.Parse(data)
	html := svc.HTML(ret, data, 0)
	return html
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
