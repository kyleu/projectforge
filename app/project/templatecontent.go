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
		return "on:\n  push:\n    tags:\n      - '*'"
	case "versions":
		return "on:\n  push:\n    tags:\n      - 'v*'"
	default:
		return "on:\n  push:\n    tags:\n      - 'DISABLED_v*'"
	}
}

func (t *TemplateContext) ConfigVarsContent() string {
	ret, err := util.MarkdownTable([]string{"Name", "Type", "Description"}, t.ConfigVars.Array(t.Key))
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
	return strings.Join(ret, "\n")
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
	if lo.Contains(t.Modules, "sqlite") && !lo.Contains(t.Modules, "postgres") {
		return "explain query plan "
	}
	return "explain "
}
