package audit

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

func (s *Service) RecordsForAudit(ctx context.Context, tx *sqlx.Tx, auditID uuid.UUID, params *filter.Params, logger util.Logger) (Records, error) {
	params = params.Sanitize("audit_record", &filter.Ordering{Column: "occurred"})
	wc := `"audit_id" = $1`
	q := database.SQLSelect(recordColumnsString, recordTableQuoted, wc, params.OrderByString(), params.Limit, params.Offset)
	ret := recordDTOs{}
	err := s.db.Select(ctx, &ret, q, tx, logger, auditID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get audit records by audit [%s]", auditID.String())
	}
	return ret.ToRecords(), nil
}

func (s *Service) GetRecord(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, logger util.Logger) (*Record, error) {
	q := database.SQLSelectSimple(recordColumnsString, recordTableQuoted, "id = $1")
	ret := &recordDTO{}
	err := s.db.Get(ctx, ret, q, tx, logger, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get audit record by id [%s]", id.String())
	}
	return ret.ToRecord(), nil
}

func (s *Service) CreateRecords(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*Record) error {
	if len(models) == 0 {
		return nil
	}
	q := database.SQLInsert(recordTableQuoted, recordColumnsQuoted, len(models), "")
	vals := make([]any, 0, len(models)*len(recordColumnsQuoted))
	for _, arg := range models {
		vals = append(vals, arg.ToData()...)
	}
	return s.db.Insert(ctx, q, tx, logger, vals...)
}
