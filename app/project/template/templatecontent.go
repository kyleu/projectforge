package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func (t *Context) CIContent() string {
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

func (t *Context) ConfigVarsContent() string {
	ret, err := util.MarkdownTable([]string{"Name", "Type", "Description"}, t.ConfigVars.Array(t.Key), t.Linebreak)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return ret
}

func (t *Context) ExtraFilesContent() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := util.NewStringSlice([]string{"\n    extra_files:"})
	lo.ForEach(t.Info.ExtraFiles, func(ef string, _ int) {
		ret.Push("      - " + ef)
	})
	return ret.Join(t.Linebreak)
}

func (t *Context) ExtraFilesDocker() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := lo.Map(t.Info.ExtraFiles, func(ef string, _ int) string {
		return fmt.Sprintf("\nCOPY %s /%s", ef, ef)
	})
	return util.StringJoin(ret, "")
}

func (t *Context) GoBinaryContent() string {
	if t.Info == nil || t.Info.GoBinary == "" || t.Info.GoBinary == goStdBin {
		return ""
	}
	return fmt.Sprintf("\n    gobinary: %q", t.Info.GoBinary)
}

func (t *Context) IgnoredSetting() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	return "\n" + util.StringJoin(lo.Map(t.Ignore, func(i string, _ int) string {
		return "      - /" + strings.TrimPrefix(i, "^")
	}), "\n")
}

func (t *Context) IgnoredQuoted() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	return util.StringJoin(lo.Map(t.Ignore, func(i string, _ int) string {
		return fmt.Sprintf(", %q", strings.TrimPrefix(i, "^"))
	}), "")
}

func (t *Context) HasModules(keys ...string) bool {
	return lo.ContainsBy(keys, func(key string) bool {
		return lo.Contains(t.Modules, key)
	})
}

func (t *Context) HasModule(key string) bool {
	return t.HasModules(key)
}

func (t *Context) Public() bool {
	return t.Build == nil || !t.Build.Private
}

func (t *Context) Private() bool {
	return t.Build == nil && t.Build.Private
}

func (t *Context) DockerPackages() string {
	if t.Info == nil || len(t.Info.DockerPackages) == 0 {
		return ""
	}
	return " " + util.StringJoin(t.Info.DockerPackages, " ")
}

func (t *Context) Acronyms() string {
	if t.Info == nil || len(t.Info.Acronyms) == 0 {
		return ""
	}
	return util.StringJoin(util.StringArrayQuoted(t.Info.Acronyms), ", ")
}

func (t *Context) CoreStruct() string {
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
