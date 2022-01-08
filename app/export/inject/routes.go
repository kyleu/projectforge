package inject

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func Routes(f *file.File, args *model.Args) error {
	out := make([]string, 0, 100)
	for _, m := range args.Models {
		pkNames := make([]string, 0, len(m.PKs()))
		for _, col := range m.PKs() {
			pkNames = append(pkNames, fmt.Sprintf("{%s}", col.Name))
		}

		for _, grp := range m.GroupedColumns() {
			pathExtra := fmt.Sprintf("/%s/{%s}", grp.Camel(), grp.Camel())
			callSuffix := fmt.Sprintf("By%s", grp.Proper())

			g := fmt.Sprintf("\tr.GET(\"/%s/%s\", %s%sList)", m.Package, grp.Camel(), m.PackageProper(), grp.Proper())
			l := fmt.Sprintf("\tr.GET(\"/%s%s\", %sList%s)", m.Package, pathExtra, m.PackageProper(), callSuffix)
			nf := fmt.Sprintf("\tr.GET(\"/%s%s/new\", %sCreateForm%s)", m.Package, pathExtra, m.PackageProper(), callSuffix)
			ns := fmt.Sprintf("\tr.POST(\"/%s%s/new\", %sCreate%s)", m.Package, pathExtra, m.PackageProper(), callSuffix)
			d := fmt.Sprintf("\tr.GET(\"/%s%s/%s\", %sDetail%s)", m.Package, pathExtra, strings.Join(pkNames, "/"), m.PackageProper(), callSuffix)
			ef := fmt.Sprintf("\tr.GET(\"/%s%s/%s/edit\", %sEditForm%s)", m.Package, pathExtra, strings.Join(pkNames, "/"), m.PackageProper(), callSuffix)
			es := fmt.Sprintf("\tr.POST(\"/%s%s/%s/edit\", %sEdit%s)", m.Package, pathExtra, strings.Join(pkNames, "/"), m.PackageProper(), callSuffix)
			out = append(out, g, l, nf, ns, d, ef, es, "")
		}

		l := fmt.Sprintf("\tr.GET(\"/%s\", %sList)", m.Package, m.PackageProper())
		nf := fmt.Sprintf("\tr.GET(\"/%s/new\", %sCreateForm)", m.Package, m.PackageProper())
		ns := fmt.Sprintf("\tr.POST(\"/%s/new\", %sCreate)", m.Package, m.PackageProper())
		d := fmt.Sprintf("\tr.GET(\"/%s/%s\", %sDetail)", m.Package, strings.Join(pkNames, "/"), m.PackageProper())
		ef := fmt.Sprintf("\tr.GET(\"/%s/%s/edit\", %sEditForm)", m.Package, strings.Join(pkNames, "/"), m.PackageProper())
		es := fmt.Sprintf("\tr.POST(\"/%s/%s/edit\", %sEdit)", m.Package, strings.Join(pkNames, "/"), m.PackageProper())
		out = append(out, l, nf, ns, d, ef, es)
		if m.IsRevision() {
			rc := m.HistoryColumn()
			msg := "\tr.GET(\"/%s/%s/%s/{%s}\", %s%s)"
			out = append(out, fmt.Sprintf(msg, m.Package, strings.Join(pkNames, "/"), rc.Name, rc.Name, m.PackageProper(), rc.Proper()))
		}
	}
	content := map[string]string{"codegen": "\n" + strings.Join(out, "\n") + "\n\t// "}
	return file.Inject(f, content)
}
