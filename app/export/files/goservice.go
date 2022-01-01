package files

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func Service(m *model.Model, args *model.Args) *file.File {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, "service")
	for _, imp := range importsForTypes("go", m.Columns.PKs().Types()...) {
		g.AddImport(imp.Type, imp.Value)
	}
	g.AddImport(golang.ImportTypeInternal, "context")
	g.AddImport(golang.ImportTypeExternal, "go.uber.org/zap")
	g.AddImport(golang.ImportTypeExternal, "github.com/pkg/errors")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/lib/filter")

	if args.HasModule("database") {
		g.AddImport(golang.ImportTypeExternal, "github.com/jmoiron/sqlx")
		g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/lib/database")
	}

	g.AddBlocks(serviceStruct(args), serviceNew(m, args), serviceList(m))
	if len(m.Columns.PKs()) > 0 {
		g.AddBlocks(serviceGetByPK(m))
	}
	if len(m.Columns.PKs()) > 1 {
		for _, pkCol := range m.Columns.PKs() {
			g.AddBlocks(serviceGetBy("GetBy"+pkCol.Proper(), m, model.Columns{pkCol}, true))
		}
	}
	if len(m.Search) > 0 {
		g.AddImport(golang.ImportTypeInternal, "strings")
		g.AddBlocks(serviceSearch(m))
	}
	g.AddBlocks(serviceAdd(m), serviceUpdate(m), serviceSave(m), serviceDelete(m))
	return g.Render()
}

func serviceStruct(args *model.Args) *golang.Block {
	ret := golang.NewBlock("Service", "struct")
	ret.W("type Service struct {")
	if args.HasModule("database") {
		ret.W("\tdb     *database.Service")
	}
	ret.W("\tlogger *zap.SugaredLogger")
	ret.W("}")
	return ret
}

func serviceNew(m *model.Model, args *model.Args) *golang.Block {
	ret := golang.NewBlock("NewService", "func")
	logRebuild := fmt.Sprintf("\tlogger = logger.With(\"svc\", %q)", m.Camel())
	if args.HasModule("database") {
		ret.W("func NewService(db *database.Service, logger *zap.SugaredLogger) *Service {")
		ret.W(logRebuild)
		ret.W("\tfilter.AllowedColumns[\"%s\"] = Columns", m.Package)
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

func serviceList(m *model.Model) *golang.Block {
	ret := golang.NewBlock("List", "func")

	ret.W("func (s *Service) List(ctx context.Context, tx *sqlx.Tx, params *filter.Params) (%s, error) {", m.ProperPlural())
	ret.W("\tret := dtos{}")
	ret.W("\tsql := database.SQLSelect(ColumnsString, Table, \"\", params.OrderByString(), params.Limit, params.Offset)")
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx)")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", m.ProperPlural())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceSearch(m *model.Model) *golang.Block {
	ret := golang.NewBlock("search", "func")
	ret.W("func (s *Service) Search(ctx context.Context, q string, tx *sqlx.Tx, params *filter.Params) (%s, error) {", m.ProperPlural())
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
			like = eq
		}
	}
	for _, s := range m.Search {
		if strings.HasPrefix(s, "=") {
			clauses = append(clauses, s+" = "+eq)
		} else {
			clauses = append(clauses, s+" like "+like)
		}
	}
	wc := strings.Join(clauses, " or ")
	ret.W("\tsql := database.SQLSelect(ColumnsString, Table, \"" + wc + "\", params.OrderByString(), params.Limit, params.Offset)")
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx, " + strings.Join(params, ", ") + ")")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceGetByPK(m *model.Model) *golang.Block {
	return serviceGetBy("Get", m, m.Columns.PKs(), false)
}

func serviceGetBy(key string, m *model.Model, cols model.Columns, returnMultiple bool) *golang.Block {
	if returnMultiple {
		return serviceGetMultiple(key, m, cols)
	}
	return serviceGetOne(key, m, cols)
}

