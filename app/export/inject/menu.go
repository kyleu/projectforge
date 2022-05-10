package inject

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

const msg = `Key: %q, Title: %q, Description: %q, Icon: %q, Route: "/%s"`

func Menu(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models))

	var adminModels, nonAdminModels model.Models
	for _, m := range args.Models {
		if m.HasTag("public") {
			nonAdminModels = append(nonAdminModels, m)
		} else {
			adminModels = append(adminModels, m)
		}
	}

	if len(nonAdminModels) > 0 {
		out = append(out, "ret = append(ret,")
	}
	for _, m := range nonAdminModels {
		out = append(out, menuFor(m, "\t")...)
	}
	if len(nonAdminModels) > 0 {
		out = append(out, ")")
	}

	if len(adminModels) > 0 {
		out = append(out, "if isAdmin {", "\tret = append(ret,")
	}
	for _, m := range adminModels {
		out = append(out, menuFor(m, "\t\t")...)
	}
	if len(adminModels) > 0 {
		out = append(out, "\t)", "}")
	}

	content := map[string]string{"codegen": "\n\t" + strings.Join(out, "\n\t") + "\n\t// "}
	return file.Inject(f, content)
}

func menuFor(m *model.Model, prefix string) []string {
	var out []string
	l := fmt.Sprintf(msg, m.Package, m.TitlePlural(), m.Description, m.Icon, m.Route())
	if len(m.GroupedColumns()) == 0 {
		out = append(out, prefix+"&menu.Item{"+l+"},")
	} else {
		out = append(out, prefix+"&menu.Item{"+l+", Children: menu.Items{")
		for _, g := range m.GroupedColumns() {
			desc := fmt.Sprintf("%s from %s", g.ProperPlural(), m.Plural())
			gl := fmt.Sprintf(msg, g.Camel(), g.ProperPlural(), desc, m.Icon, m.Route()+"/"+g.Camel())
			out = append(out, prefix+"\t&menu.Item{"+gl+"},")
		}
		out = append(out, prefix+"}},")
	}
	return out
}
