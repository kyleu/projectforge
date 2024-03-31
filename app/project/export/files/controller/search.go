package controller

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Search(args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile("search", []string{"app", "lib", "search"}, "generated")
	if args.Models.HasSearch() {
		g.AddImport(helper.ImpContext, helper.ImpApp, helper.ImpAppUtil, helper.ImpCutil, helper.ImpLo, helper.ImpSearchResult)
	}
	lo.ForEach(args.Models, func(m *model.Model, _ int) {
		if m.HasSearches() {
			g.AddImport(helper.AppImport(m.PackageWithGroup("")))
		}
	})
	g.AddBlocks(searchBlock(args))
	return g.Render(addHeader, linebreak)
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
	ret.Pushf("\t%sFunc := func(ctx context.Context, params *Params, as *app.State, page *cutil.PageState, logger util.Logger) (result.Results, error) {", m.Package)
	if !m.HasTag("public") {
		ret.Push("\t\tif !page.Admin {")
		ret.Push("\t\t\treturn nil, nil")
		ret.Push("\t\t}")
	}
	ret.Pushf("\t\tprm := params.PS.Sanitized(%q, logger).WithLimit(5)", m.Package)
	const msg = "\t\tmodels, err := as.Services.%s.Search(ctx, params.Q, nil, prm%s, logger)"
	ret.Pushf(msg, m.Proper(), m.SoftDeleteSuffix())
	ret.Push("\t\tif err != nil {")
	ret.Push("\t\t\treturn nil, err")
	ret.Push("\t\t}")

	ret.Pushf("\t\treturn lo.Map(models, func(m *%s.%s, _ int) *result.Result {", m.Package, m.Proper())
	data := "m"
	if m.HasTag("big") {
		data = "nil"
	}
	icon := fmt.Sprintf("%q", m.Icon)
	if icons := m.Columns.WithFormat("icon"); len(icons) == 1 {
		icon = fmt.Sprintf("%s.%s", data, icons[0].Proper())
	}
	ret.Pushf("\t\t\treturn result.NewResult(%q, m.String(), m.WebPath(), m.TitleString(), %s, m, %s, params.Q)", m.Package, icon, data)
	ret.Push("\t\t}), nil")
	ret.Push("\t}")
	return ret.Slice
}
