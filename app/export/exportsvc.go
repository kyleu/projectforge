package export

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
)

func exportServiceFile(m *Model, args *Args) *file.File {
	g := NewGoFile(m.Package, []string{"app", m.Package}, "service")
	for _, imp := range m.Columns.PKs().Types().Imports() {
		g.AddImport(imp.Type, imp.Value)
	}

	g.AddImport(ImportTypeInternal, "context")
	g.AddImport(ImportTypeExternal, "go.uber.org/zap")
	g.AddImport(ImportTypeExternal, "github.com/pkg/errors")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/filter")

	if args.HasModule("database") {
		g.AddImport(ImportTypeExternal, "github.com/jmoiron/sqlx")
		g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/database")
	}

	g.AddBlocks(serviceStruct(args), serviceNew(m, args), serviceList(m))
	if len(m.Columns.PKs()) > 0 {
		g.AddBlocks(serviceGetByPK(m))
	}
	if len(m.Columns.PKs()) > 1 {
		for _, pkCol := range m.Columns.PKs() {
			g.AddBlocks(serviceGetBy("GetBy"+pkCol.proper(), m, Columns{pkCol}, true))
		}
	}
	if len(m.Search) > 0 {
		g.AddImport(ImportTypeInternal, "strings")
		g.AddBlocks(serviceSearch(m))
	}
	return g.Render()
}

func serviceStruct(args *Args) *Block {
	ret := NewBlock("Service", "struct")
	ret.W("type Service struct {")
	if args.HasModule("database") {
		ret.W("\tdb     *database.Service")
	}
	ret.W("\tlogger *zap.SugaredLogger")
	ret.W("}")
	return ret
}

func serviceNew(m *Model, args *Args) *Block {
	ret := NewBlock("NewService", "func")
	logRebuild := fmt.Sprintf("\tlogger = logger.With(\"svc\", %q)", m.camel())
	if args.HasModule("database") {
		ret.W("func NewService(db *database.Service, logger *zap.SugaredLogger) *Service {")
		ret.W(logRebuild)
		ret.W("\treturn &Service{db: db, logger: logger}")
		ret.W("}")
	} else {
		ret.W("func NewService(logger *zap.SugaredLogger) *Service {")
		ret.W(logRebuild)
		ret.W("\treturn &Service{logger: logger}")
		ret.W("}")
	}
	return ret
}

func serviceList(m *Model) *Block {
	ret := NewBlock("List", "func")

	ret.WF("func (s *Service) List(ctx context.Context, tx *sqlx.Tx, params *filter.Params) (%s, error) {", m.properPlural())
	ret.W("\tret := dtos{}")
	ret.W("\tsql := database.SQLSelect(ColumnsString, Table, \"\", params.OrderByString(), params.Limit, params.Offset)")
	ret.W("\terr := s.db.Select(ctx, ret, sql, tx)")
	ret.W("\tif err != nil {")
	ret.WF("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", m.properPlural())
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.properPlural())
	ret.W("}")
	return ret
}

func serviceSearch(m *Model) *Block {
	ret := NewBlock("search", "func")
	ret.WF("func (s *Service) Search(ctx context.Context, q string, tx *sqlx.Tx, params *filter.Params) (%s, error) {", m.properPlural())
	ret.W("\tret := dtos{}")

	var clauses []string
	hasEqual, hasLike := false, false
	for _, s := range m.Search {
		if strings.HasPrefix(s, "=") {
			hasEqual = true
		} else {
			hasLike = true
		}
	}
	eq := "$1"
	like := "$2"
	var params []string
	if hasEqual {
		params = append(params, "strings.ToLower(q)")
	}
	if hasLike {
		params = append(params, `"%%"+strings.ToLower(q)+"%%"`)
		if !hasEqual {
			like = "$1"
		}
	}
	for _, s := range m.Search {
		if strings.HasPrefix(s, "=") {
			clauses = append(clauses, s+" = "+eq)
		} else {
			clauses = append(clauses, s+" like "+like)
		}
	}
	wc := strings.Join(clauses, " and ")
	ret.W("\tsql := database.SQLSelect(ColumnsString, Table, \"" + wc + "\", params.OrderByString(), params.Limit, params.Offset)")
	ret.W("\terr := s.db.Select(ctx, ret, sql, tx, " + strings.Join(params, ", ") + ")")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.properPlural())
	ret.W("}")
	return ret
}

func serviceGetByPK(m *Model) *Block {
	return serviceGetBy("Get", m, m.Columns.PKs(), false)
}

func serviceGetBy(key string, m *Model, cols Columns, returnMultiple bool) *Block {
	if returnMultiple {
		return serviceGetMultiple(key, m, cols)
	}
	return serviceGetOne(key, m, cols)
}

func serviceGetOne(key string, m *Model, cols Columns) *Block {
	if key == "" {
		key = "GetBy" + cols.Smushed()
	}
	ret := NewBlock(key, "func")
	ret.WF("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s) (*%s, error) {", key, cols.Args(), m.proper())
	ret.W("\tret := &dto{}")
	ret.WF("\tsql := database.SQLSelectSimple(ColumnsString, Table, %q)", cols.WhereClause())
	ret.WF("\terr := s.db.Get(ctx, ret, sql, tx, %s)", strings.Join(cols.camelNames(), ", "))
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.proper())
	ret.W("}")
	return ret
}

func serviceGetMultiple(key string, m *Model, cols Columns) *Block {
	if key == "" {
		key = "GetBy" + cols.Smushed()
	}
	ret := NewBlock(key, "func")
	ret.WF("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s, params *filter.Params) (%s, error) {", key, cols.Args(), m.properPlural())
	ret.W("\tret := dtos{}")
	ret.WF("\tsql := database.SQLSelect(ColumnsString, Table, %q, params.OrderByString(), params.Limit, params.Offset)", cols.WhereClause())
	ret.WF("\terr := s.db.Select(ctx, ret, sql, tx, %s)", strings.Join(cols.camelNames(), ", "))
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.properPlural())
	ret.W("}")
	return ret
}
