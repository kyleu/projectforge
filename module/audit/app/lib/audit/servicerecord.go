package audit

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
)

func (s *Service) RecordsForAudit(ctx context.Context, tx *sqlx.Tx, auditID uuid.UUID, params *filter.Params) (Records, error) {
	params = recordFilters(params)
	wc := `"audit_id" = $1`
	q := database.SQLSelect(recordColumnsString, recordTableQuoted, wc, params.OrderByString(), params.Limit, params.Offset)
	ret := recordDTOs{}
	err := s.db.Select(ctx, &ret, q, tx, auditID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get audit records by audit [%s]", auditID.String())
	}
	return ret.ToRecords(), nil
}

func (s *Service) GetRecord(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*Record, error) {
	q := database.SQLSelectSimple(recordColumnsString, recordTableQuoted, "id = $1")
	ret := &recordDTO{}
	err := s.db.Get(ctx, ret, q, tx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get audit record by id [%s]", id.String())
	}
	return ret.ToRecord(), nil
}
