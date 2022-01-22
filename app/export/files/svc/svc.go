package svc

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
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

	g.AddBlocks(serviceStruct(args), serviceNew(m, args), serviceDefaultFilters(m))
	return g.Render(addHeader)
}

func serviceStruct(args *model.Args) *golang.Block {
	ret := golang.NewBlock("Service", "struct")
	ret.W("type Service struct {")
	ret.W("\tdb     *database.Service")
	ret.W("\tlogger *zap.SugaredLogger")
	ret.W("}")
	return ret
}

func serviceNew(m *model.Model, args *model.Args) *golang.Block {
	ret := golang.NewBlock("NewService", "func")
	ret.W("func NewService(db *database.Service, logger *zap.SugaredLogger) *Service {")
	ret.W("\tlogger = logger.With(\"svc\", %q)", m.Camel())
	ret.W("\tfilter.AllowedColumns[\"%s\"] = columns", m.Package)
	ret.W("\treturn &Service{db: db, logger: logger}")
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
