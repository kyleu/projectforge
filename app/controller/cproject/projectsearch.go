package cproject

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/search"
	"projectforge.dev/projectforge/app/lib/search/result"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectSearch(rc *fasthttp.RequestCtx) {
	controller.Act("project.search", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		q := string(rc.URI().QueryArgs().Peek("q"))
		params := &search.Params{
			Q:  q,
			PS: nil,
		}

		res, err := searchProject(prj, q, as, ps.Logger)
		if err != nil {
			return "", errors.Wrapf(err, "unable to search project [%s]", prj.Key)
		}

		ps.Title = fmt.Sprintf("[%s] Project Results", prj.Title())
		ps.Data = res
		// page := &vsearch.Results{Params: params, Results: res, SearchPath: fmt.Sprintf("/p/%s/search", prj.Key)}
		page := &vproject.Search{Project: prj, Params: params, Results: res}
		return controller.Render(rc, as, page, ps, "projects", prj.Key, "Search")
	})
}

func ProjectSearchAll(rc *fasthttp.RequestCtx) {
	controller.Act("project.search.all", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
		if len(tags) > 0 {
			prjs = prjs.WithTags(tags...)
		}
		q := string(rc.URI().QueryArgs().Peek("q"))
		params := &search.Params{
			Q:  q,
			PS: nil,
		}

		ret := map[string]result.Results{}

		for _, prj := range prjs {
			res, err := searchProject(prj, q, as, ps.Logger)
			if err != nil {
				return "", errors.Wrapf(err, "unable to search project [%s]", prj.Key)
			}

			ret[prj.Key] = res
		}
		ps.Title = "Project Search Results"
		ps.Data = ret
		page := &vproject.SearchAll{Params: params, Projects: prjs, Tags: tags, Results: ret}
		return controller.Render(rc, as, page, ps, "projects", "Search")
	})
}

func searchProject(prj *project.Project, q string, as *app.State, logger util.Logger) (result.Results, error) {
	if q == "" {
		return nil, nil
	}
	var res result.Results
	fs := as.Services.Projects.GetFilesystem(prj)
	files, err := fs.ListFilesRecursive("", append([]string{".png$"}, prj.Ignore...), logger)
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

	u := fmt.Sprintf("/p/%s/fs/%s", prjKey, path)
	return &result.Result{Type: "file", ID: fn, Title: fn, Icon: "star", URL: u, Matches: matches, Data: nil}
}
