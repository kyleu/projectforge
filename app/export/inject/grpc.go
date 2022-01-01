package inject

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func GRPC(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models)*6)
	for _, m := range args.Models {
		out = append(out, grpcList(m), grpcDetail(m), grpcAdd(m), grpcUpdate(m), grpcSave(m), grpcDelete(m))
	}
	content := map[string]string{"codegen": "\n" + strings.Join(out, "\n") + "\n\t// "}
	return file.Inject(f, content)
}

func grpcList(m *model.Model) string {
	f := golang.NewBlock("list", "inject")
	f.W("\tcase \"%s.list\":", m.Package)
	f.W("\t\treturn %sList(p)", m.PackageProper())
	return f.Render()
}

func grpcDetail(m *model.Model) string {
	f := golang.NewBlock("detail", "inject")
	f.W("\tcase \"%s.detail\":", m.Package)
	f.W("\t\treturn %sDetail(p)", m.PackageProper())
	return f.Render()
}

func grpcAdd(m *model.Model) string {
	f := golang.NewBlock("add", "inject")
	f.W("\tcase \"%s.add\":", m.Package)
	f.W("\t\treturn %sAdd(p)", m.PackageProper())
	return f.Render()
}

func grpcUpdate(m *model.Model) string {
	f := golang.NewBlock("update", "inject")
	f.W("\tcase \"%s.update\":", m.Package)
	f.W("\t\treturn %sUpdate(p)", m.PackageProper())
	return f.Render()
}

func grpcSave(m *model.Model) string {
	f := golang.NewBlock("save", "inject")
	f.W("\tcase \"%s.save\":", m.Package)
	f.W("\t\treturn %sSave(p)", m.PackageProper())
	return f.Render()
}

func grpcDelete(m *model.Model) string {
	f := golang.NewBlock("delete", "inject")
	f.W("\tcase \"%s.delete\":", m.Package)
	f.W("\t\treturn %sDelete(p)", m.PackageProper())
	return f.Render()
}
