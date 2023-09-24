package project

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
	switch t.Info.CI {
	case "all":
		return "on: push"
	case "tags":
		return "on:" + t.Linebreak + "  push:" + t.Linebreak + "    tags:" + t.Linebreak + "      - '*'"
	case "versions":
		return "on:" + t.Linebreak + "  push:" + t.Linebreak + "    tags:" + t.Linebreak + "      - 'v*'"
	default:
		return "on:" + t.Linebreak + "  push:" + t.Linebreak + "    tags:" + t.Linebreak + "      - 'DISABLED_v*'"
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
	ret := []string{"\n    extra_files:"}
	lo.ForEach(t.Info.ExtraFiles, func(ef string, _ int) {
		ret = append(ret, "      - "+ef)
	})
	return strings.Join(ret, t.Linebreak)
}

func (t *TemplateContext) ExtraFilesDocker() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := make([]string, 0, len(t.Info.ExtraFiles))
	lo.ForEach(t.Info.ExtraFiles, func(ef string, _ int) {
		ret = append(ret, fmt.Sprintf("\nCOPY %s /%s", ef, ef))
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
	ret := make([]string, 0, len(t.Ignore))
	lo.ForEach(t.Ignore, func(i string, _ int) {
		ret = append(ret, "/"+strings.TrimPrefix(i, "^"))
	})
	return " --skip-dirs \"" + strings.Join(ret, "|") + "\""
}

func (t *TemplateContext) IgnoredQuoted() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	ret := make([]string, 0, len(t.Ignore))
	lo.ForEach(t.Ignore, func(i string, _ int) {
		ret = append(ret, fmt.Sprintf(", %q", strings.TrimPrefix(i, "^")))
	})
	return strings.Join(ret, "")
}

func (t *TemplateContext) ExplainPrefix() string {
	if t.HasModules(util.DatabasePostgreSQL, util.DatabaseSQLite) {
		return "explain query plan "
	}
	return "explain "
}

func (t *TemplateContext) HasModules(keys ...string) bool {
	return lo.ContainsBy(keys, func(key string) bool {
		return lo.Contains(t.Modules, key)
	})
}

func (t *TemplateContext) HasModule(key string) bool {
	return t.HasModules(key)
}
