package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func (t *TemplateContext) CIContent() string {
	if t.Info == nil {
		return ""
	}
	common := func(x string) string {
		return "on:" + t.Linebreak + "  push:" + t.Linebreak + "    tags:" + t.Linebreak + "      - '" + x + "'"
	}
	switch t.Info.CI {
	case "all":
		return "on: push"
	case "tags":
		return common("*")
	case "versions":
		return common("v*")
	default:
		return common("DISABLED_v*")
	}
}

func (t *TemplateContext) ConfigVarsContent() string {
	ret, err := util.MarkdownTable([]string{"Name", "Type", "Description"}, t.ConfigVars.Array(t.Key), t.Linebreak)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return ret
}

func (t *TemplateContext) ExtraFilesContent() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := util.NewStringSlice([]string{"\n    extra_files:"})
	lo.ForEach(t.Info.ExtraFiles, func(ef string, _ int) {
		ret.Push("      - " + ef)
	})
	return ret.Join(t.Linebreak)
}

func (t *TemplateContext) ExtraFilesDocker() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := lo.Map(t.Info.ExtraFiles, func(ef string, _ int) string {
		return fmt.Sprintf("\nCOPY %s /%s", ef, ef)
	})
	return strings.Join(ret, "")
}

func (t *TemplateContext) GoBinaryContent() string {
	if t.Info == nil || t.Info.GoBinary == "" || t.Info.GoBinary == goStdBin {
		return ""
	}
	return fmt.Sprintf("\n    gobinary: %q", t.Info.GoBinary)
}

func (t *TemplateContext) IgnoredSetting() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	return "\n  exclude-dirs:\n" + strings.Join(lo.Map(t.Ignore, func(i string, _ int) string {
		return "    - /" + strings.TrimPrefix(i, "^")
	}), "\n")
}

func (t *TemplateContext) IgnoredQuoted() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	return strings.Join(lo.Map(t.Ignore, func(i string, _ int) string {
		return fmt.Sprintf(", %q", strings.TrimPrefix(i, "^"))
	}), "")
}

func (t *TemplateContext) TypeScriptProjectContent() string {
	return strings.Join(lo.Map(lo.Filter(t.Info.Deployments, func(x string, _ int) bool {
		return strings.HasPrefix(x, "ts:")
	}), func(x string, _ int) string {
		return strings.Join([]string{
			"",
			"",
			"esbuild.build({",
			`  entryPoints: ["src/game/game.ts"],`,
			"  bundle: true,",
			"  minify: true,",
			"  sourcemap: true,",
			`  outfile: "../assets/game.js",`,
			"  logLevel: \"info\"",
			"}).catch((e) => console.error(e.message));",
		}, "\n")
	}), "\n")
}

func (t *TemplateContext) HasModules(keys ...string) bool {
	return lo.ContainsBy(keys, func(key string) bool {
		return lo.Contains(t.Modules, key)
	})
}

func (t *TemplateContext) HasModule(key string) bool {
	return t.HasModules(key)
}

func (t *TemplateContext) Public() bool {
	return t.Build == nil || !t.Build.Private
}

func (t *TemplateContext) Private() bool {
	return t.Build == nil && t.Build.Private
}

func (t *TemplateContext) Acronyms() string {
	if t.Info == nil || len(t.Info.Acronyms) == 0 {
		return ""
	}
	return strings.Join(util.StringArrayQuoted(t.Info.Acronyms), ", ")
}

func (t *TemplateContext) CoreStruct() string {
	ret := &util.StringSlice{}
	if t.HasModule("") {
		ret.Push("...")
	}
	if t.HasModule("") {
		ret.Push("...")
	}
	if t.HasModule("") {
		ret.Push("...")
	}
	if t.HasModule("") {
		ret.Push("...")
	}
	if t.HasModule("") {
		ret.Push("...")
	}
	if t.HasModule("") {
		ret.Push("...")
	}
	return ret.Join("\n")
}
