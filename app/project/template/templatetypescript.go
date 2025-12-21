package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func (t *Context) TypeScriptDeployments() []string {
	return lo.Map(lo.Filter(t.Info.Deployments, func(x string, _ int) bool {
		return strings.HasPrefix(x, "ts:")
	}), func(x string, _ int) string {
		return strings.TrimPrefix(x, "ts:")
	})
}

func (t *Context) TypeScriptProjectWarning() string {
	tsDeps := t.TypeScriptDeployments()
	if len(tsDeps) == 0 {
		return ""
	}
	return "// ESBuild watch doesn't work with multiple projects, use `watchman` or another alternative.\n  "
}

func (t *Context) TypeScriptProjectContent() string {
	tsDeps := t.TypeScriptDeployments()
	if len(tsDeps) == 0 {
		return ""
	}
	convert := func(x string, _ int) string {
		key := x
		if idx := strings.LastIndex(x, "/"); idx != -1 {
			key = x[idx+1:]
		}
		ss := &util.StringSlice{}
		ss.Pushf(`  await esbuild.build({...options, entryPoints: ["src/%s/%s.ts"], outfile: "../assets/%s.js"});`, x, key, key)
		return ss.String()
	}
	return "\n" + util.StringJoin(lo.Map(tsDeps, convert), "\n")
}

func (t *Context) TypeScriptProjectContentStatic() string {
	tsDeps := t.TypeScriptDeployments()
	if len(tsDeps) == 0 {
		return ""
	}
	convert := func(x string, _ int) string {
		key := x
		if idx := strings.LastIndex(x, "/"); idx != -1 {
			key = x[idx+1:]
		}
		ss := &util.StringSlice{}
		ss.Push("")
		ss.Push("esbuild.build({")
		ss.Pushf(`  entryPoints: ["src/%s/%s.ts"],`, x, key)
		ss.Push("  bundle: true,")
		ss.Push("  minify: true,")
		ss.Push("  sourcemap: true,")
		ss.Pushf(`  outfile: "../assets/%s.js",`, key)
		ss.Push(`  logLevel: "info"`)
		ss.Push("});")
		return ss.String()
	}
	return "\n" + util.StringJoin(lo.Map(tsDeps, convert), "\n")
}

func (t *Context) NPMDependencies() string {
	if t.Info == nil || len(t.Info.Dependencies) == 0 {
		return ""
	}
	ss := util.NewStringSlice("", `  "dependencies": {`)
	keys := lo.Keys(t.Info.Dependencies)
	lo.ForEach(lo.Filter(keys, func(x string, _ int) bool {
		return strings.HasPrefix(x, "npm:")
	}), func(x string, idx int) {
		v := t.Info.Dependencies[x]
		x = strings.TrimPrefix(x, "npm:")
		ss.Pushf("    %q: %q%s", x, v, util.Choose(idx == len(keys)-1, "", ","))
	})
	ss.Push("  },")
	return ss.String()
}

func (t *Context) TypeScriptPaths() string {
	paths := map[string]string{}
	if t.HasModule("numeric") {
		paths["@numeric/*"] = "numeric/*"
	}
	for _, pth := range t.TypeScriptDeployments() {
		key := pth
		if idx := strings.LastIndex(pth, "/"); idx != -1 {
			key = pth[idx+1:]
		}
		paths[fmt.Sprintf("@%s/*", key)] = fmt.Sprintf("%s/*", pth)
	}
	if len(paths) == 0 {
		return ""
	}
	ret := util.NewStringSlice("", `    "paths": {`)
	for _, k := range util.ArraySorted(lo.Keys(paths)) {
		ret.Pushf("      %q: [%q]%s", k, paths[k], util.Choose(len(ret.Slice) == len(paths)+1, "", ","))
	}
	ret.Push("    },")
	return ret.String()
}
