package database

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (s *Service) Query(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) (*sqlx.Rows, error) {
	ctx, span := s.newSpan(ctx, "query", q)
	defer span.End()
	s.logQuery("running raw query", q, values)
	if tx == nil {
		return s.db.QueryxContext(ctx, q, values...)
	}
	return tx.QueryxContext(ctx, q, values...)
}

func (s *Service) QueryRows(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(ctx, q, tx, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	ret := []map[string]interface{}{}
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

func (s *Service) QuerySingleRow(ctx context.Context, q string, tx *sqlx.Tx, values ...interface{}) (map[string]interface{}, error) {
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
	ctx, span := s.newSpan(ctx, "select", q)
	defer span.End()
	s.logQuery(fmt.Sprintf("selecting rows of type [%T]", dest), q, values)
	if tx == nil {
		return s.db.SelectContext(ctx, dest, q, values...)
	}
	return tx.SelectContext(ctx, dest, q, values...)
}

func (s *Service) Get(ctx context.Context, dto interface{}, q string, tx *sqlx.Tx, values ...interface{}) error {
	ctx, span := s.newSpan(ctx, "get", q)
	defer span.End()
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
	ctx, span := s.newSpan(ctx, "single-int", q)
	defer span.End()
	x := &singleIntResult{}
	err := s.Get(ctx, x, q, tx, values...)
	if err != nil {
		return -1, errors.Wrap(err, "returned value is not an integer")
	}
	if x.X == nil {
		return 0, nil
	}
	return *x.X, nil
}
