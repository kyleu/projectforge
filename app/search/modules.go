package search

import (
	"context"
	"strings"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/module"
)

func searchModules(ctx context.Context, st *app.State, p *Params) (Results, error) {
	var ret Results
	for _, mod := range st.Services.Modules.Modules() {
		if m := moduleMatches(mod, p.Q); len(m) > 0 {
			res := &Result{ID: mod.Key, Type: "module", Title: mod.Title(), Icon: mod.IconSafe(), URL: "/m/" + mod.Key, Matches: MatchesFrom(m), Data: mod}
			ret = append(ret, res)
		}
	}

	return ret, nil
}

func moduleMatches(mod *module.Module, q string) []string {
	var ret []string
	ql := strings.ToLower(q)
	f := func(k string, v string) {
		if strings.Contains(strings.ToLower(v), ql) {
			ret = append(ret, k+": "+v)
		}
	}
	f("key", mod.Key)
	f("name", mod.Name)
	f("authorName", mod.AuthorName)
	f("authorEmail", mod.AuthorEmail)
	f("description", mod.Description)
	f("license", mod.License)
	f("sourcecode", mod.Sourcecode)
	return ret
}
