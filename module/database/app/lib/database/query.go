package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Query(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) (*sqlx.Rows, error) {
	const op = "query"
	now, ctx, span, logger := s.newSpan(ctx, dbPrefix+op, q, logger)
	var ret *sqlx.Rows
	var err error
	defer s.complete(q, op, span, now, logger, err)
	f := s.logQuery(ctx, "running raw query", q, logger, values...)
	if tx == nil {
		ret, err = s.db.QueryxContext(ctx, q, values...)
		defer f(0, "ran raw query without transaction", err)
		return ret, err
	}
	ret, err = tx.QueryxContext(ctx, q, values...)
	defer f(0, "ran raw query with transaction", err)
	return ret, err
}

func (s *Service) QueryRows(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) ([]util.ValueMap, error) {
	rows, err := s.Query(ctx, q, tx, logger, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var ret []util.ValueMap
	for rows.Next() {
		x, err := MapScan(rows)
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan output row")
		}
		ret = append(ret, x)
	}

	return ret, nil
}

func (s *Service) Query2DArray(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) ([][]any, error) {
	var err error
	var slice []any

	rows, err := s.Query(ctx, q, tx, logger, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var ret [][]any
	for rows.Next() {
		slice, err = rows.SliceScan()
		if err != nil {
			return nil, errors.Wrap(err, "unable to scan output row")
		}
		ret = append(ret, slice)
	}

	return ret, nil
}

func (s *Service) QueryKVMap(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) (util.ValueMap, error) {
	msg := `must provide two columns, a string key named "k" and value "v"`
	maps, err := s.QueryRows(ctx, q, tx, logger, values...)
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

func (s *Service) QuerySingleRow(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) (util.ValueMap, error) {
	rows, err := s.QueryRows(ctx, q, tx, logger, values...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.Wrap(sql.ErrNoRows, "no rows returned from query")
	}
	if len(rows) > 1 {
		return nil, errors.New("more than one row returned from query")
	}
	return rows[0], nil
}

func (s *Service) Select(ctx context.Context, dest any, q string, tx *sqlx.Tx, logger util.Logger, values ...any) error {
	const op = "select"
	now, ctx, span, logger := s.newSpan(ctx, dbPrefix+op, q, logger)
	var err error
	defer s.complete(q, op, span, now, logger, err)
	f := s.logQuery(ctx, fmt.Sprintf("selecting rows of type [%T]", dest), q, logger, values...)
	if tx == nil {
		err = s.db.SelectContext(ctx, dest, q, values...)
		defer f(util.LengthAny(dest), "ran select query without transaction", err, util.ArrayFromAnyOK[any](dest)...)
		return err
	}
	err = tx.SelectContext(ctx, dest, q, values...)
	defer f(util.LengthAny(dest), "ran select query with transaction", err, util.ArrayFromAnyOK[any](dest)...)
	return err
}

func (s *Service) Get(ctx context.Context, row any, q string, tx *sqlx.Tx, logger util.Logger, values ...any) error {
	const op = "get"
	now, ctx, span, logger := s.newSpan(ctx, dbPrefix+op, q, logger)
	var err error
	defer s.complete(q, op, span, now, logger, err)
	f := s.logQuery(ctx, fmt.Sprintf("getting single row of type [%T]", row), q, logger, values...)
	if tx == nil {
		err = s.db.GetContext(ctx, row, q, values...)
		var count int
		if row != nil {
			count = 1
		}
		defer f(count, "ran [get] query without transaction", err, row)
	} else {
		err = tx.GetContext(ctx, row, q, values...)
		var count int
		if row != nil {
			count = 1
		}
		defer f(count, "ran [get] query with transaction", err, row)
	}
	return err
}

type singleIntResult struct {
	X *int64 `db:"x"`
}

func (s *Service) SingleInt(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) (int64, error) {
	const op = "single-int"
	now, ctx, span, logger := s.newSpan(ctx, dbPrefix+op, q, logger)
	var err error
	defer s.complete(q, op, span, now, logger, err)
	x := &singleIntResult{}
	err = s.Get(ctx, x, q, tx, logger, values...)
	if err != nil {
		return -1, err
	}
	if x.X == nil {
		return 0, nil
	}
	return *x.X, nil
}

type singleBoolResult struct {
	X bool `db:"x"`
}

func (s *Service) SingleBool(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) (bool, error) {
	const op = "single-bool"
	now, ctx, span, logger := s.newSpan(ctx, dbPrefix+op, q, logger)
	var err error
	defer s.complete(q, op, span, now, logger, err)
	x := &singleBoolResult{}
	err = s.Get(ctx, x, q, tx, logger, values...)
	if err != nil {
		return false, err
	}
	return x.X, nil
}
