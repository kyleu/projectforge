package svc

import (
	"fmt"
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
	"strings"
)

func Services(args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile("app", []string{"app"}, "generated")
	g.AddImport(helper.ImpContext, helper.ImpAppDatabase, helper.ImpAppUtil)
	if args.HasModule("audit") {
		g.AddImport(helper.ImpAudit)
	}
	for _, m := range args.Models {
		g.AddImport(helper.AppImport(m.PackageWithGroup("")))
	}

	svcSize := 0
	lo.ForEach(args.Models, func(m *model.Model, _ int) {
		if len(m.Proper()) > svcSize {
			svcSize = len(m.Proper())
		}
	})

	initParamsArr := []string{"dbSvc *database.Service"}
	if args.HasModule("readonlydb") {
		initParamsArr = append(initParamsArr, "dbReadSvc *database.Service")
	}
	if args.HasModule("audit") {
		initParamsArr = append(initParamsArr, "audSvc *audit.Service")
	}
	initParams := strings.Join(initParamsArr, ", ")

	svcs := make([]string, 0, len(args.Models))
	refs := make([]string, 0, len(args.Models))
	lo.ForEach(args.Models, func(m *model.Model, _ int) {
		svcs = append(svcs, fmt.Sprintf("%s *%s.Service", util.StringPad(m.Proper(), svcSize), m.Package))

		refParamsArr := []string{"dbSvc"}
		if args.HasModule("readonlydb") && !m.HasTag("audit") {
			refParamsArr = append(refParamsArr, "dbReadSvc")
		}
		if args.HasModule("audit") && m.HasTag("audit") {
			refParamsArr = append(refParamsArr, "audSvc")
		}
		refParams := strings.Join(refParamsArr, ", ")

		refs = append(refs, fmt.Sprintf("%s %s.NewService(%s),", util.StringPad(m.Proper()+":", svcSize+1), m.Package, refParams))
	})

	g.AddBlocks(servicesStruct(svcs), servicesInitFn(refs, initParams))
	return g.Render(addHeader, linebreak)
}

func servicesStruct(svcs []string) *golang.Block {
	ret := golang.NewBlock("genStruct", "struct")
	ret.W("type GeneratedServices struct {")
	for _, svc := range svcs {
		ret.W("\t" + svc)
	}
	ret.W("}")
	return ret
}

func servicesInitFn(refs []string, params string) *golang.Block {
	ret := golang.NewBlock("initGeneratedServices", "func")
	ret.W("func initGeneratedServices(ctx context.Context, %s, logger util.Logger) GeneratedServices {", params)
	ret.W("\treturn GeneratedServices{")
	for _, svc := range refs {
		ret.W("\t\t" + svc)
	}
	ret.W("\t}")
	ret.W("}")
	return ret
}
