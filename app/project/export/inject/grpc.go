package inject

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/grpc"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func GRPC(f *file.File, args *model.Args, linebreak string) error {
	if len(args.Models) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models)*6)
	for _, m := range args.Models {
		fileArgs, err := grpc.GetGRPCFileArgs(m, args)
		if err != nil {
			return errors.Wrap(err, "invalid arguments")
		}
		lo.ForEach(fileArgs, func(fa *grpc.FileArgs, _ int) {
			out = append(out, grpcAll(m, fa, linebreak)...)
		})
	}
	content := map[string]string{"codegen": linebreak + strings.Join(out, linebreak) + linebreak + "\t// "}
	return file.Inject(f, content)
}

func grpcAll(m *model.Model, fa *grpc.FileArgs, linebreak string) []string {
	ret := []string{grpcList(m, fa, linebreak)}
	if m.HasSearches() {
		ret = append(ret, grpcSearch(m, fa, linebreak))
	}
	ret = append(ret,
		grpcDetail(m, fa, linebreak), grpcCreate(m, fa, linebreak), grpcUpdate(m, fa, linebreak),
		grpcSave(m, fa, linebreak), grpcDelete(m, fa, linebreak),
	)
	return ret
}

func grpcList(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("list", "inject")
	f.W("\tcase \"%s.list%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sList%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}

func grpcSearch(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("search", "inject")
	f.W("\tcase \"%s.search%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sSearch%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}

func grpcDetail(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("detail", "inject")
	f.W("\tcase \"%s.detail%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sDetail%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}

func grpcCreate(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("create", "inject")
	f.W("\tcase \"%s.create%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sCreate%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}

func grpcUpdate(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("update", "inject")
	f.W("\tcase \"%s.update%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sUpdate%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}

func grpcSave(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("save", "inject")
	f.W("\tcase \"%s.save%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sSave%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}

func grpcDelete(m *model.Model, fa *grpc.FileArgs, linebreak string) string {
	f := golang.NewBlock("delete", "inject")
	f.W("\tcase \"%s.delete%s\":", m.Package, fa.KeySuffix())
	f.W("\t\treturn %sDelete%s(p)", m.Proper(), fa.APISuffix())
	return f.Render(linebreak)
}
