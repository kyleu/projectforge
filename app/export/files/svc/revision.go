package svc

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func ServiceRevision(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, "servicerevision")
	g.AddImport(helper.ImpFmt, helper.ImpStrings, helper.ImpContext, helper.ImpFilter, helper.ImpSQLx, helper.ImpErrors, helper.ImpDatabase)
	ar, err := serviceGetAllRevisions(m)
	if err != nil {
		return nil, err
	}
	gr, err := serviceGetRevision(m)
	if err != nil {
		return nil, err
	}
	gnr, err := serviceGetCurrentRevisions(m)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(ar, gr, gnr)
	return g.Render()
}

func serviceGetAllRevisions(m *model.Model) (*golang.Block, error) {
	hc := m.HistoryColumns(true)
	pks := m.PKs()

	ret := golang.NewBlock(fmt.Sprintf("GetAll%s", hc.Col.ProperPlural()), "func")
	decl := "func (s *Service) GetAll%s(ctx context.Context, tx *sqlx.Tx, %s, params *filter.Params, includeDeleted bool) (%s, error) {"
	ret.W(decl, hc.Col.ProperPlural(), pks.Args(), m.ProperPlural())
	ret.W("\tparams = filters(params)")
	placeholders := make([]string, 0, len(m.PKs()))
	for idx, pk := range m.PKs() {
		placeholders = append(placeholders, fmt.Sprintf("\\\"%s\\\" = $%d", pk.Name, idx+1))
	}
	ret.W("\twc := \"%s\"", strings.Join(placeholders, " and "))
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	err := addJoinClause(ret, m, hc)
	if err != nil {
		return nil, err
	}
	ret.W("\tq := database.SQLSelect(columnsString, tablesJoinedParam, wc, params.OrderByString(), params.Limit, params.Offset)")
	ret.W("\tret := dtos{}")
	ret.W("\terr := s.db.Select(ctx, &ret, q, tx, %s)", strings.Join(pks.Names(), ", "))
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", util.StringToPlural(m.Proper()))
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", util.StringToPlural(m.Proper()))
	ret.W("}")
	return ret, nil
}

func serviceGetRevision(m *model.Model) (*golang.Block, error) {
	revCol := m.HistoryColumn()
	ret := golang.NewBlock(fmt.Sprintf("Get%s", revCol.Proper()), "func")
	decl := "func (s *Service) Get%s(ctx context.Context, tx *sqlx.Tx, %s, %s int) (*%s, error) {"
	ret.W(decl, revCol.Proper(), m.PKs().Args(), revCol.Camel(), m.Proper())
	placeholders := make([]string, 0, len(m.PKs()))
	for idx, pk := range m.PKs() {
		placeholders = append(placeholders, fmt.Sprintf("\\\"%s\\\" = $%d", pk.Name, idx+1))
	}
	ret.W("\twc := \"%s and \\\"%s\\\" = $%d\"", strings.Join(placeholders, " and "), revCol.Name, len(m.PKs())+1)
	ret.W("\tret := &dto{}")
	err := addJoinClause(ret, m, m.HistoryColumns(true))
	if err != nil {
		return nil, err
	}
	ret.W("\tq := database.SQLSelectSimple(columnsString, tablesJoinedParam, wc)")
	ret.W("\terr := s.db.Get(ctx, ret, q, tx, %s, %s)", strings.Join(m.PKs().Names(), ", "), revCol.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.Proper())
	ret.W("}")
	return ret, nil
}

