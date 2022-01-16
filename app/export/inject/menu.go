package inject

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func Menu(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models))
	msg := `Key: %q, Title: %q, Description: %q, Icon: %q, Route: "/%s"`
	for _, m := range args.Models {
		l := fmt.Sprintf(msg, m.Package, m.TitlePlural(), m.Description, m.Icon, m.Package)
		if len(m.GroupedColumns()) == 0 {
			out = append(out, "&menu.Item{"+l+"},")
		} else {
			out = append(out, "&menu.Item{"+l+", Children: menu.Items{")
			for _, g := range m.GroupedColumns() {
				desc := fmt.Sprintf("%s from %s", g.ProperPlural(), m.Plural())
				gl := fmt.Sprintf(msg, g.Camel(), g.ProperPlural(), desc, m.Icon, m.Package+"/"+g.Camel())
				out = append(out, "\t&menu.Item{"+gl+"},")
			}
			out = append(out, "}},")
		}
	}
	content := map[string]string{"codegen": "\n\tret = append(ret,\n\t\t" + strings.Join(out, "\n\t\t") + "\n\t)\n\t// "}
	return file.Inject(f, content)
}
