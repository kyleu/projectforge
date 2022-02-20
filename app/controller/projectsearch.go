package controller

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/kyleu/projectforge/app/lib/search"
	"github.com/kyleu/projectforge/app/lib/search/result"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vproject"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ProjectSearch(rc *fasthttp.RequestCtx) {
	act("project.search", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		q := string(rc.URI().QueryArgs().Peek("q"))
		params := &search.Params{
			Q:  q,
			PS: nil,
		}

		res, err := searchProject(prj, q, as)
		if err != nil {
			return "", errors.Wrapf(err, "unable to search project [%s]", prj.Key)
		}

		ps.Title = fmt.Sprintf("[%s] Search Results", prj.Title())
		ps.Data = res
		// page := &vsearch.Results{Params: params, Results: res, SearchPath: fmt.Sprintf("/p/%s/search", prj.Key)}
		page := &vproject.Search{Project: prj, Params: params, Results: res}
		return render(rc, as, page, ps, "projects", prj.Key, "Search")
	})
}

func ProjectSearchAll(rc *fasthttp.RequestCtx) {
	act("project.search.all", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()

		q := string(rc.URI().QueryArgs().Peek("q"))
		params := &search.Params{
			Q:  q,
			PS: nil,
		}

		ret := map[string]result.Results{}

		for _, prj := range prjs {
			res, err := searchProject(prj, q, as)
			if err != nil {
				return "", errors.Wrapf(err, "unable to search project [%s]", prj.Key)
			}

			ret[prj.Key] = res
		}
		ps.Title = "Search Results"
		ps.Data = ret
		page := &vproject.SearchAll{Params: params, Projects: prjs, Results: ret}
		return render(rc, as, page, ps, "projects", "Search")
	})
}

func searchProject(prj *project.Project, q string, as *app.State) (result.Results, error) {
	if q == "" {
		return nil, nil
	}
	var res result.Results
	fs := as.Services.Projects.GetFilesystem(prj)
	files, err := fs.ListFilesRecursive("", append([]string{".png$"}, prj.Ignore...))
	if err != nil {
		return nil, err
	}

	for _, path := range files {
		if len(res) > 100 {
			continue
		}
		content, err := fs.ReadFile(path)
		if err != nil {
			return nil, err
		}
		x := newProjectResult(q, prj.Key, path, content)
		if x != nil {
			res = append(res, x)
		}
	}
	return res, nil
}

func newProjectResult(q string, prjKey string, path string, content []byte) *result.Result {
	if len(content) > 1024*1024*4 {
		return nil
	}
	if !utf8.Valid(content) {
		return nil
	}

	_, fn := util.StringSplitLast(path, '/', true)

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
		URL:     fmt.Sprintf("/p/%s/fs/%s", prjKey, path),
		Matches: matches,
		Data:    nil,
	}
}
