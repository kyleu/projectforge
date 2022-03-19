package inject

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

func Search(f *file.File, args *model.Args) error {
	var hasSearch bool
	for _, m := range args.Models {
		if len(m.Search) > 0 {
			hasSearch = true
			break
		}
	}
	if !hasSearch {
		return nil
	}
	out := make([]string, 0, len(args.Models))
	funcs := make([]string, 0, len(args.Models))
	for _, m := range args.Models {
		if len(m.Search) > 0 {
			out = append(out, searchModel(m))
			funcs = append(funcs, fmt.Sprintf("%sFunc", m.Package))
		}
	}
	out = append(out, "\tallProviders = append(allProviders, "+strings.Join(funcs, ", ")+")")
	content := map[string]string{"codegen": "\n" + strings.Join(out, "\n") + "\n\t// "}
	return file.Inject(f, content)
}

func searchModel(m *model.Model) string {
	f := golang.NewBlock("search", "inject")
	f.W("\t%sFunc := func(ctx context.Context, as *app.State, params *Params) (result.Results, error) {", m.Package)
	const msg = "\t\tmodels, err := as.Services.%s.Search(ctx, params.Q, nil, params.PS.Get(%q, nil, as.Logger)%s)"
	f.W(msg, m.Proper(), m.Package, m.SoftDeleteSuffix())
	f.W("\t\tif err != nil {")
	f.W("\t\t\treturn nil, errors.Wrap(err, \"\")")
	f.W("\t\t}")
	f.W("\t\tres := make(result.Results, 0, len(models))")
	f.W("\t\tfor _, m := range models {")
	f.W("\t\t\tres = append(res, result.NewResult(%q, m.String(), m.WebPath(), m.String(), %q, m, params.Q))", m.Package, m.Icon)
	f.W("\t\t}")
	f.W("\t\treturn res, nil")
	f.W("\t}")
	return f.Render()
}
