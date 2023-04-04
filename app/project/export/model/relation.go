package model

import (
	"strings"

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
	for _, s := range r.Src {
		c := src.Columns.Get(s)
		url += "+`/`+" + c.ToGoString(prefix)
	}
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
	for _, x := range r {
		if x.ContainsSource(colName) {
			return true
		}
	}
	return false
}

func colsFor(cols []string, m *Model) Columns {
	ret := make(Columns, 0, len(cols))
	for _, x := range cols {
		ret = append(ret, m.Columns.Get(x))
	}
	return ret
}
