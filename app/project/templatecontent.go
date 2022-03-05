package project

import (
	"fmt"
	"strings"
)

func (t *TemplateContext) CIContent() string {
	if t.Info == nil {
		return ""
	}
	switch t.Info.CI {
	case "all":
		return "on: push"
	case "tags":
		return "on:\n  push:\n    tags"
	case "versions":
		return "on:\n  push:\n    tags:\n      - 'v*'"
	default:
		return "on:\n  push:\n    tags:\n      - 'DISABLED_v*'"
	}
}

func (t *TemplateContext) ExtraFilesContent() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := []string{"\n    extra_files:"}
	for _, ef := range t.Info.ExtraFiles {
		ret = append(ret, "      - "+ef)
	}
	return strings.Join(ret, "\n")
}

func (t *TemplateContext) ExtraFilesDocker() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := make([]string, 0, len(t.Info.ExtraFiles))
	for _, ef := range t.Info.ExtraFiles {
		ret = append(ret, "\nCOPY "+ef+" /")
	}
	return strings.Join(ret, "")
}

func (t *TemplateContext) GoBinaryContent() string {
	if t.Info == nil || t.Info.GoBinary == "" || t.Info.GoBinary == "go" {
		return ""
	}
	return fmt.Sprintf("\n    gobinary: %q", t.Info.GoBinary)
}