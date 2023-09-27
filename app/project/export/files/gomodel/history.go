package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func History(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, m.Camel()+"history")
	lo.ForEach(helper.ImportsForTypes("go", "", m.Columns.Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpJSON, helper.ImpUUID, helper.ImpAppUtil, helper.ImpLo)
	mh, err := modelHistory(m, args.Enums)
	if err != nil {
		return nil, err
	}
	row, err := modelHistoryRow(m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mh, modelHistoryToData(m), modelHistories(m), row, modelHistoryRowToHistory(m), modelHistoryRows(m))
	return g.Render(addHeader, linebreak)
}

func modelHistory(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"History", "struct")
	ret.W("type History struct {")
	maxCount := m.PKs().MaxCamelLength() + len(m.Camel())
	if maxCount < 7 {
		maxCount = 7
	}
	tmax := 13
	ret.W("\t%s %s `json:\"id\"`", util.StringPad("ID", maxCount), util.StringPad("uuid.UUID", tmax))
	for _, pk := range m.PKs() {
		gt, err := pk.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, tmax)
		ret.W("\t%s %s `json:\"%s%s\"`", util.StringPad(m.Proper()+pk.Proper(), maxCount), goType, m.Camel(), pk.Proper())
	}
	ret.W("\t%s util.ValueMap `json:\"o,omitempty\"`", util.StringPad("Old", maxCount))
	ret.W("\t%s util.ValueMap `json:\"n,omitempty\"`", util.StringPad("New", maxCount))
	ret.W("\t%s util.Diffs    `json:\"c,omitempty\"`", util.StringPad("Changes", maxCount))
	ret.W("\t%s time.Time     `json:\"created\"`", util.StringPad("Created", maxCount))
	ret.W("}")
	return ret, nil
}

func modelHistoryToData(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"HistoryToData", "func")
	ret.W("func (h *History) ToData() []any {")
	ret.W("\treturn []any{")
	ret.W("\t\th.ID,")
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		ret.W("\t\th.%s%s,", m.Proper(), pk.Proper())
	})
	ret.W("\t\tutil.ToJSONBytes(h.Old, true),")
	ret.W("\t\tutil.ToJSONBytes(h.New, true),")
	ret.W("\t\tutil.ToJSONBytes(h.Changes, true),")
	ret.W("\t\th.Created,")
	ret.W("\t}")
	ret.W("}")
	return ret
}

func modelHistories(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Histories", "struct")
	ret.W("type Histories []*History")
	return ret
}

func modelHistoryRow(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"HistoryRow", "struct")
	ret.W("type historyRow struct {")
	maxCount := m.PKs().MaxCamelLength() + len(m.Camel())
	if maxCount < 7 {
		maxCount = 7
	}
	tmax := 15
	ret.W("\t%s %s `db:\"id\"`", util.StringPad("ID", maxCount), util.StringPad("uuid.UUID", tmax))
	for _, pk := range m.PKs() {
		gt, err := pk.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, tmax)
		ret.W("\t%s %s `db:\"%s_%s\"`", util.StringPad(m.Proper()+pk.Proper(), maxCount), goType, m.Name, pk.Name)
	}
	ret.W("\t%s json.RawMessage `db:\"o\"`", util.StringPad("Old", maxCount))
	ret.W("\t%s json.RawMessage `db:\"n\"`", util.StringPad("New", maxCount))
	ret.W("\t%s json.RawMessage `db:\"c\"`", util.StringPad("Changes", maxCount))
	ret.W("\t%s time.Time       `db:\"created\"`", util.StringPad("Created", maxCount))
	ret.W("}")
	return ret, nil
}

func modelHistoryRowToHistory(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"HistoryRowToHistory", "func")
	ret.W("func (h *historyRow) ToHistory() *History {")
	ret.W("\to := util.ValueMap{}")
	ret.W("\t_ = util.FromJSON(h.Old, &o)")
	ret.W("\tn := util.ValueMap{}")
	ret.W("\t_ = util.FromJSON(h.New, &n)")
	ret.W("\tc := util.Diffs{}")
	ret.W("\t_ = util.FromJSON(h.Changes, &c)")
	pkCalls := make([]string, 0, len(m.PKs()))
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		pkCalls = append(pkCalls, fmt.Sprintf("%s%s: h.%s%s", m.Proper(), pk.Proper(), m.Proper(), pk.Proper()))
	})
	ret.W("\treturn &History{ID: h.ID, %s, Old: o, New: n, Changes: c, Created: h.Created}", strings.Join(pkCalls, ", "))
	ret.W("}")
	return ret
}

func modelHistoryRows(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"HistoryRows", "func")
	ret.W("type historyRows []*historyRow")
	ret.WB()
	ret.W("func (h historyRows) ToHistories() Histories {")
	ret.W("\treturn lo.Map(h, func(x *historyRow, _ int) *History {")
	ret.W("\t\treturn x.ToHistory()")
	ret.W("\t})")
	ret.W("}")
	return ret
}
