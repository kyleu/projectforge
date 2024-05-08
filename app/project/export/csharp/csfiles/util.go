package csfiles

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
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

func CSNamespace(prefix string, dflt string, path ...string) string {
	if dflt != "Shared" {
		dflt = util.StringToCamel(prefix) + dflt
	}
	return strings.Join(append([]string{dflt}, path...), ".")
}
