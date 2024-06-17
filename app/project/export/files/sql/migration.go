package sql

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Migration(m *model.Model, args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, m.Name+util.ExtSQL)
	drop, err := sqlDrop(m, args.Database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(drop)
	sc, err := sqlCreate(m, args.Models, args.Database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(sc)
	return g.Render(linebreak)
}

func sqlDrop(m *model.Model, database string) (*golang.Block, error) {
	ret := golang.NewBlock("SQLDrop", "sql")
	ret.W(sqlFunc(m.Proper() + helper.TextDrop))
	if database == util.DatabaseSQLServer {
		ret.W("if exists (select * from sysobjects where name='%s' and xtype='U')", m.Table())
		ret.W("drop table %q;", m.Table())
	} else {
		ret.W("drop table if exists %q;", m.Table())
	}
	ret.W(sqlEnd())
	return ret, nil
}

func sqlCreate(m *model.Model, models model.Models, database string) (*golang.Block, error) {
	ret := golang.NewBlock("SQLCreate", "sql")
	ret.W(sqlFunc(m.Proper() + "Create"))
	if database == util.DatabaseSQLServer {
		ret.W("if not exists (select * from sysobjects where name='%s' and xtype='U')", m.Table())
		ret.W("create table %q (", m.Table())
	} else {
		ret.W("create table if not exists %q (", m.Table())
	}
	for _, col := range m.Columns.NotDerived() {
		st, err := col.ToSQLType(database)
		if err != nil {
			return nil, err
		}
		ret.W("  %q %s,", col.SQL(), st)
	}
	sqlRelations(ret, m, models)
	lo.ForEach(m.Columns, func(col *model.Column, _ int) {
		if col.HasTag("unique") {
			ret.W("  unique (%q),", col.SQL())
		}
	})
	ret.W("  primary key (%s)", strings.Join(m.PKs().NamesQuoted(), ", "))
	ret.W(");")

	pks := m.PKs()

	// var indexes [][]string
	lo.ForEach(m.Columns, func(col *model.Column, _ int) {
		if (col.PK && len(pks) > 1) || col.Indexed {
			addIndex(database, ret, m.Table(), col.SQL())
		}
	})
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		cols := rel.SrcColumns(m)
		if slices.Equal(cols.Names(), m.PKs().Names()) {
			return
		}
		for _, c := range cols {
			if !(c.PK || c.Indexed) {
				addIndex(database, ret, m.Table(), cols.Names()...)
				break
			}
		}
	})
	lo.ForEach(m.Indexes, func(idx *model.Index, _ int) {
		ret.W(idx.SQL())
	})
	ret.W(sqlEnd())
	return ret, nil
}

func addIndex(database string, ret *golang.Block, tbl string, names ...string) {
	name := fmt.Sprintf("%s__%s_idx", tbl, strings.Join(names, "_"))
	quoted := lo.Map(names, func(n string, _ int) string {
		return fmt.Sprintf("%q", n)
	})
	ret.WB()
	if database == util.DatabaseSQLServer {
		ret.W("if not exists (select * from sys.indexes where name='%s' and object_id=object_id('%s'))", tbl, name)
		ret.W("create index %q on %q (%s);", name, tbl, strings.Join(quoted, ", "))
	} else {
		ret.W("create index if not exists %q on %q (%s);", name, tbl, strings.Join(quoted, ", "))
	}
}
