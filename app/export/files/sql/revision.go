package sql

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func sqlDropRevision(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("SQLDrop", "sql")
	ret.W("-- {%% func " + m.Proper() + "Drop() %%}")
	if m.IsRevision() {
		ret.W("drop table if exists %q;", fmt.Sprintf("%s_%s", m.Package, m.HistoryColumn().Name))
	}
	ret.W("drop table if exists %q;", m.Package)
	ret.W("-- {%% endfunc %%}")
	return ret, nil
}

func sqlCreateRevision(m *model.Model) (*golang.Block, error) {
	pks := m.PKs()
	hc := m.HistoryColumns(true)

	ret := golang.NewBlock("SQLCreateRev", "sql")
	ret.W("-- {%% func " + m.Proper() + "Create() %%}")

	// core
	ret.W("create table if not exists %q (", m.Package)
	for _, col := range hc.Const {
		ret.W("  %q %s,", col.Name, col.ToSQLType())
	}
	ret.W("  primary key (%s)", strings.Join(pks.NamesQuoted(), ", "))
	ret.W(");")

	if len(pks) > 1 {
		for _, pk := range pks {
			ret.W("")
			ret.W("create index if not exists %q on %q(%q);", fmt.Sprintf("%s__%s_idx", m.Package, pk.Name), m.Package, pk.Name)
		}
	}

	// revision
	revTblName := m.Package + "_" + hc.Col.Name
	ret.W("")
	ret.W("create table if not exists %q (", revTblName)
	for _, col := range hc.Var {
		ret.W("  %q %s,", col.Name, col.ToSQLType())
	}

	revPKs := hc.Var.PKs()
	revPKsWithRev := append(model.Columns{}, revPKs...)
	revPKsWithRev = append(revPKsWithRev, hc.Col)

	bareRefs := strings.Join(revPKs.NamesQuoted(), ", ")
	ret.W("  foreign key (%s) references %s(%s),", bareRefs, m.Name, strings.Join(pks.NamesQuoted(), ", "))

	ret.W("  primary key (%s)", strings.Join(revPKsWithRev.NamesQuoted(), ", "))
	ret.W(");")

	msg := "create index if not exists \"%s__%s_idx\" on %q(%s);"
	ret.W(msg, revTblName, strings.Join(revPKs.Names(), "_"), revTblName, strings.Join(revPKs.NamesQuoted(), ", "))

	for _, pk := range revPKsWithRev {
		if !pk.HasTag(model.RevisionType) {
			addIndex(ret, revTblName, pk.Name)
		}
	}

	ret.W("-- {%% endfunc %%}")
	return ret, nil
}
