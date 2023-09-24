package sql

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func sqlCreateRevision(m *model.Model, modules []string, now string, models model.Models) (*golang.Block, error) {
	pks := m.PKs()
	hc := m.HistoryColumns(true)

	ret := golang.NewBlock("SQLCreateRev", "sql")
	ret.W("-- {%% func " + m.Proper() + "Create() %%}")

	// core
	ret.W("create table if not exists %q (", m.Name)
	for _, col := range hc.Const {
		st, err := col.ToSQLType()
		if err != nil {
			return nil, err
		}
		ret.W("  %q %s,", col.Name, st)
	}
	ret.W("  primary key (%s)", strings.Join(pks.NamesQuoted(), ", "))
	ret.W(");")

	if len(pks) > 1 {
		lo.ForEach(pks, func(pk *model.Column, _ int) {
			addIndex(ret, m.Name, pk.Name)
		})
	}

	// revision
	ret.WB()
	revTblName := m.Name + "_" + hc.Col.Name
	ret.W("create table if not exists %q (", revTblName)
	for _, col := range hc.Var {
		st, err := col.ToSQLType()
		if err != nil {
			return nil, err
		}
		ret.W("  %q %s,", col.Name, st)
	}

	revPKs := hc.Var.PKs()
	revPKsWithRev := append(model.Columns{}, revPKs...)
	revPKsWithRev = append(revPKsWithRev, hc.Col)

	bareRefs := strings.Join(revPKs.NamesQuoted(), ", ")
	ret.W("  foreign key (%s) references %q (%s),", bareRefs, m.Name, strings.Join(pks.NamesQuoted(), ", "))
	sqlRelations(ret, m, models)
	ret.W("  primary key (%s)", strings.Join(revPKsWithRev.NamesQuoted(), ", "))
	ret.W(");")

	if len(revPKs) > 1 {
		addIndex(ret, revTblName, revPKs.Names()...)
	}

	lo.ForEach(revPKsWithRev, func(pk *model.Column, _ int) {
		if !pk.HasTag(model.RevisionType) {
			addIndex(ret, revTblName, pk.Name)
		}
	})

	lo.ForEach(m.Indexes, func(idx *model.Index, _ int) {
		ret.W(idx.SQL())
	})

	sqlHistory(ret, m, modules, now)
	ret.W("-- {%% endfunc %%}")
	return ret, nil
}
