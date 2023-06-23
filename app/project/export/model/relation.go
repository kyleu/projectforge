package model

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Relation struct {
	Name  string   `json:"name"`
	Src   []string `json:"src"`
	Table string   `json:"table"`
	Tgt   []string `json:"tgt"`
}

func (r *Relation) SrcColumns(m *Model) Columns {
	return colsFor(r.Src, m)
}

func (r *Relation) SrcQuoted() any {
	return strings.Join(util.StringArrayQuoted(r.Src), ", ")
}

func (r *Relation) TgtColumns(m *Model) Columns {
	return colsFor(r.Tgt, m)
}

func (r *Relation) TgtQuoted() any {
	return strings.Join(util.StringArrayQuoted(r.Tgt), ", ")
}

func (r *Relation) WebPath(src *Model, tgt *Model, prefix string) any {
	url := "`/" + tgt.Route() + "`"
	lo.ForEach(r.Src, func(s string, _ int) {
		c := src.Columns.Get(s)
		url += "+`/`+" + c.ToGoString(prefix)
	})
	return url
}

func (r *Relation) Reverse(name string) *Relation {
	return &Relation{Name: r.Name, Src: r.Tgt, Table: name, Tgt: r.Src}
}

func (r *Relation) ContainsSource(colName string) bool {
	return len(r.Src) == 1 && r.Src[0] == colName
}

type Relations []*Relation

func (r Relations) ContainsSource(colName string) bool {
	return lo.ContainsBy(r, func(x *Relation) bool {
		return x.ContainsSource(colName)
	})
}

func colsFor(cols []string, m *Model) Columns {
	return lo.Map(cols, func(x string, _ int) *Column {
		return m.Columns.Get(x)
	})
}
