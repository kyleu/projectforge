package export

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
)

func injectRoutes(f *file.File, args *Args) error {
	var out []string
	for _, m := range args.Models {
		l := fmt.Sprintf("r.GET(\"/%s\", %sList)", m.Package, m.packageProper())
		pkNames := make([]string, 0, len(m.Columns.PKs()))
		for _, col := range m.Columns.PKs() {
			pkNames = append(pkNames, fmt.Sprintf("{%s}", col.Name))
		}
		d := fmt.Sprintf("r.GET(\"/%s/%s\", %sDetail)", m.Package, strings.Join(pkNames, "/"), m.packageProper())
		out = append(out, l, d)
	}
	content := map[string]string{"codegen": "\n\t" + strings.Join(out, "\n\t") + "\n\t// "}
	return file.Inject(f, content)
}

func injectMenu(f *file.File, args *Args) error {
	if len(args.Models) == 0 {
		return nil
	}
	var out []string
	msg := `&menu.Item{Key: %q, Title: %q, Description: %q, Icon: %q, Route: "/%s"},`
	for _, m := range args.Models {
		l := fmt.Sprintf(msg, m.Package, m.proper(), m.Description, m.Icon, m.Package)
		out = append(out, l)
	}
	content := map[string]string{"codegen": "\n\tret = append(ret,\n\t\t" + strings.Join(out, "\n\t\t") + "\n\t)\n\t// "}
	return file.Inject(f, content)
}
