package svc

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

func ServiceAll(m *model.Model, args *model.Args, addHeader bool) (file.Files, error) {
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

func Service(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, "service")
	g.AddImport(helper.ImpLogging, helper.ImpFilter, helper.ImpDatabase)

	isAudit := args.HasModule("audit") && m.HasTag("audit")
	if isAudit {
		g.AddImport(helper.ImpAudit)
	}

	g.AddBlocks(serviceStruct(isAudit), serviceNew(m, isAudit), serviceDefaultFilters(m))
	return g.Render(addHeader)
}

func serviceStruct(isAudit bool) *golang.Block {
	ret := golang.NewBlock("Service", "struct")
	ret.W("type Service struct {")
	ret.W("\tdb     *database.Service")
	if isAudit {
		ret.W("\taudit  *audit.Service")
	}
	ret.W("\tlogger *zap.SugaredLogger")
	ret.W("}")
	return ret
}

func serviceNew(m *model.Model, isAudit bool) *golang.Block {
	ret := golang.NewBlock("NewService", "func")
	newSuffix, callSuffix := "", ""
	if isAudit {
		newSuffix = "aud *audit.Service, "
		callSuffix = "audit: aud, "
	}
	ret.W("func NewService(db *database.Service, %slogger *zap.SugaredLogger) *Service {", newSuffix)
	ret.W("\tlogger = logger.With(\"svc\", %q)", m.Camel())
	ret.W("\tfilter.AllowedColumns[\"%s\"] = columns", m.Package)
	ret.W("\treturn &Service{db: db, %slogger: logger}", callSuffix)
	ret.W("}")
	return ret
}

func serviceDefaultFilters(m *model.Model) *golang.Block {
	ret := golang.NewBlock("NewService", "func")
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

func getSuffix(m *model.Model) string {
	if m.IsSoftDelete() {
		return ", includeDeleted bool"
	}
	return ""
}
