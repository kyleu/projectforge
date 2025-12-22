package clib

import (
	"fmt"
	"net/http"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/search"
	"{{{ .Package }}}/views/vsearch"
)

const searchKey = "search"

func Search(w http.ResponseWriter, r *http.Request) {
	controller.Act(searchKey, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		q := r.URL.Query().Get("q")
		params := &search.Params{Q: q, PS: ps.Params}
		results, errs := search.Search(ps.Context, as, params, ps)
		ps.SetTitleAndData("Search Results", results)
		if q != "" {
			ps.Title = fmt.Sprintf("[%s] %s", q, ps.Title)
		}
		if len(results) == 1 && results[0].URL != "" {
			return controller.FlashAndRedir(true, "single search result found", results[0].URL, ps)
		}
		ps.DefaultNavIcon = searchKey
		bc := []string{"Search||/search"}
		if q != "" {
			bc = append(bc, q)
		}
		return controller.Render(r, as, &vsearch.Results{Params: params, Results: results, Errors: errs}, ps, bc...)
	})
}
