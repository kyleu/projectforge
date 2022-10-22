package sql

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func Migration(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, m.Name+".sql")
	if m.IsRevision() {
		drop, err := sqlDrop(m)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(drop)
		cr, err := sqlCreateRevision(m, args.Modules)
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
		sc, err := sqlCreate(m, args.Modules)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(sc)
	}
	return g.Render(addHeader)
}

func sqlDrop(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("SQLDrop", "sql")
	ret.W("-- {%% func " + m.Proper() + "Drop() %%}")
	if m.IsHistory() {
		ret.W("drop table if exists %q;", fmt.Sprintf("%s_history", m.Name))
	}
	if m.IsRevision() {
		ret.W("drop table if exists %q;", fmt.Sprintf("%s_%s", m.Name, m.HistoryColumn().Name))
	}
	ret.W("drop table if exists %q;", m.Name)
	ret.W("-- {%% endfunc %%}")
	return ret, nil
}

func sqlCreate(m *model.Model, modules []string) (*golang.Block, error) {
	ret := golang.NewBlock("SQLCreate", "sql")
	ret.W("-- {%% func " + m.Proper() + "Create() %%}")
	ret.W("create table if not exists %q (", m.Name)
	for _, col := range m.Columns {
		st, err := col.ToSQLType()
		if err != nil {
			return nil, err
		}
		ret.W("  %q %s,", col.Name, st)
	}
	sqlRelations(ret, m)
	for _, col := range m.Columns {
		if col.HasTag("unique") {
			ret.W("  unique (%q),", col.Name)
		}
	}
	ret.W("  primary key (%s)", strings.Join(m.PKs().NamesQuoted(), ", "))
	ret.W(");")

	if pks := m.PKs(); len(pks) > 1 {
		for _, pk := range pks {
			addIndex(ret, m.Name, pk.Name)
		}
	}
	for _, rel := range m.Relations {
		cols := rel.SrcColumns(m)
		if len(cols) == 1 && cols[0].PK {
			continue
		}
		addIndex(ret, m.Name, cols.Names()...)
	}
	sqlHistory(ret, m, modules)
	ret.W("-- {%% endfunc %%}")
	return ret, nil
}

func addIndex(ret *golang.Block, tbl string, names ...string) {
	name := fmt.Sprintf("%s__%s_idx", tbl, strings.Join(names, "_"))
	quoted := make([]string, 0, len(names))
	for _, n := range names {
		quoted = append(quoted, fmt.Sprintf("%q", n))
	}
	ret.W("")
	ret.W("create index if not exists %q on %q (%s);", name, tbl, strings.Join(quoted, ", "))
}
