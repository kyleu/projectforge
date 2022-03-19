package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Insert(ctx context.Context, q string, tx *sqlx.Tx, values ...any) error {
	s.logQuery("inserting row", q, values)
	aff, err := s.execUnknown(ctx, "insert", q, tx, values...)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.Errorf("no rows affected by insert using sql [%s] and %d values", q, len(values))
	}
	return nil
}

func (s *Service) Update(ctx context.Context, q string, tx *sqlx.Tx, expected int, values ...any) (int, error) {
	return s.process(ctx, "update", "updating", "updated", q, tx, expected, values...)
}

func (s *Service) UpdateOne(ctx context.Context, q string, tx *sqlx.Tx, values ...any) error {
	_, err := s.Update(ctx, q, tx, 1, values...)
	return err
}

func (s *Service) Delete(ctx context.Context, q string, tx *sqlx.Tx, expected int, values ...any) (int, error) {
	return s.process(ctx, "delete", "deleting", "deleted", q, tx, expected, values...)
}

func (s *Service) DeleteOne(ctx context.Context, q string, tx *sqlx.Tx, values ...any) error {
	_, err := s.Delete(ctx, q, tx, 1, values...)
	if err != nil {
		return errors.Wrap(err, errMessage("delete", q, values)+"")
	}
	return err
}

func (s *Service) Exec(ctx context.Context, q string, tx *sqlx.Tx, expected int, values ...any) (int, error) {
	return s.process(ctx, "exec", "executing", "executed", q, tx, expected, values...)
}

func (s *Service) execUnknown(ctx context.Context, op string, q string, tx *sqlx.Tx, values ...any) (int, error) {
	if op == "" {
		op = "unknown"
	}
	now, ctx, span := s.newSpan(ctx, "db:"+op, q)
	var err error
	defer s.complete(q, op, span, now, err)
	var ret sql.Result
	if tx == nil {
		ret, err = s.db.ExecContext(ctx, q, values...)
	} else {
		ret, err = tx.ExecContext(ctx, q, values...)
	}
	if err != nil {
		return 0, errors.Wrap(err, errMessage(op, q, values)+"")
	}
	aff, _ := ret.RowsAffected()
	return int(aff), nil
}

func (s *Service) process(ctx context.Context, op string, key string, past string, q string, tx *sqlx.Tx, expected int, values ...any) (int, error) {
	if s.logger != nil {
		s.logQuery(fmt.Sprintf("%s [%d] rows", key, expected), q, values)
	}

	aff, err := s.execUnknown(ctx, op, q, tx, values...)
	if err != nil {
		return 0, errors.Wrap(err, errMessage(past, q, values))
	}
	if expected > -1 && aff != expected {
		const msg = "expected [%d] %s row(s), but [%d] records affected from sql [%s] with values [%s]"
		return aff, errors.Errorf(msg, expected, past, aff, q, valueStrings(values))
	}
	return aff, nil
}

func valueStrings(values []any) string {
	ret := util.StringArrayFromInterfaces(values, 256)
	return strings.Join(ret, ", ")
}
