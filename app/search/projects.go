package search

import (
	"context"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/project"
)

func searchProjects(ctx context.Context, st *app.State, p *Params) (Results, error) {
	var ret Results
	for _, prj := range st.Services.Projects.Projects() {
		if projectMatches(prj, p.Q) {
			ret = append(ret, &Result{
				ID:      "x",
				Type:    "x",
				Title:   "x",
				Icon:    "x",
				URL:     "x",
				Matches: nil,
				Data:    nil,
			})
		}
	}

	return ret, nil
}

func projectMatches(prj *project.Project, q string) bool {
	return false
}

