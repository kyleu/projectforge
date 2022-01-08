package inject

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/files/grpc"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/pkg/errors"
)

func GRPC(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models)*6)
	for _, m := range args.Models {
		fileArgs, err := grpc.GetGRPCFileArgs(m, args)
		if err != nil {
			return errors.Wrap(err, "invalid arguments")
		}
		for _, fa := range fileArgs {
			out = append(out, grpcAll(m, fa)...)
		}
	}
	content := map[string]string{"codegen": "\n" + strings.Join(out, "\n") + "\n\t// "}
	return file.Inject(f, content)
}

func grpcAll(m *model.Model, fa *grpc.GRPCFileArgs) []string {
	return []string{
		grpcList(m, fa), grpcSearch(m, fa), grpcDetail(m, fa),
		grpcCreate(m, fa), grpcUpdate(m, fa), grpcSave(m, fa), grpcDelete(m, fa),
	}
}

func grpcList(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("list", "inject")
	f.W("\tcase \"%s.list%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sList%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}

func grpcSearch(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("search", "inject")
	f.W("\tcase \"%s.search%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sSearch%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}

func grpcDetail(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("detail", "inject")
	f.W("\tcase \"%s.detail%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sDetail%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}

func grpcCreate(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("create", "inject")
	f.W("\tcase \"%s.create%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sCreate%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}

func grpcUpdate(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("update", "inject")
	f.W("\tcase \"%s.update%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sUpdate%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}

func grpcSave(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("save", "inject")
	f.W("\tcase \"%s.save%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sSave%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}

func grpcDelete(m *model.Model, fa *grpc.GRPCFileArgs) string {
	f := golang.NewBlock("delete", "inject")
	f.W("\tcase \"%s.delete%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sDelete%s(p)", m.PackageProper(), fa.APISuffix())
	return f.Render()
}
