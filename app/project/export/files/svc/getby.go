package svc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const declSubscript = " [%%v]"

func writeGetBy(key string, cols model.Columns, doExtra []string, name string, dbRef string, m *model.Model, args *metamodel.Args, g *golang.File) error {
	if name == "" {
		name = helper.TextGetBy + strings.Join(cols.ProperNames(), "")
	}
	lo.ForEach(helper.ImportsForTypes("go", "", cols.Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	returnMultiple := lo.NoneBy(cols, func(x *model.Column) bool {
		return x.HasTag("unique")
	})
	sb, err := serviceGetBy(name, m, cols, returnMultiple, dbRef, args.Enums, args.Database, g)
	if err != nil {
		return err
	}
	g.AddBlocks(sb)
	if slices.Contains(doExtra, key) {
		if len(cols) == 1 {
			n := cols[0].ProperPlural()
			if cols[0].ProperPlural() == cols[0].Proper() {
				n += "Set"
			}
			pb, err := serviceGetMultipleSingleCol(m, helper.TextGetBy+n, cols[0], dbRef, args.Enums)
			if err != nil {
				return err
			}
			g.AddBlocks(pb)
		}
	}
	return nil
}

func serviceGetByPK(m *model.Model, dbRef string, enums enum.Enums, database string, g *golang.File) (*golang.Block, error) {
	return serviceGetBy("Get", m, m.PKs(), false, dbRef, enums, database, g)
}

func serviceGetBy(
	key string, m *model.Model, cols model.Columns, returnMultiple bool, dbRef string, enums enum.Enums, database string, g *golang.File,
) (*golang.Block, error) {
	if returnMultiple {
		return serviceGetByCols(key, m, cols, dbRef, enums, database)
	}
	return serviceGet(key, g, m, cols, dbRef, enums)
}

func serviceGetByCols(key string, m *model.Model, cols model.Columns, dbRef string, enums enum.Enums, database string) (*golang.Block, error) {
	if key == "" {
		key = helper.TextGetBy + cols.Smushed()
	}
	ret := golang.NewBlock(key, "func")
	args, err := cols.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s, params *filter.Params%s, logger util.Logger) (%s, error) {"
	msg = fmt.Sprintf(msg, key, args, getSuffix(m), m.ProperPlural())
	ret.W(msg)
	ret.W("\tparams = filters(params)")
	var placeholder string
	if database == util.DatabaseSQLServer {
		placeholder = "@"
	}
	ret.WF("\twc := %q", cols.WhereClause(0, placeholder))
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.WF("\tq := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset, s.db.Type)", tableClause)
	ret.W("\tret := rows{}")
	ret.WF("\terr := s.%s.Select(ctx, &ret, q, tx, logger, %s)", dbRef, strings.Join(cols.CamelNames(), ", "))
	ret.W("\tif err != nil {")
	sj := strings.Join(cols.CamelNames(), ", ")
	decls := make([]string, 0, len(cols))
	lo.ForEach(cols, func(c *model.Column, _ int) {
		decls = append(decls, c.Camel()+declSubscript)
	})
	ret.WF("\t\treturn nil, errors.Wrapf(err, \"unable to get %s by %s\", %s)", m.TitlePlural(), strings.Join(decls, ", "), sj)
	ret.W("\t}")
	ret.WF("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret, nil
}
