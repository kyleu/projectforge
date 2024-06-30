package cproject

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func List(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		execs := as.Services.Exec.Execs
		tags := util.StringSplitAndTrim(r.URL.Query().Get("tags"), ",")
		if len(tags) > 0 {
			prjs = prjs.WithTags(tags...)
		}
		switch r.URL.Query().Get("sort") {
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
		switch r.URL.Query().Get("fmt") {
		case "ports":
			msgs := lo.Map(prjs, func(p *project.Project, _ int) string {
				return fmt.Sprintf("%s: %d", p.Key, p.Port)
			})
			_, _ = ps.W.Write([]byte(strings.Join(msgs, util.StringDefaultLinebreak)))
			return "", nil
		case "versions":
			msgs := lo.Map(prjs, func(p *project.Project, _ int) string {
				return fmt.Sprintf("%s: %s", p.Key, p.Version)
			})
			_, _ = ps.W.Write([]byte(strings.Join(msgs, util.StringDefaultLinebreak)))
			return "", nil
		default:
			page := &vproject.List{Projects: prjs, Execs: execs, Tags: tags, Icon: "code"}
			return controller.Render(r, as, page, ps, "projects")
		}
	})
}
