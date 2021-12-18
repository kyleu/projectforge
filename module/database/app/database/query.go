package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Query(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) (*sqlx.Rows, error) {
	op := "query"
	now, ctx, span := s.newSpan(ctx, "db:"+op, q)
	var ret *sqlx.Rows
	var err error
	defer s.complete(q, op, span, now, err)
	s.logQuery("running raw query", q, values)
	if tx == nil {
		ret, err = s.db.QueryxContext(ctx, q, values...)
		return ret, err
	}
	ret, err = tx.QueryxContext(ctx, q, values...)
	return ret, err
}

func (s *Service) QueryRows(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) ([]util.ValueMap, error) {
	rows, err := s.Query(ctx, q, tx, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var ret []util.ValueMap
	for rows.Next() {
		x := map[string]interface{}{}
		err = rows.MapScan(x)
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan output row")
		}
		ret = append(ret, x)
	}

	return ret, nil
}

func (s *Service) Query2DArray(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) ([][]interface{}, error) {
	var err error
	var slice []interface{}

	rows, err := s.Query(ctx, q, tx, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var ret [][]interface{}
	for rows.Next() {
		slice, err = rows.SliceScan()
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan output row")
		}
		ret = append(ret, slice)
	}

	return ret, nil
}

func (s *Service) QueryKVMap(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) (util.ValueMap, error) {
	msg := `must provide two columns, a string key named "k" and value "v"`
	maps, err := s.QueryRows(ctx, q, tx, values...)
	if err != nil {
		return nil, err
	}
	ret := make(util.ValueMap, len(maps))
	for _, m := range maps {
		kx, ok := m["k"]
		if !ok {
			return nil, errors.New(msg)
		}
		k, ok := kx.(string)
		if !ok {
			return nil, errors.New(msg)
		}
		v, ok := m["v"]
		if !ok {
			return nil, errors.New(msg)
		}
		ret[k] = v
	}
	return ret, nil
}

func (s *Service) QuerySingleRow(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) (util.ValueMap, error) {
	rows, err := s.QueryRows(ctx, q, tx, values...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("no rows returned from query")
	}
	if len(rows) < 1 {
		return nil, errors.New("more than one row returned from query")
	}
	return rows[0], nil
}

func (s *Service) Select(ctx context.Context, dest interface{}, q string, tx *sqlx.Tx, values ...interface{}) error {
	op := "select"
	now, ctx, span := s.newSpan(ctx, "db:"+op, q)
	var err error
	defer s.complete(q, op, span, now, err)
	s.logQuery(fmt.Sprintf("selecting rows of type [%T]", dest), q, values)
	if tx == nil {
		err = s.db.SelectContext(ctx, dest, q, values...)
		return err
	}
	err = tx.SelectContext(ctx, dest, q, values...)
	return err
}

func (s *Service) Get(ctx context.Context, dto interface{}, q string, tx *sqlx.Tx, values ...interface{}) error {
	op := "get"
	now, ctx, span := s.newSpan(ctx, "db:"+op, q)
	var err error
	defer s.complete(q, op, span, now, err)
	s.logQuery(fmt.Sprintf("getting single row of type [%T]", dto), q, values)
	if tx == nil {
		return s.db.GetContext(ctx, dto, q, values...)
	}
	return tx.GetContext(ctx, dto, q, values...)
}

type singleIntResult struct {
	X *int64 `db:"x"`
}

func (s *Service) SingleInt(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) (int64, error) {
	op := "single-int"
	now, ctx, span := s.newSpan(ctx, "db:"+op, q)
	var err error
	defer s.complete(q, op, span, now, err)
	x := &singleIntResult{}
	err = s.Get(ctx, x, q, tx, values...)
	if err != nil {
		return -1, err
	}
	if x.X == nil {
		return 0, nil
	}
	return *x.X, nil
}
