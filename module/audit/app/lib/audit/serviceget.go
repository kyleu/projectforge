package audit

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

func (s *Service) List(ctx context.Context, tx *sqlx.Tx, params *filter.Params, logger util.Logger) (Audits, error) {
	params = filters(params)
	q := database.SQLSelect(columnsString, tableQuoted, "", params.OrderByString(), params.Limit, params.Offset, s.db.Type)
	ret := rows{}
	err := s.db.Select(ctx, &ret, q, tx, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get audits")
	}
	return ret.ToAudits(), nil
}

func (s *Service) Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, logger util.Logger) (*Audit, error) {
	wc := defaultWC
	ret := &row{}
	q := database.SQLSelectSimple(columnsString, tableQuoted, s.db.Type, wc){{{ if .SQLServerOnly }}}
	err := s.db.Get(ctx, ret, q, tx, logger, id.String()){{{ else }}}
	err := s.db.Get(ctx, ret, q, tx, logger, id){{{ end }}}
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get audit by id [%v]", id)
	}
	return ret.ToAudit(), nil
}

const searchClause = `(
  lower(id) like {{{ .Placeholder 1 }}} or lower(app) like {{{ .Placeholder 1 }}} or lower(act) like {{{ .Placeholder 1 }}} or
  lower(client) like {{{ .Placeholder 1 }}} or lower(server) like {{{ .Placeholder 1 }}} or lower(user) like {{{ .Placeholder 1 }}} or
  lower(metadata::text) like {{{ .Placeholder 1 }}} or lower(message) like {{{ .Placeholder 1 }}}
)`

func (s *Service) Search(ctx context.Context, query string, tx *sqlx.Tx, params *filter.Params, logger util.Logger) (Audits, error) {
	params = filters(params)
	wc := searchClause
	q := database.SQLSelect(columnsString, tableQuoted, wc, params.OrderByString(), params.Limit, params.Offset, s.db.Type)
	ret := rows{}
	err := s.db.Select(ctx, &ret, q, tx, logger, "%"+strings.ToLower(query)+"%")
	if err != nil {
		return nil, err
	}
	return ret.ToAudits(), nil
}

func filters(orig *filter.Params) *filter.Params {
	return orig.Sanitize("audit", &filter.Ordering{Column: "started"})
}
