package svc

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const (
	serviceAssignmentToken = ":="
	chunkedArgs            = "ctx context.Context, tx *sqlx.Tx, chunkSize int, progress *util.Progress, logger util.Logger, models ...*%s) error {"
)

func ServiceMutate(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "servicemutate")
	lo.ForEach(helper.ImportsForTypes("go", "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil, helper.ImpContext, helper.ImpSQLx, helper.ImpAppDatabase, helper.ImpLo)
	g.AddImport(m.Imports.Supporting("servicemutate")...)

	if add, err := serviceCreate(g, m, args.Audit(m)); err == nil {
		g.AddBlocks(add, serviceCreateChunked(m))
	} else {
		return nil, err
	}
	if upd, err := serviceUpdate(g, m, args.Audit(m), args.Database); err == nil {
		g.AddBlocks(upd)
	} else {
		return nil, err
	}
	if save, err := serviceSave(g, m, args.Audit(m)); err == nil {
		g.AddBlocks(save, serviceSaveChunked(m))
	} else {
		return nil, err
	}
	if m.IsSoftDelete() {
		g.AddImport(helper.ImpAppUtil)
		sdel, err := serviceSoftDelete(m, args)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(sdel, serviceSoftDeleteWhere(m), serviceAddDeletedClause(m))
	} else {
		del, err := serviceDelete(m, args)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(del, serviceDeleteWhere(m))
	}
	return g.Render(linebreak)
}

func serviceCreate(g *golang.File, m *model.Model, audit bool) (*golang.Block, error) {
	ret := golang.NewBlock("Create", "func")
	ret.WF("func (s *Service) Create(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if err := serviceAddCreatedUpdated(g, m, ret, false); err != nil {
		return nil, err
	}
	ret.W("\tq := database.SQLInsert(tableQuoted, columnsQuoted, len(models), s.db.Type)")

	if audit {
		ret.W("\tvals := make([]any, 0, len(models)*len(columnsQuoted))")
		ret.W("\tfor _, arg := range models {")
		msg := "\t\t_, _, err := s.audit.ApplyObjSimple(ctx, \"%s.create\", \"created new %s\", nil, arg, %q, nil, logger)"
		ret.WF(msg, m.Proper(), m.TitleLower(), m.Name)
		ret.WE(2)
		ret.W("\t\tvals = append(vals, arg.ToData()...)")
		ret.W("\t}")
	} else {
		ret.WF("\tvals := lo.FlatMap(models, func(arg *%s, _ int) []any {", m.Proper())
		ret.W("\t\treturn arg.ToData()")
		ret.W("\t})")
	}
	if m.HasTag("events") {
		ret.W("\terr := s.db.Insert(ctx, q, tx, logger, vals...)")
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Create(ctx, tx, logger, models...); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [create] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
		ret.W("\treturn err")
	} else {
		ret.W("\treturn s.db.Insert(ctx, q, tx, logger, vals...)")
	}

	ret.W("}")
	return ret, nil
}

func serviceCreateChunked(m *model.Model) *golang.Block {
	ret := golang.NewBlock("CreateChunked", "func")
	ret.WF("func (s *Service) CreateChunked("+chunkedArgs, m.Proper())
	ret.W("\tfor idx, chunk := range lo.Chunk(models, chunkSize) {")
	ret.W("\t\tif logger != nil {")
	ret.W("\t\t\tcount := ((idx + 1) * chunkSize) - 1")
	ret.W("\t\t\tif len(models) < count {")
	ret.W("\t\t\t\tcount = len(models)")
	ret.W("\t\t\t}")
	ret.WF("\t\t\tlogger.Infof(\"creating %s [%%%%d-%%%%d]\", idx*chunkSize, count)", m.TitlePluralLower())
	ret.W("\t\t}")
	ret.W("\t\tif err := s.Create(ctx, tx, logger, chunk...); err != nil {")
	ret.W("\t\t\treturn err")
	ret.W("\t\t}")
	ret.W("\t\tprogress.Increment(len(chunk), logger)")
	ret.W("\t}")
	if m.HasTag("events") {
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Create(ctx, tx, logger, models...); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [createChunked] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
	}
	ret.W("\treturn nil")
	ret.W("}")
	return ret
}

func serviceUpdate(g *golang.File, m *model.Model, audit bool, database string) (*golang.Block, error) {
	ret := golang.NewBlock("Update", "func")
	ret.WF("func (s *Service) Update(ctx context.Context, tx *sqlx.Tx, model *%s, logger util.Logger) error {", m.Proper())

	cols := m.Columns.NotDerived()
	if cc := cols.WithTag("created"); len(cc) > 0 || audit {
		g.AddImport(helper.ImpErrors)
		ret.WF("\tcurr, err := s.Get(ctx, tx, %s%s, logger)", m.PKs().ToRefs(helper.TextModelPrefix), m.SoftDeleteSuffix())
		ret.W("\tif err != nil {")
		ret.WF("\t\treturn errors.Wrapf(err, \"can't get original %s [%%%%s]\", model.String())", m.TitleLower())
		ret.W("\t}")
		lo.ForEach(cc, func(c *model.Column, _ int) {
			ret.WF("\tmodel.%s = curr.%s", c.Proper(), c.Proper())
		})
	}

	for _, updated := range cols.WithTag("updated") {
		err := serviceSetVal(g, updated, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	pks := m.PKs()
	pkVals := make([]string, 0, len(pks))
	lo.ForEach(pks, func(pk *model.Column, _ int) {
		pkVals = append(pkVals, helper.TextModelPrefix+pk.Proper())
	})
	var placeholder string
	if database == util.DatabaseSQLServer {
		placeholder = "@"
	}
	ret.WF("\tq := database.SQLUpdate(tableQuoted, columnsQuoted, %q, s.db.Type)", pks.WhereClause(len(cols.NotDerived()), placeholder))
	ret.W("\tdata := model.ToData()")
	ret.WF("\tdata = append(data, %s)", util.StringJoin(pkVals, ", "))
	token := "="
	if len(cols.WithTag("created")) == 0 && (!audit) {
		token = serviceAssignmentToken
	}
	ret.WF("\t_, err %s s.db.Update(ctx, q, tx, 1, logger, data...)", token)
	ret.WE(1)
	if audit {
		serviceAuditApply(g, m, ret)
	}
	if m.HasTag("events") {
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Save(ctx, tx, logger, model); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [update] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
	}
	ret.W("\treturn nil")
	ret.W("}")
	return ret, nil
}

func serviceAuditApply(g *golang.File, m *model.Model, ret *golang.Block) {
	g.AddImport(helper.ImpFmt)
	ret.WF("\tmsg := fmt.Sprintf(\"updated %s [%%%%s]\", model.String())", m.Title())
	ret.WF("\t_, _, err = s.audit.ApplyObjSimple(ctx, \"%s.update\", msg, curr, model, %q, nil, logger)", m.Proper(), m.Name)
	ret.WE(1)
}

//nolint:unused
func serviceUpdateIfNeeded(g *golang.File, m *model.Model, database string) (*golang.Block, error) {
	ret := golang.NewBlock("UpdateIfNeeded", "func")
	ret.WF("func (s *Service) UpdateIfNeeded(ctx context.Context, tx *sqlx.Tx, model *%s, logger util.Logger) error {", m.Proper())

	if cc := m.Columns.WithTag("created"); len(cc) > 0 {
		g.AddImport(helper.ImpErrors)
		ret.WF("\tcurr, err := s.Get(ctx, tx, %s%s, logger)", m.PKs().ToRefs(helper.TextModelPrefix), m.SoftDeleteSuffix())
		ret.W("\tif curr == nil || err != nil {")
		ret.W("\t\treturn s.Create(ctx, tx, logger, model)")
		ret.W("\t}")
		lo.ForEach(cc, func(c *model.Column, _ int) {
			ret.WF("\tmodel.%s = curr.%s", c.Proper(), c.Proper())
		})
	}

	for _, updated := range m.Columns.WithTag("updated") {
		err := serviceSetVal(g, updated, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	pks := m.PKs()
	pkVals := make([]string, 0, len(pks))
	lo.ForEach(pks, func(pk *model.Column, _ int) {
		pkVals = append(pkVals, helper.TextModelPrefix+pk.Proper())
	})
	var placeholder string
	if database == util.DatabaseSQLServer {
		placeholder = "@"
	}
	ret.WF("\tq := database.SQLUpdate(tableQuoted, columnsQuoted, %q, s.db.Type)", pks.WhereClause(len(m.Columns), placeholder))
	ret.W("\tdata := model.ToData()")
	ret.WF("\tdata = append(data, %s)", util.StringJoin(pkVals, ", "))
	token := "="
	if len(m.Columns.WithTag("created")) == 0 {
		token = serviceAssignmentToken
	}
	ret.WF("\t_, err %s s.db.Update(ctx, q, tx, 1, logger, data...)", token)
	ret.WE(1)
	ret.W("\treturn nil")
	ret.W("}")
	return ret, nil
}

func serviceSave(g *golang.File, m *model.Model, audit bool) (*golang.Block, error) {
	ret := golang.NewBlock("Save", "func")
	ret.WF("func (s *Service) Save(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if err := serviceAddCreatedUpdated(g, m, ret, false); err != nil {
		return nil, err
	}
	q := util.StringJoin(m.PKs().SQLQuoted(), ", ")
	var pkOpt string
	ret.WF("\tq := database.SQLUpsert(tableQuoted%s, columnsQuoted, len(models), []string{%s}, columnsQuoted, s.db.Type)", pkOpt, q)
	ret.WF("\tdata := lo.FlatMap(models, func(model *%s, _ int) []any {", m.Proper())
	ret.W("\t\treturn model.ToData()")
	ret.W("\t})")
	if audit {
		if len(m.PKs()) == 1 {
			pk := m.PKs()[0]
			ret.WF("\tcurr, err := s.GetMultiple(ctx, tx, nil, logger, %s(models).%s()...)", m.ProperPlural(), pk.ProperPlural())
		} else {
			ret.WF("\tcurr, err := s.GetMultiple(ctx, tx, nil, logger, %s(models).ToPKs()...)", m.ProperPlural())
		}
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\tfor _, arg := range models {")
		ret.WF("\t\tif x := curr.Get(%s); x != nil {", m.PKs().ToRefs("arg."))
		msg := "\t\t\t_, _, err := s.audit.ApplyObjSimple(ctx, \"%s.create\", \"created new %s\", x, arg, %q, nil, logger)"
		ret.WF(msg, m.Proper(), m.Camel(), m.Name)
		ret.W("\t\t\tif err != nil {")
		ret.W("\t\t\t\treturn err")
		ret.W("\t\t\t}")
		ret.W("\t\t}")
		ret.W("\t}")
	}
	if m.HasTag("events") {
		ret.W("\terr := s.db.Insert(ctx, q, tx, logger, data...)")
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Save(ctx, tx, logger, models...); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [save] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
		ret.W("\treturn err")
	} else {
		ret.W("\treturn s.db.Insert(ctx, q, tx, logger, data...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceSaveChunked(m *model.Model) *golang.Block {
	ret := golang.NewBlock("SaveChunked", "func")
	ret.WF("func (s *Service) SaveChunked("+chunkedArgs, m.Proper())
	ret.W("\tfor idx, chunk := range lo.Chunk(models, chunkSize) {")
	ret.W("\t\tif logger != nil {")
	ret.W("\t\t\tcount := ((idx + 1) * chunkSize) - 1")
	ret.W("\t\t\tif len(models) < count {")
	ret.W("\t\t\t\tcount = len(models)")
	ret.W("\t\t\t}")
	ret.WF("\t\t\tlogger.Infof(\"saving %s [%%%%d-%%%%d]\", idx*chunkSize, count)", m.TitlePluralLower())
	ret.W("\t\t}")
	ret.W("\t\tif err := s.Save(ctx, tx, logger, chunk...); err != nil {")
	ret.W("\t\t\treturn err")
	ret.W("\t\t}")
	ret.W("\t\tprogress.Increment(len(chunk), logger)")
	ret.W("\t}")
	if m.HasTag("events") {
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Save(ctx, tx, logger, models...); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [saveChunked] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
	}
	ret.W("\treturn nil")
	ret.W("}")
	return ret
}

func serviceAddCreatedUpdated(g *golang.File, m *model.Model, ret *golang.Block, loadCurr bool) error {
	createdCols := m.Columns.WithTag("created")
	updatedCols := m.Columns.WithTag("updated")
	if len(createdCols) > 0 || len(updatedCols) > 0 {
		ret.WF("\tlo.ForEach(models, func(model *%s, _ int) {", m.Proper())
		err := serviceLoadCreated(g, ret, m, createdCols, loadCurr)
		if err != nil {
			return err
		}
		for _, updated := range updatedCols {
			err := serviceSetVal(g, updated, ret, 2)
			if err != nil {
				return err
			}
		}
		ret.W("\t})")
	}
	return nil
}

func serviceLoadCreated(g *golang.File, ret *golang.Block, m *model.Model, createdCols model.Columns, loadCurr bool) error {
	if len(createdCols) > 0 {
		if loadCurr {
			ret.WF("\t\tcurr, e := s.Get(ctx, tx, %s%s)", m.PKs().ToRefs(helper.TextModelPrefix), m.SoftDeleteSuffix())
			ret.W("\t\tif e == nil && curr != nil {")
			lo.ForEach(createdCols, func(created *model.Column, _ int) {
				ret.WF("\t\t\tmodel.%s = curr.%s", created.Proper(), created.Proper())
			})
			ret.W("\t\t} else {")
			for _, created := range createdCols {
				err := serviceSetVal(g, created, ret, 3)
				if err != nil {
					return err
				}
			}
			ret.W("\t\t}")
		} else {
			for _, created := range createdCols {
				err := serviceSetVal(g, created, ret, 2)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func serviceSetVal(g *golang.File, c *model.Column, ret *golang.Block, indent int) error {
	ind := util.StringRepeat("\t", indent)
	if c.Type.Key() == types.KeyTimestamp || c.Type.Key() == types.KeyTimestampZoned || c.Type.Key() == types.KeyDate {
		if c.Nullable {
			g.AddImport(helper.ImpAppUtil)
			ret.WF(ind+"model.%s = util.TimeCurrentP()", c.Proper())
		} else {
			g.AddImport(helper.ImpAppUtil)
			ret.WF(ind+"model.%s = util.TimeCurrent()", c.Proper())
		}
	} else {
		return errors.Errorf("unhandled type [%s]", c.Type.Key())
	}
	return nil
}