func serviceGetOne(key string, m *model.Model, cols model.Columns) *golang.Block {
	if key == "" {
		key = "GetBy" + cols.Smushed()
	}
	ret := golang.NewBlock(key, "func")
	ret.W("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s) (*%s, error) {", key, cols.Args(), m.Proper())
	ret.W("\tret := &dto{}")
	ret.W("\tsql := database.SQLSelectSimple(ColumnsString, Table, %q)", cols.WhereClause())
	ret.W("\terr := s.db.Get(ctx, ret, sql, tx, %s)", strings.Join(cols.CamelNames(), ", "))
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.Proper())
	ret.W("}")
	return ret
}

func serviceGetMultiple(key string, m *model.Model, cols model.Columns) *golang.Block {
	if key == "" {
		key = "GetBy" + cols.Smushed()
	}
	ret := golang.NewBlock(key, "func")
	ret.W("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s, params *filter.Params) (%s, error) {", key, cols.Args(), m.ProperPlural())
	ret.W("\tret := dtos{}")
	ret.W("\tsql := database.SQLSelect(ColumnsString, Table, %q, params.OrderByString(), params.Limit, params.Offset)", cols.WhereClause())
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx, %s)", strings.Join(cols.CamelNames(), ", "))
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceAdd(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Add", "func")
	ret.W("func (s *Service) Add(ctx context.Context, tx *sqlx.Tx, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")
	ret.W("\tq := database.SQLInsert(Table, Columns, len(models), \"\")")
	ret.W("\tvals := make([]interface{}, 0, len(models)*len(Columns))")
	ret.W("\tfor _, arg := range models {")
	ret.W("\t\tvals = append(vals, arg.ToData()...)")
	ret.W("\t}")
	ret.W("\treturn s.db.Insert(ctx, q, tx, vals...)")
	ret.W("}")
	return ret
}

func serviceUpdate(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Update", "func")
	ret.W("func (s *Service) Update(ctx context.Context, tx *sqlx.Tx, model *%s) error {", m.Proper())
	pks := m.Columns.PKs()
	wc := make([]string, 0, len(pks))
	pkVals := make([]string, 0, len(pks))
	for idx, pk := range pks {
		wc = append(wc, fmt.Sprintf("%s = $%d", pk.Name, len(m.Columns)+idx+1))
		pkVals = append(pkVals, "model."+pk.Proper())
	}
	ret.W("\tq := database.SQLUpdate(Table, Columns, \"%s\", \"\")", strings.Join(wc, " and "))
	ret.W("\tdata := model.ToData()")
	ret.W("\tdata = append(data, %s)", strings.Join(pkVals, ", "))
	ret.W("\t_, ret := s.db.Update(ctx, q, tx, 1, data...)")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}

func serviceSave(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Save", "func")
	ret.W("func (s *Service) Save(ctx context.Context, tx *sqlx.Tx, models ...*%s) error {", m.Proper())
	q := strings.Join(m.Columns.PKs().NamesQuoted(), ", ")
	ret.W("\tq := database.SQLUpsert(Table, Columns, len(models), []string{%s}, Columns, \"\")", q)
	ret.W("\tvar data []interface{}")
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tdata = append(data, model.ToData()...)")
	ret.W("\t}")
	ret.W("\treturn s.db.Insert(ctx, q, tx, data...)")
	ret.W("}")
	return ret
}

func serviceDelete(m *model.Model) *golang.Block {
	pks := m.Columns.PKs()
	whereClauses := make([]string, 0, len(pks))
	for idx, pk := range pks {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", pk.Name, idx+1))
	}
	ret := golang.NewBlock("Delete", "func")
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s) error {", pks.Args())
	ret.W("\tq := database.SQLDelete(Table, %q)", strings.Join(whereClauses, " and "))
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, 1, %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	return ret
}
