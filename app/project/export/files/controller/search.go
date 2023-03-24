package controller

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func Search(args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile("search", []string{"app", "lib", "search"}, "generated")
	if args.Models.HasSearch() {
		g.AddImport(helper.ImpContext, helper.ImpApp, helper.ImpAppUtil, helper.ImpCutil, helper.AppImport("app/lib/search/result"))
	}
	g.AddBlocks(searchBlock(args))
	return g.Render(addHeader)
}

func searchBlock(args *model.Args) *golang.Block {
	ret := golang.NewBlock("menu", "func")
	keys := make([]string, 0, len(args.Models))
	ret.W("//nolint:gocognit")
	ret.W("func generatedSearch() []Provider {")
	for _, m := range args.Models {
		if len(m.AllSearches()) > 0 {
			keys = append(keys, m.Package+"Func")
			out := searchModel(m)
			for _, line := range out {
				ret.W(line)
			}
		}
	}
	ret.W("\treturn []Provider{" + strings.Join(keys, ", ") + "}")
	ret.W("}")
	return ret
}

func searchModel(m *model.Model) []string {
	var ret []string
	add := func(s string, args ...any) {
		ret = append(ret, fmt.Sprintf(s, args...))
	}
	add("\t%sFunc := func(ctx context.Context, params *Params, as *app.State, page *cutil.PageState, logger util.Logger) (result.Results, error) {", m.Package)
	if !m.HasTag("public") {
		add("\t\tif !page.Admin {")
		add("\t\t\treturn nil, nil")
		add("\t\t}")
	}
	add("\t\tprm := params.PS.Get(%q, nil, logger).Sanitize(%q).WithLimit(5)", m.Package, m.Package)
	const msg = "\t\tmodels, err := as.Services.%s.Search(ctx, params.Q, nil, prm%s, logger)"
	add(msg, m.Proper(), m.SoftDeleteSuffix())
	add("\t\tif err != nil {")
	add("\t\t\treturn nil, err")
	add("\t\t}")
	add("\t\tres := make(result.Results, 0, len(models))")
	add("\t\tfor _, m := range models {")
	data := "m"
	if m.HasTag("big") {
		data = "nil"
	}
	add("\t\t\tres = append(res, result.NewResult(%q, m.String(), m.WebPath(), m.String(), %q, m, %s, params.Q))", m.Package, m.Icon, data)
	add("\t\t}")
	add("\t\treturn res, nil")
	add("\t}")
	return ret
}
