package model

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	HistoryType  = "history"
	RevisionType = "revision"
)

type HistoryMap struct {
	Col   *Column `json:"col"`
	Const Columns `json:"const,omitempty"`
	Var   Columns `json:"var,omitempty"`
	Err   error   `json:"-"`
}

func (m *Model) HistoryColumns(coreColumns bool) *HistoryMap {
	if coreColumns && m.historyMapDB != nil {
		return m.historyMapDB
	}
	if !coreColumns && m.historyMap != nil {
		return m.historyMap
	}
	var pk Columns
	var c Columns
	var v Columns

	revCol, err := m.Columns.OneWithTag(RevisionType)
	if err != nil {
		return &HistoryMap{Err: err}
	}
	if revCol.Nullable {
		return &HistoryMap{Err: errors.Errorf("revision column [%s] must not be nullable", revCol.Name)}
	}
	revCurrentCol := &Column{
		Name:       "current_" + revCol.Name,
		Type:       revCol.Type,
		SQLDefault: revCol.SQLDefault,
		Tags:       []string{"current_revision"},
	}

	lo.ForEach(m.Columns, func(col *Column, _ int) {
		if col.Name == revCol.Name {
			return
		}
		switch {
		case col.PK:
			pk = append(pk, col)
		case col.HasTag("const") || col.HasTag("updated") || col.HasTag("deleted"):
			c = append(c, col)
		default:
			v = append(v, col)
		}
	})

	cRet := Columns{}
	cRet = append(cRet, pk...)
	if coreColumns {
		cRet = append(cRet, revCurrentCol)
	} else {
		cRet = append(cRet, revCol)
	}
	cRet = append(cRet, c...)

	vRet := Columns{}
	lo.ForEach(pk, func(col *Column, _ int) {
		if coreColumns {
			col = &Column{
				Name:       fmt.Sprintf("%s_%s", m.Name, col.Name),
				Type:       col.Type,
				PK:         col.PK,
				Nullable:   col.Nullable,
				Search:     col.Search,
				SQLDefault: col.SQLDefault,
				Tags:       nil,
			}
		}
		vRet = append(vRet, col)
	})
	vRet = append(vRet, revCol)
	vRet = append(vRet, v...)

	if coreColumns {
		m.historyMapDB = &HistoryMap{Col: revCol, Const: cRet, Var: vRet}
		return m.historyMapDB
	}
	m.historyMap = &HistoryMap{Col: revCol, Const: cRet, Var: vRet}
	return m.historyMap
}

func (m *Model) HistoryColumn() *Column {
	return m.HistoryColumns(true).Col
}
