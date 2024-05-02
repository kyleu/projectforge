package csfiles

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const (
	ImpMVC         = "Microsoft.AspNetCore.Mvc"
	ImpEF          = "Microsoft.EntityFrameworkCore"
	ImpControllers = "Shared.Controllers"
	ImpServices    = "Shared.Services"
)

func ToCSharpType(col *model.Column) any {
	switch col.Type.Key() {
	case "uuid":
		return "Guid"
	default:
		return fmt.Sprintf("unknown-type[%s]", col.Type.String())
	}
}

func ToCSharpViewType(col *model.Column) any {
	switch col.Type.Key() {
	case "uuid":
		return ":guid"
	default:
		return ""
	}
}

func CSNamespace(prefix string, grp string, path ...string) string {
	ns := util.StringToTitle(grp)
	if ns != "Shared" {
		ns = util.StringToTitle(prefix) + ns
	}
	return strings.Join(append([]string{ns}, path...), ".")
}
