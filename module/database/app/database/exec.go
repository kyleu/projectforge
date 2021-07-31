package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Insert(q string, tx *sqlx.Tx, values ...interface{}) error {
	s.logQuery("inserting row", q, values)
	aff, err := s.execUnknown(q, tx, values...)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.Errorf("no rows affected by insert using sql [%s] and %d values", q, len(values))
	}
	return nil
}

func (s *Service) Update(q string, tx *sqlx.Tx, expected int, values ...interface{}) (int, error) {
	return s.process("updating", "updated", q, tx, expected, values...)
}

func (s *Service) UpdateOne(q string, tx *sqlx.Tx, values ...interface{}) error {
	_, err := s.Update(q, tx, 1, values...)
	return err
}

func (s *Service) Delete(q string, tx *sqlx.Tx, expected int, values ...interface{}) (int, error) {
	return s.process("deleting", "deleted", q, tx, expected, values...)
}

func (s *Service) DeleteOne(q string, tx *sqlx.Tx, values ...interface{}) error {
	_, err := s.Delete(q, tx, 1, values...)
	if err != nil {
		return errors.Wrap(err, errMessage("delete", q, values)+"")
	}
	return err
}

func (s *Service) Exec(q string, tx *sqlx.Tx, expected int, values ...interface{}) (int, error) {
	return s.process("executing", "executed", q, tx, expected, values...)
}

func (s *Service) execUnknown(q string, tx *sqlx.Tx, values ...interface{}) (int, error) {
	var err error
	var ret sql.Result
	if tx == nil {
		r, e := s.db.Exec(q, values...)
		ret = r
		err = e
	} else {
		r, e := tx.Exec(q, values...)
		ret = r
		err = e
	}
	if err != nil {
		return 0, errors.Wrap(err, errMessage("exec", q, values)+"")
	}
	aff, _ := ret.RowsAffected()
	// if err != nil {
	// 	return 0, err
	// }
	return int(aff), nil
}

func (s *Service) process(key string, past string, q string, tx *sqlx.Tx, expected int, values ...interface{}) (int, error) {
	if s.logger != nil {
		s.logQuery(fmt.Sprintf("%s [%d] rows", key, expected), q, values)
	}

	aff, err := s.execUnknown(q, tx, values...)
	if err != nil {
		return 0, errors.Wrap(err, errMessage(past, q, values))
	}
	if expected > -1 && aff != expected {
		const msg = "expected [%d] %s row(s), but [%d] records affected from sql [%s] with values [%s]"
		return aff, errors.Errorf(msg, expected, past, aff, q, valueStrings(values))
	}
	return aff, nil
}

func valueStrings(values []interface{}) string {
	ret := util.StringArrayFromInterfaces(values)
	return strings.Join(ret, ", ")
}
