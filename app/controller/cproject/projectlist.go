package cproject

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectList(rc *fasthttp.RequestCtx) {
	controller.Act("project.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		execs := as.Services.Exec.Execs
		tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
		if len(tags) > 0 {
			prjs = prjs.WithTags(tags...)
		}
		switch string(rc.QueryArgs().Peek("sort")) {
		case "package":
			slices.SortFunc(prjs, func(l *project.Project, r *project.Project) int {
				return cmp.Compare(l.Package, r.Package)
			})
		case "port":
			slices.SortFunc(prjs, func(l *project.Project, r *project.Project) int {
				return cmp.Compare(l.Port, r.Port)
			})
		}
		ps.SetTitleAndData("All Projects", prjs)
		switch string(rc.QueryArgs().Peek("fmt")) {
		case "ports":
			msgs := lo.Map(prjs, func(p *project.Project, _ int) string {
				return fmt.Sprintf("%s: %d", p.Key, p.Port)
			})
			_, _ = rc.WriteString(strings.Join(msgs, util.StringDefaultLinebreak))
			return "", nil
		case "versions":
			msgs := lo.Map(prjs, func(p *project.Project, _ int) string {
				return fmt.Sprintf("%s: %s", p.Key, p.Version)
			})
			_, _ = rc.WriteString(strings.Join(msgs, util.StringDefaultLinebreak))
			return "", nil
		default:
			return controller.Render(rc, as, &vproject.List{Projects: prjs, Execs: execs, Tags: tags}, ps, "projects")
		}
	})
}
