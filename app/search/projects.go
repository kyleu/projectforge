package search

import (
	"context"
	"strings"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/project"
)

func searchProjects(ctx context.Context, st *app.State, p *Params) (Results, error) {
	var ret Results
	for _, prj := range st.Services.Projects.Projects() {
		if m := projectMatches(prj, p.Q); len(m) > 0 {
			res := &Result{ID: prj.Key, Type: "project", Title: prj.Title(), Icon: prj.SafeIcon(), URL: "/p/" + prj.Key, Matches: MatchesFrom(m), Data: prj}
			ret = append(ret, res)
		}
	}

	return ret, nil
}

func projectMatches(prj *project.Project, q string) []string {
	var ret []string
	ql := strings.ToLower(q)
	f := func(k string, v string) {
		if strings.Contains(strings.ToLower(v), ql) {
			ret = append(ret, k+": "+v)
		}
	}
	f("key", prj.Key)
	f("name", prj.Name)
	f("version", prj.Version)
	f("package", prj.Package)
	if prj.Info != nil {
		f("authorID", prj.Info.AuthorID)
		f("authorName", prj.Info.AuthorName)
		f("authorEmail", prj.Info.AuthorEmail)
		f("license", prj.Info.License)
		f("sourcecode", prj.Info.Sourcecode)
		f("description", prj.Info.Description)
		f("summary", prj.Info.Summary)
		f("org", prj.Info.Org)
	}
	return ret
}
