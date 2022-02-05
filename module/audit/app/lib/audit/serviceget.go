package audit

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
)

func (s *Service) List(ctx context.Context, tx *sqlx.Tx, params *filter.Params) (Audits, error) {
	params = filters(params)
	q := database.SQLSelect(columnsString, tableQuoted, "", params.OrderByString(), params.Limit, params.Offset)
	ret := dtos{}
	err := s.db.Select(ctx, &ret, q, tx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get audits")
	}
	return ret.ToAudits(), nil
}

func (s *Service) Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*Audit, error) {
	wc := defaultWC
	ret := &dto{}
	q := database.SQLSelectSimple(columnsString, tableQuoted, wc)
	err := s.db.Get(ctx, ret, q, tx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get audit by id [%v]", id)
	}
	return ret.ToAudit(), nil
}

const searchClause = `(
  lower(id) like $1 or lower(app) like $1 or lower(act) like $1 or
  lower(client) like $1 or lower(server) like $1 or lower(user) like $1 or
  lower(metadata::text) like $1 or lower(message) like $1
)`

func (s *Service) Search(ctx context.Context, query string, tx *sqlx.Tx, params *filter.Params) (Audits, error) {
	params = filters(params)
	wc := searchClause
	q := database.SQLSelect(columnsString, tableQuoted, wc, params.OrderByString(), params.Limit, params.Offset)
	ret := dtos{}
	err := s.db.Select(ctx, &ret, q, tx, "%"+strings.ToLower(query)+"%")
	if err != nil {
		return nil, err
	}
	return ret.ToAudits(), nil
}

func filters(orig *filter.Params) *filter.Params {
	return orig.Sanitize("audit", &filter.Ordering{Column: "started"})
}
