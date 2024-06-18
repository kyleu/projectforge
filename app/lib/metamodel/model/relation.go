package model

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var RelationFieldDescs = util.FieldDescs{
	{Key: "name", Title: "Name", Description: "The name of the relation"},
	{Key: "src", Title: "Source", Description: "The source columns of the relation"},
	{Key: "table", Title: "Table", Description: "The target table of the relation"},
	{Key: "tgt", Title: "Target", Description: "The target columns of the relation"},
}

type Relation struct {
	Name  string   `json:"name"`
	Src   []string `json:"src"`
	Table string   `json:"table"`
	Tgt   []string `json:"tgt"`
}

func (r *Relation) SrcColumns(m *Model) Columns {
	return colsFor(r.Src, m)
}

func (r *Relation) SrcQuoted() string {
	return strings.Join(util.StringArrayQuoted(r.Src), ", ")
}

func (r *Relation) TgtColumns(m *Model) Columns {
	return colsFor(r.Tgt, m)
}

func (r *Relation) TgtQuoted() string {
	return strings.Join(util.StringArrayQuoted(r.Tgt), ", ")
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

func (r Relations) Get(name string) *Relation {
	return lo.FindOrElse(r, nil, func(x *Relation) bool {
		return x.Name == name
	})
}

func colsFor(cols []string, m *Model) Columns {
	return lo.Map(cols, func(x string, _ int) *Column {
		return m.Columns.Get(x)
	})
}
