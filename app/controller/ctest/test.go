package ctest

import (
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vtest"
)

func TestList(w http.ResponseWriter, r *http.Request) {
	controller.Act("test.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("Tests", []string{"bootstrap", "diff"})
		return controller.Render(r, as, &vtest.List{}, ps, "Tests")
	})
}

func TestRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("test.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Test ["+key+"]", key)
		bc := []string{"Tests||/test"}

		var page layout.Page
		switch key {
		case "diff":
			bc = append(bc, "Diff")
			page, err = diffTest(ps)
		case "bootstrap":
			bc = append(bc, "Bootstrap")
			page, err = bootstrapTest(as, ps)
		case "search":
			bc = append(bc, "Search")
			page, err = searchTest(as, r, ps)
		case "jsonschema-init":
			bc = append(bc, "JSON Schema Init")
			page, err = schemaInit(ps)
		case "jsonschema-test":
			bc = append(bc, "JSON Schema")
			page, err = schemaTest(ps)
		default:
			return "", errors.Errorf("invalid test [%s]", key)
		}
		if err != nil {
			return "", err
		}
		return controller.Render(r, as, page, ps, bc...)
	})
}
