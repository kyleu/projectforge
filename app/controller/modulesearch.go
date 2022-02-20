package controller

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/kyleu/projectforge/app/lib/search"
	"github.com/kyleu/projectforge/app/lib/search/result"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vmodule"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ModuleSearch(rc *fasthttp.RequestCtx) {
	act("module.search", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}

		q := string(rc.URI().QueryArgs().Peek("q"))
		params := &search.Params{
			Q:  q,
			PS: nil,
		}

		var res result.Results
		if q != "" {
			fs := as.Services.Modules.GetFilesystem(mod.Key)
			files, err := fs.ListFilesRecursive("", []string{".png$"})
			if err != nil {
				return "", err
			}

			for _, path := range files {
				if len(res) > 100 {
					continue
				}
				content, err := fs.ReadFile(path)
				if err != nil {
					return "", err
				}
				x := newModuleResult(q, mod.Key, path, content)
				if x != nil {
					res = append(res, x)
				}
			}
		}

		ps.Title = fmt.Sprintf("[%s] Search Results", mod.Title())
		ps.Data = mod
		// page := &vsearch.Results{Params: params, Results: res, SearchPath: fmt.Sprintf("/p/%s/search", mod.Key)}
		page := &vmodule.Search{Module: mod, Params: params, Results: res}
		return render(rc, as, page, ps, "modules", mod.Key, "Search")
	})
}

func newModuleResult(q string, modKey string, path string, content []byte) *result.Result {
	if len(content) > 1024*1024*4 {
		return nil
	}
	if !utf8.Valid(content) {
		return nil
	}

	fn := path
	if strings.Contains(path, "/") {
		_, fn = util.StringSplitLast(path, '/', true)
	}

	lines := strings.Split(string(content), "\n")

	var matches result.Matches
	for idx, line := range lines {
		m := result.MatchesFor(fmt.Sprint(idx+1), line, q)
		matches = append(matches, m...)
	}

	if len(matches) == 0 {
		return nil
	}

	return &result.Result{
		Type:    "file",
		ID:      fn,
		Title:   fn,
		Icon:    "star",
		URL:     fmt.Sprintf("/m/%s/fs/%s", modKey, path),
		Matches: matches,
		Data:    nil,
	}
}
