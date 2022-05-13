package database

import (
	"context"

	"{{{ .Package }}}/queries/schema"
	"{{{ .Package }}}/app/util"
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

type sizeDTO struct {
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

func (s *sizeDTO) ToSize() *TableSize {
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

type sizeDTOs []*sizeDTO

func (ts sizeDTOs) ToSizes() TableSizes {
	ret := make(TableSizes, 0, len(ts))
	for _, t := range ts {
		ret = append(ret, t.ToSize())
	}
	return ret
}

func (s *Service) Sizes(ctx context.Context, logger util.Logger) (TableSizes, error) {
	q := schema.SizeInfo()
	ret := sizeDTOs{}
	err := s.Select(ctx, &ret, q, nil, logger)
	if err != nil {
		return nil, err
	}
	return ret.ToSizes(), nil
}
