package inject

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func Routes(f *file.File, args *model.Args) error {
	out := make([]string, 0, len(args.Models))
	for _, m := range args.Models {
		pkNames := make([]string, 0, len(m.Columns.PKs()))
		for _, col := range m.Columns.PKs() {
			pkNames = append(pkNames, fmt.Sprintf("{%s}", col.Name))
		}

		l := fmt.Sprintf("r.GET(\"/%s\", %sList)", m.Package, m.PackageProper())
		nf := fmt.Sprintf("r.GET(\"/%s/new\", %sCreateForm)", m.Package, m.PackageProper())
		ns := fmt.Sprintf("r.POST(\"/%s/new\", %sCreate)", m.Package, m.PackageProper())
		d := fmt.Sprintf("r.GET(\"/%s/%s\", %sDetail)", m.Package, strings.Join(pkNames, "/"), m.PackageProper())
		ef := fmt.Sprintf("r.GET(\"/%s/%s/edit\", %sEditForm)", m.Package, strings.Join(pkNames, "/"), m.PackageProper())
		es := fmt.Sprintf("r.POST(\"/%s/%s/edit\", %sEdit)", m.Package, strings.Join(pkNames, "/"), m.PackageProper())
		out = append(out, l, nf, ns, d, ef, es)
	}
	content := map[string]string{"codegen": "\n\t" + strings.Join(out, "\n\t") + "\n\t// "}
	return file.Inject(f, content)
}
