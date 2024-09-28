package svc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const tableClause = "tableQuoted"

func ServiceGet(m *model.Model, args *model.Args, linebreak string) (*file.File, error) {
	dbRef := args.DBRef()
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "serviceget")
	lo.ForEach(helper.ImportsForTypes("go", "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	if len(m.PKs()) > 1 {
		g.AddImport(helper.ImpFmt)
	}
	g.AddImport(helper.ImpAppUtil, helper.ImpContext, helper.ImpErrors, helper.ImpSQLx, helper.ImpFilter, helper.ImpAppDatabase, helper.ImpLo)
	imps, err := helper.SpecialImports(m.IndexedColumns(true), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("serviceget")...)
	pkLen := len(m.PKs())
	if pkLen > 0 {
		gbpk, err := serviceGetByPK(m, dbRef, args.Enums, args.Database)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(gbpk)
		if pkLen == 1 {
			x, err := serviceGetMultipleSingleCol(m, "GetMultiple", m.PKs()[0], dbRef, args.Enums)
			if err != nil {
				return nil, err
			}
			g.AddBlocks(x)
		} else {
			g.AddBlocks(serviceGetMultipleManyCols(m, "GetMultiple", m.PKs(), dbRef))
		}
	}
	getBys := map[string]model.Columns{}
	titles := map[string]string{}
	var doExtra []string
	add := func(c *model.Column) {
		getBys[c.Name] = model.Columns{c}
		if !slices.Contains(doExtra, c.Name) {
			doExtra = append(doExtra, c.Name)
		}
	}
	if pkLen > 1 {
		lo.ForEach(m.PKs(), func(pkCol *model.Column, _ int) {
			add(pkCol)
		})
	}
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		cols := rel.SrcColumns(m)
		colStr := strings.Join(cols.Names(), ",")
		getBys[colStr] = cols
	})
	lo.ForEach(m.IndexedColumns(false), func(col *model.Column, _ int) {
		add(col)
	})
	lo.ForEach(m.Columns, func(col *model.Column, _ int) {
		lo.ForEach(col.Tags, func(tag string, _ int) {
			if strings.HasPrefix(tag, "fn:") {
				fn := strings.TrimPrefix(tag, "fn:")
				getBys[fn] = append(getBys[fn], col)
				titles[fn] = fn
			}
		})
	})
	for _, key := range util.ArraySorted(lo.Keys(getBys)) {
		err := writeGetBy(key, getBys[key], doExtra, titles[key], dbRef, m, args, g)
		if err != nil {
			return nil, err
		}
	}
	g.AddBlocks(serviceRandom(m, args.Database))
	return g.Render(linebreak)
}

func serviceGet(key string, m *model.Model, cols model.Columns, dbRef string, enums enum.Enums) (*golang.Block, error) {
	if key == "" {
		key = helper.TextGetBy + cols.Smushed()
	}
	ret := golang.NewBlock(key, "func")
	args, err := cols.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s%s, logger util.Logger) (*%s, error) {"
	ret.WF(msg, key, args, getSuffix(m), m.Proper())
	if slices.Equal(m.PKs().Names(), cols.Names()) {
		ret.W("\twc := defaultWC(0)")
	} else {
		wc := make([]string, 0, len(cols))
		lo.ForEach(cols, func(col *model.Column, idx int) {
			wc = append(wc, fmt.Sprintf("%q = $%d", col.SQL(), idx+1))
		})
		ret.WF("\twc := %q", strings.Join(wc, " and "))
	}
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := &row{}")
	ret.WF("\tq := database.SQLSelectSimple(columnsString, %s, s.db.Type, wc)", tableClause)
	ret.WF("\terr := s.%s.Get(ctx, ret, q, tx, logger, %s)", dbRef, strings.Join(cols.CamelNames(), ", "))
	ret.W("\tif err != nil {")
	sj := strings.Join(cols.CamelNames(), ", ")
	decls := make([]string, 0, len(cols))
	lo.ForEach(cols, func(c *model.Column, _ int) {
		decls = append(decls, c.Camel()+declSubscript)
	})
	ret.WF("\t\treturn nil, errors.Wrapf(err, \"unable to get %s by %s\", %s)", m.Camel(), strings.Join(decls, ", "), sj)
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.Proper())
	ret.W("}")
	return ret, nil
}

func serviceRandom(m *model.Model, database string) *golang.Block {
	ret := golang.NewBlock("Random", "func")
	ret.WF("func (s *Service) Random(ctx context.Context, tx *sqlx.Tx, logger util.Logger) (*%s, error) {", m.Proper())
	ret.W("\tret := &row{}")
	rnd := "random()"
	if database == util.DatabaseSQLServer {
		rnd = "newid()"
	}
	ret.WF("\tq := database.SQLSelect(columnsString, tableQuoted, \"\", %q, 1, 0, s.db.Type)", rnd)
	ret.W("\terr := s.db.Get(ctx, ret, q, tx, logger)")
	ret.W("\tif err != nil {")
	ret.WF("\t\treturn nil, errors.Wrap(err, \"unable to get random %s\")", m.TitlePluralLower())
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.Proper())
	ret.W("}")
	return ret
}