func serviceGetCurrentRevisions(m *model.Model) (*golang.Block, error) {
	revCol := m.HistoryColumn()
	pks := m.PKs()
	pkWCStr := make([]string, 0, len(pks))
	pkWCIdx := make([]string, 0, len(pks))
	pkModelRefs := make([]string, 0, len(pks))
	pkComps := make([]string, 0, len(pks))
	for idx, pk := range pks {
		pkWCStr = append(pkWCStr, fmt.Sprintf("%q = $%%%%d", pk.Name))
		if len(pks) == 1 {
			pkWCIdx = append(pkWCIdx, "i+1")
		} else {
			pkWCIdx = append(pkWCIdx, fmt.Sprintf("(i*%d)+%d", len(pks), idx+1))
		}
		pkModelRefs = append(pkModelRefs, fmt.Sprintf("model.%s", pk.Proper()))
		pkComps = append(pkComps, fmt.Sprintf("x.%s == model.%s", pk.Proper(), pk.Proper()))
	}

	ret := golang.NewBlock(fmt.Sprintf("GetCurrent%s", revCol.ProperPlural()), "func")
	decl := "func (s *Service) getCurrent%s(ctx context.Context, tx *sqlx.Tx, models ...*%s) (map[string]%s, error) {"
	ret.W(decl, revCol.ProperPlural(), m.Proper(), revCol.Type.ToGoType(false))
	ret.W("\tstmts := make([]string, 0, len(models))")
	ret.W("\tfor i := range models {")
	ret.W("\t\tstmts = append(stmts, fmt.Sprintf(`%s`, %s))", strings.Join(pkWCStr, " and "), strings.Join(pkWCIdx, ", "))
	ret.W("\t}")
	ret.W("\tq := database.SQLSelectSimple(`%s, \"current_%s\"`, tableQuoted, strings.Join(stmts, \" or \"))", strings.Join(pks.NamesQuoted(), ", "), revCol.Name)
	ret.W("\tvals := make([]interface{}, 0, len(models))")
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tvals = append(vals, %s)", strings.Join(pkModelRefs, ", "))
	ret.W("\t}")
	ret.W("\tvar results []*struct {")
	maxColLength := pks.MaxGoKeyLength()
	maxTypeLength := pks.MaxGoTypeLength()
	currRevStr := fmt.Sprintf("Current%s", revCol.Proper())
	if maxColLength < len(currRevStr) {
		maxColLength = len(currRevStr)
	}
	for _, pk := range pks {
		ret.W("\t\t%s %s `db:%q`", util.StringPad(pk.Proper(), maxColLength), util.StringPad(pk.ToGoType(), maxTypeLength), pk.Name)
	}
	ret.W("\t\t%s int    `db:\"current_%s\"`", currRevStr, revCol.Name)
	ret.W("\t}")
	ret.W("\terr := s.db.Select(ctx, &results, q, tx, vals...)")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", m.ProperPlural())
	ret.W("\t}")
	ret.W("")
	ret.W("\tret := make(map[string]int, len(models))")
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tcurr := 0")
	ret.W("\t\tfor _, x := range results {")
	ret.W("\t\t\tif %s {", strings.Join(pkComps, " && "))
	ret.W("\t\t\t\tcurr = x.Current%s", revCol.Proper())
	ret.W("\t\t\t}")
	ret.W("\t\t}")
	ret.W("\t\tret[model.String()] = curr")
	ret.W("\t}")
	ret.W("\treturn ret, nil")
	ret.W("}")
	return ret, nil
}

func addJoinClause(ret *golang.Block, m *model.Model, hc *model.HistoryMap) error {
	joinClause := fmt.Sprintf("%%%%q %s join %%%%q %sr on ", m.FirstLetter(), m.FirstLetter())
	var joins []string
	for idx, col := range hc.Const {
		if col.PK {
			rCol := hc.Var[idx]
			if !(rCol.PK || rCol.HasTag(model.RevisionType)) {
				return errors.Errorf("invalid revision column [%s] at index [%d]", rCol.Name, idx)
			}
			joins = append(joins, fmt.Sprintf("%s.%s = %sr.%s", m.FirstLetter(), col.Name, m.FirstLetter(), rCol.Name))
		}
	}
	joinClause += strings.Join(joins, " and ")
	ret.W("\ttablesJoinedParam := fmt.Sprintf(%q, table, table%s)", joinClause, hc.Col.Proper())
	return nil
}
