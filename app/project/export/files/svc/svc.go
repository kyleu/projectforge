package svc

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	golang2 "projectforge.dev/projectforge/app/project/export/golang"
	model2 "projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func ServiceAll(m *model2.Model, args *model2.Args, addHeader bool) (file.Files, error) {
	x, err := Service(m, args, addHeader)
	if err != nil {
		return nil, err
	}
	g, err := ServiceGet(m, args, addHeader)
	if err != nil {
		return nil, err
	}
	mt, err := ServiceMutate(m, args, addHeader)
	if err != nil {
		return nil, err
	}
	ret := file.Files{x, g, mt}
	if m.IsRevision() {
		r, err := ServiceRevision(m, args, addHeader)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r)
	}
	if m.IsHistory() {
		r, err := ServiceHistory(m, args, addHeader)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func Service(m *model2.Model, args *model2.Args, addHeader bool) (*file.File, error) {
	g := golang2.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "service")
	g.AddImport(helper.ImpFilter, helper.ImpDatabase)

	isRO := args.HasModule("readonlydb")
	isAudit := args.HasModule("audit") && m.HasTag("audit")
	if isAudit {
		g.AddImport(helper.ImpAudit)
	}

	g.AddBlocks(serviceStruct(isRO, isAudit), serviceNew(m, isRO, isAudit), serviceDefaultFilters(m))
	return g.Render(addHeader)
}

func serviceStruct(isRO bool, isAudit bool) *golang2.Block {
	ret := golang2.NewBlock("Service", "struct")
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

func serviceNew(m *model2.Model, isRO bool, isAudit bool) *golang2.Block {
	ret := golang2.NewBlock("NewService", "func")
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

func serviceDefaultFilters(m *model2.Model) *golang2.Block {
	ret := golang2.NewBlock("NewService", "func")
	ret.W("func filters(orig *filter.Params) *filter.Params {")
	ords := make([]string, 0, len(m.Ordering))
	for _, ord := range m.Ordering {
		if ord.Asc {
			ords = append(ords, fmt.Sprintf("&filter.Ordering{Column: %q, Asc: true}", ord.Column))
		} else {
			ords = append(ords, fmt.Sprintf("&filter.Ordering{Column: %q}", ord.Column))
		}
	}
	if len(ords) == 0 {
		ret.W("\treturn orig.Sanitize(%q)", m.Package)
	} else {
		ret.W("\treturn orig.Sanitize(%q, %s)", m.Package, strings.Join(ords, ", "))
	}
	ret.W("}")
	return ret
}

func getSuffix(m *model2.Model) string {
	if m.IsSoftDelete() {
		return incDel
	}
	return ""
}
