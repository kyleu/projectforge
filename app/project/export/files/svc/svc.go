package svc

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const incDel = ", includeDeleted bool"

func ServiceAll(m *model.Model, args *model.Args, addHeader bool, linebreak string) (file.Files, error) {
	x, err := Service(m, args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	g, err := ServiceGet(m, args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	l, err := ServiceList(m, args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	mt, err := ServiceMutate(m, args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	return file.Files{x, g, l, mt}, nil
}

func Service(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "service")
	g.AddImport(helper.ImpFilter, helper.ImpAppDatabase, helper.ImpAppSvc)
	g.AddImport(m.Imports.Supporting("service")...)

	isRO := args.HasModule("readonlydb")
	isAudit := args.HasModule("audit") && m.HasTag("audit")
	if isAudit {
		g.AddImport(helper.ImpAudit)
	}

	g.AddBlocks(typeAssert(g, m, args.Enums), serviceStruct(isRO, isAudit), serviceNew(m, isRO, isAudit), serviceDefaultFilters(m))
	return g.Render(addHeader, linebreak)
}

func typeAssert(g *golang.File, m *model.Model, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("assert", "type")
	suffix := ""
	if m.IsSoftDelete() {
		suffix += "SoftDelete"
	}
	var assertion string
	if len(m.PKs()) == 1 {
		lo.ForEach(helper.ImportsForTypes("go", "", m.PKs().Types()...), func(imp *model.Import, _ int) {
			g.AddImport(imp)
		})
		pk := m.PKs()[0]
		args, _ := pk.ToGoType(m.Package, enums)
		assertion = fmt.Sprintf("svc.Service%sID[*%s, %s, %s] = (*Service)(nil)", suffix, m.Proper(), m.ProperPlural(), args)
	} else {
		assertion = fmt.Sprintf("svc.Service%s[*%s, %s] = (*Service)(nil)", suffix, m.Proper(), m.ProperPlural())
	}
	if m.HasSearches() {
		ret.W("var (")
		ret.W("\t_ %s", assertion)
		ss := fmt.Sprintf("svc.ServiceSearch[%s]", m.ProperPlural())
		ret.W("\t_ %s = (*Service)(nil)", util.StringPad(ss, len(assertion)-18))
		ret.W(")")
	} else {
		ret.W("var _ %s", assertion)
	}
	return ret
}

func serviceStruct(isRO bool, isAudit bool) *golang.Block {
	ret := golang.NewBlock("Service", "struct")
	ret.W("type Service struct {")
	size := 2
	if isAudit {
		size = 5
	}
	if isRO {
		size = 6
	}
	ret.W("\t%s *database.Service", util.StringPad("db", size))
	if isRO {
		ret.W("\t%s *database.Service", util.StringPad("dbRead", size))
	}
	if isAudit {
		ret.W("\t%s *audit.Service", util.StringPad("audit", size))
	}
	ret.W("}")
	return ret
}

func serviceNew(m *model.Model, isRO bool, isAudit bool) *golang.Block {
	ret := golang.NewBlock("NewService", "func")
	newSuffix, callSuffix := "", ""
	if isRO {
		newSuffix = ", dbRead *database.Service"
		callSuffix = ", dbRead: dbRead"
	}
	if isAudit {
		newSuffix = ", aud *audit.Service"
		callSuffix = ", audit: aud"
	}
	ret.W("func NewService(db *database.Service%s) *Service {", newSuffix)
	ret.W("\tfilter.AllowedColumns[\"%s\"] = columns", m.Package)
	ret.W("\treturn &Service{db: db%s}", callSuffix)
	ret.W("}")
	return ret
}

func serviceDefaultFilters(m *model.Model) *golang.Block {
	ret := golang.NewBlock("NewService", "func")
	ret.W("func filters(orig *filter.Params) *filter.Params {")
	ords := make([]string, 0, len(m.Ordering))
	lo.ForEach(m.Ordering, func(ord *filter.Ordering, _ int) {
		if ord.Asc {
			ords = append(ords, fmt.Sprintf("&filter.Ordering{Column: %q, Asc: true}", ord.Column))
		} else {
			ords = append(ords, fmt.Sprintf("&filter.Ordering{Column: %q}", ord.Column))
		}
	})
	if len(ords) == 0 {
		ret.W("\treturn orig.Sanitize(%q)", m.Package)
	} else {
		ret.W("\treturn orig.Sanitize(%q, %s)", m.Package, strings.Join(ords, ", "))
	}
	ret.W("}")
	return ret
}

func getSuffix(m *model.Model) string {
	if m.IsSoftDelete() {
		return incDel
	}
	return ""
}
