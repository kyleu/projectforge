package sql

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func Migration(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, m.Package+".sql")
	if m.IsRevision() {
		drop, err := sqlDrop(m)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(drop)
		cr, err := sqlCreateRevision(m)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(cr)
	} else {
		drop, err := sqlDrop(m)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(drop)
		g.AddBlocks(sqlCreate(m))
	}
	return g.Render()
}

func sqlDrop(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("SQLDrop", "sql")
	ret.W("-- {%% func " + m.Proper() + "Drop() %%}")
	if m.IsRevision() {
		ret.W("drop table if exists %q;", fmt.Sprintf("%s_%s", m.Name, m.HistoryColumn().Name))
	}
	ret.W("drop table if exists %q;", m.Name)
	ret.W("-- {%% endfunc %%}")
	return ret, nil
}

func sqlCreate(m *model.Model) *golang.Block {
	ret := golang.NewBlock("SQLCreate", "sql")
	ret.W("-- {%% func " + m.Proper() + "Create() %%}")
	ret.W("create table if not exists %q (", m.Name)
	for _, col := range m.Columns {
		ret.W("  %q %s,", col.Name, col.ToSQLType())
	}
	ret.W("  primary key (%s)", strings.Join(m.PKs().NamesQuoted(), ", "))
	ret.W(");")

	if pks := m.PKs(); len(pks) > 1 {
		for _, pk := range pks {
			addIndex(ret, m.Name, pk.Name)
		}
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}

func addIndex(ret *golang.Block, tbl string, names ...string) {
	name := fmt.Sprintf("%s__%s_idx", tbl, strings.Join(names, "_"))
	msg := "create index if not exists %q on %q(%q);"
	ret.W("")
	ret.W(msg, name, tbl, strings.Join(names, ", "))
}
