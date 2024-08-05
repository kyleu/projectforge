package controller

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Search(args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile("search", []string{"app", "lib", "search"}, "generated")
	if args.Models.HasSearch() {
		g.AddImport(helper.ImpContext, helper.ImpApp, helper.ImpAppUtil, helper.ImpCutil, helper.ImpSearchResult)
	}
	g.AddBlocks(searchBlock(args))
	return g.Render(linebreak)
}

func searchBlock(args *model.Args) *golang.Block {
	ret := golang.NewBlock("menu", "func")
	keys := make([]string, 0, len(args.Models))
	ret.W("func generatedSearch() []Provider {")
	lo.ForEach(args.Models, func(m *model.Model, _ int) {
		if m.HasSearches() {
			keys = append(keys, m.Package+"Func")
			out := searchModel(m)
			lo.ForEach(out, func(line string, _ int) {
				ret.W(line)
			})
		}
	})
	ret.W("\treturn []Provider{" + strings.Join(keys, ", ") + "}")
	ret.W("}")
	return ret
}

func searchModel(m *model.Model) []string {
	ret := &util.StringSlice{}
	fMsg := "\t%sFunc := func(ctx context.Context, params *Params, as *app.State, page *cutil.PageState, logger util.Logger) (result.Results, error) {"
	ret.Pushf(fMsg, m.Package)
	if !m.HasTag("public") {
		ret.Push("\t\tif !page.Admin {")
		ret.Push("\t\t\treturn nil, nil")
		ret.Push("\t\t}")
	}
	ret.Pushf("\t\tprm := params.PS.Sanitized(%q, logger).WithLimit(5)", m.Package)
	const msg = "\t\treturn as.Services.%s.SearchEntries(ctx, params.Q, nil, prm%s, logger)"
	ret.Pushf(msg, m.Proper(), m.SoftDeleteSuffix())
	ret.Push("\t}")
	return ret.Slice
}
