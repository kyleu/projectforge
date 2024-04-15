package cmodule

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/search"
	"projectforge.dev/projectforge/app/lib/search/result"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vmodule"
)

func ModuleSearch(w http.ResponseWriter, r *http.Request) {
	controller.Act("module.search", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(r, as, ps)
		if err != nil {
			return "", err
		}

		q := r.URL.Query().Get("q")
		params := &search.Params{Q: q, PS: nil}

		var res result.Results
		if q != "" {
			fs := as.Services.Modules.GetFilesystem(mod.Key)
			files, err := fs.ListFilesRecursive("", []string{".png$"}, ps.Logger)
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

		ps.SetTitleAndData(fmt.Sprintf("[%s] Module Results", mod.Title()), mod)
		page := &vmodule.Search{Module: mod, Params: params, Results: res}
		return controller.Render(r, as, page, ps, "modules", mod.Key, "Search**archive")
	})
}

func newModuleResult(q string, modKey string, path string, content []byte) *result.Result {
	if len(content) > 1024*1024*4 {
		return nil
	}
	if !utf8.Valid(content) {
		return nil
	}
	_, fn := util.StringSplitPath(path)

	lines := util.StringSplitLines(string(content))

	matches := lo.FlatMap(lines, func(line string, idx int) []*result.Match {
		return result.MatchesFor(fmt.Sprint(idx+1), line, q)
	})
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
