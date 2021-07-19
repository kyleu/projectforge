package database

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (s *Service) Query(q string, tx *sqlx.Tx, values ...interface{}) (*sqlx.Rows, error) {
	s.logQuery("running raw query", q, values)
	if tx == nil {
		return s.db.Queryx(q, values...)
	}
	return tx.Queryx(q, values...)
}

func (s *Service) QueryRows(q string, tx *sqlx.Tx, values ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(q, tx, values...)
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

func (s *Service) QuerySingleRow(q string, tx *sqlx.Tx, values ...interface{}) (map[string]interface{}, error) {
	rows, err := s.QueryRows(q, tx, values...)
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

func (s *Service) Select(dest interface{}, q string, tx *sqlx.Tx, values ...interface{}) error {
	s.logQuery(fmt.Sprintf("selecting rows of type [%T]", dest), q, values)
	if tx == nil {
		return s.db.Select(dest, q, values...)
	}
	return tx.Select(dest, q, values...)
}

func (s *Service) Get(dto interface{}, q string, tx *sqlx.Tx, values ...interface{}) error {
	s.logQuery(fmt.Sprintf("getting single row of type [%T]", dto), q, values)
	if tx == nil {
		return s.db.Get(dto, q, values...)
	}
	return tx.Get(dto, q, values...)
}

type singleIntResult struct {
	X *int64 `db:"x"`
}

func (s *Service) SingleInt(q string, tx *sqlx.Tx, values ...interface{}) (int64, error) {
	x := &singleIntResult{}
	var err error
	if tx == nil {
		err = s.db.Get(x, q, values...)
	} else {
		err = tx.Get(x, q, values...)
	}
	if err != nil {
		return -1, errors.Wrap(err, "returned value is not an integer")
	}
	if x.X == nil {
		return 0, nil
	}
	return *x.X, nil
}
