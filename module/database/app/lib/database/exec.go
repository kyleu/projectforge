package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"

	"{{{ .Package }}}/app/util"
)

const dbPrefix = "db:"

func (s *Service) Insert(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) error {
	if s.ReadOnly {
		return errors.Errorf("cannot run [insert] statements in read-only database [%s]", q)
	}
	f := s.logQuery(ctx, "inserting row", q, logger, values)
	aff, err := s.execUnknown(ctx, "insert", q, tx, logger, values...)
	defer f(aff, fmt.Sprintf("inserted [%d] rows", aff), err)
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.Errorf("no rows affected by insert using sql [%s] and %d values", q, len(values))
	}
	return nil
}

func (s *Service) InsertMap(ctx context.Context, table string, m util.ValueMap, tx *sqlx.Tx, logger util.Logger, values ...any) error {
	columns := util.ArraySorted(lo.Filter(maps.Keys(m), func(x string, _ int) bool {
		return !strings.HasPrefix(x, "~")
	}))
	q := SQLInsert("\"order\"", columns, 1, s.Type)
	vals := lo.Map(columns, func(k string, _ int) any {
		return m[k]
	})
	return s.Insert(ctx, q, tx, logger, vals...)
}

func (s *Service) Update(ctx context.Context, q string, tx *sqlx.Tx, expected int, logger util.Logger, values ...any) (int, error) {
	return s.process(ctx, "update", "updating", "updated", q, tx, expected, logger, values...)
}

func (s *Service) UpdateOne(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) error {
	_, err := s.Update(ctx, q, tx, 1, logger, values...)
	return err
}

func (s *Service) Delete(ctx context.Context, q string, tx *sqlx.Tx, expected int, logger util.Logger, values ...any) (int, error) {
	return s.process(ctx, "delete", "deleting", "deleted", q, tx, expected, logger, values...)
}

func (s *Service) DeleteOne(ctx context.Context, q string, tx *sqlx.Tx, logger util.Logger, values ...any) error {
	_, err := s.Delete(ctx, q, tx, 1, logger, values...)
	if err != nil {
		return sqlErrMessage(err, "delete", q, values)
	}
	return err
}

func (s *Service) Exec(ctx context.Context, q string, tx *sqlx.Tx, expected int, logger util.Logger, values ...any) (int, error) {
	return s.process(ctx, "exec", "executing", "executed", q, tx, expected, logger, values...)
}

func (s *Service) execUnknown(ctx context.Context, op string, q string, tx *sqlx.Tx, logger util.Logger, values ...any) (int, error) {
	if op == "" {
		op = "unknown"
	}
	now, ctx, span, logger := s.newSpan(ctx, dbPrefix+op, q, logger)
	var err error
	defer s.complete(q, op, span, now, logger, err)
	var ret sql.Result
	if tx == nil {
		ret, err = s.db.ExecContext(ctx, q, values...)
	} else {
		ret, err = tx.ExecContext(ctx, q, values...)
	}
	if err != nil {
		return 0, sqlErrMessage(err, op, q, values)
	}
	aff, _ := ret.RowsAffected()
	return int(aff), nil
}

func (s *Service) process(
	ctx context.Context, op string, key string, past string, q string, tx *sqlx.Tx, expected int, logger util.Logger, values ...any,
) (int, error) {
	if s.ReadOnly {
		return 0, errors.Errorf("cannot run [%q] statements in read-only database [%s]", op, q)
	}
	msg := fmt.Sprintf("%s [%d] rows", key, expected)
	f := s.logQuery(ctx, msg, q, logger, values)
	aff, err := s.execUnknown(ctx, op, q, tx, logger, values...)
	defer f(aff, msg, err)
	if err != nil {
		return 0, sqlErrMessage(err, past, q, values)
	}
	if expected > -1 && aff != expected {
		const msg = "expected [%d] %s row(s), but [%d] records affected from sql [%s] with values [%s]"
		return aff, errors.Errorf(msg, expected, past, aff, q, valueStrings(values))
	}
	return aff, nil
}

func valueStrings(values []any) string {
	ret := util.StringArrayFromAny(values, 256)
	return strings.Join(ret, ", ")
}
