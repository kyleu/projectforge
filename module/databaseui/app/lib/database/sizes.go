package database

import (
	"context"{{{ if .SQLite }}}{{{ if .PostgreSQL }}}{{{ else }}}
	"fmt"{{{ end }}}{{{ end }}}

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/queries/schema"
)

type NumAndString struct {
	Num    int64  `json:"num,omitempty"`
	String string `json:"string,omitempty"`
}

type TableSize struct {
	Schema string        `json:"schema"`
	Name   string        `json:"name"`
	Rows   string        `json:"rows,omitempty"`
	Total  *NumAndString `json:"total,omitempty"`
	Index  *NumAndString `json:"index,omitempty"`
	Toast  *NumAndString `json:"toast,omitempty"`
	Table  *NumAndString `json:"table,omitempty"`
}

type TableSizes []*TableSize

type sizeRow struct {
	TableSchema string `db:"table_schema"`
	TableName   string `db:"table_name"`
	RowEstimate string `db:"row_estimate"`
	Total       int64  `db:"total"`
	TotalPretty string `db:"total_pretty"`
	Index       int64  `db:"index"`
	IndexPretty string `db:"index_pretty"`
	Toast       int64  `db:"toast"`
	ToastPretty string `db:"toast_pretty"`
	Table       int64  `db:"table"`
	TablePretty string `db:"table_pretty"`
}

func (s *sizeRow) ToSize() *TableSize {
	return &TableSize{
		Schema: s.TableSchema,
		Name:   s.TableName,
		Rows:   s.RowEstimate,
		Total:  &NumAndString{Num: s.Total, String: s.TotalPretty},
		Index:  &NumAndString{Num: s.Index, String: s.IndexPretty},
		Toast:  &NumAndString{Num: s.Toast, String: s.ToastPretty},
		Table:  &NumAndString{Num: s.Table, String: s.TablePretty},
	}
}

type sizeRows []*sizeRow

func (ts sizeRows) ToSizes() TableSizes {
	return lo.Map(ts, func(t *sizeRow, _ int) *TableSize {
		return t.ToSize()
	})
}

func (s *Service) Sizes(ctx context.Context, logger util.Logger) (TableSizes, error) {
	q := schema.SizeInfo()
	ret := sizeRows{}
	err := s.Select(ctx, &ret, q, nil, logger)
	if err != nil {
		return nil, err
	}{{{ if .SQLite }}}{{{ if .PostgreSQL }}}{{{ else }}}
	for _, r := range ret {
		sizeq := fmt.Sprintf("select count(*) as x from %q", r.TableName)
		rowEstimate, err := s.SingleInt(ctx, sizeq, nil, logger, r.TableName)
		if err != nil {
			return nil, err
		}
		r.RowEstimate = fmt.Sprint(rowEstimate)
	}{{{ end }}}{{{ end }}}
	return ret.ToSizes(), nil
}
