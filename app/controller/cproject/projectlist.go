package cproject

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectList(rc *fasthttp.RequestCtx) {
	controller.Act("project.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		execs := as.Services.Exec.Execs
		tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
		if len(tags) > 0 {
			prjs = prjs.WithTags(tags...)
		}
		ps.Title = "Project Listing"

		switch string(rc.QueryArgs().Peek("sort")) {
		case "package":
			slices.SortFunc(prjs, func(l *project.Project, r *project.Project) bool {
				return l.Package < r.Package
			})
		case "port":
			slices.SortFunc(prjs, func(l *project.Project, r *project.Project) bool {
				return l.Port < r.Port
			})
		}

		ps.Title = "All Projects"
		ps.Data = prjs
		switch string(rc.QueryArgs().Peek("fmt")) {
		case "ports":
			msgs := lo.Map(prjs, func(p *project.Project, _ int) string {
				return fmt.Sprintf("%s: %d", p.Key, p.Port)
			})
			_, _ = rc.WriteString(strings.Join(msgs, "\n"))
			return "", nil
		case "versions":
			msgs := lo.Map(prjs, func(p *project.Project, _ int) string {
				return fmt.Sprintf("%s: %s", p.Key, p.Version)
			})
			_, _ = rc.WriteString(strings.Join(msgs, "\n"))
			return "", nil
		case "go":
			_, _ = rc.WriteString(strings.Join(mkGoSvcs(prjs), "\n"))
			return "", nil
		default:
			return controller.Render(rc, as, &vproject.List{Projects: prjs, Execs: execs, Tags: tags}, ps, "projects")
		}
	})
}

func mkGoSvcs(prjs project.Projects) []string {
	ret := []string{"package library", "", "var ("}
	w := func(s string, args ...any) {
		ret = append(ret, fmt.Sprintf("  "+s, args...))
	}
	lo.ForEach(prjs, func(p *project.Project, _ int) {
		tags := make([]string, 0, len(p.Tags)+2)
		tags = append(tags, "go", "v2")
		tags = append(tags, p.Tags...)
		if p.HasModule("database") {
			tags = append(tags, "database")
		}
		if p.HasModule("grpc") {
			tags = append(tags, "grpc")
		}
		if p.HasModule("temporal") {
			tags = append(tags, "temporal")
		}
		w("%s = &svc.Svc{", p.NameSafe())
		w("	Key:         %q,", p.Key)
		w("	Name:        %q,", p.Name)
		w("	Description: %q,", p.DescriptionSafe())
		w("	Repo:        %q,", p.Info.Sourcecode)
		w("	Icon:        %q,", p.Icon)
		w("	ColorLight:  %q,", p.Theme.Light.NavBackground)
		w("	ColorDark:   %q,", p.Theme.Dark.NavBackground)
		w("	Owners:      []string{%q},", p.Info.AuthorIDSafe())
		w("	Tags:        []string{%s},", strings.Join(util.StringArrayQuoted(tags), ", "))
		if p.HasModule("grpc") {
			w("	Ports:       map[string]int{%q: %d, %q: %d},", "http", p.Port, "grpc", p.Port+10)
		} else {
			w("	Ports:       map[string]int{%q: %d},", "http", p.Port)
		}
		if len(p.Info.EnvVars) > 0 {
			w("	EnvVars:     toEnvMap([]string{")
			lo.ForEach(p.Info.EnvVars, func(x string, _ int) {
				w("		%q,", x)
			})
			w("	}),")
		}
		w("}")
		w("")
	})
	w("goServices = svc.Svcs{")
	lo.ForEach(prjs, func(p *project.Project, _ int) {
		w("  %s,", p.NameSafe())
	})
	w("}")
	ret = append(ret, ")")
	return ret
}
