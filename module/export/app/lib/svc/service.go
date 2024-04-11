package svc

import (
	"context"

	"github.com/jmoiron/sqlx"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

type ServiceBase[Mdl Model, Seq any] interface {
	ListSQL(ctx context.Context, tx *sqlx.Tx, sql string, logger util.Logger, values ...any) (Seq, error)
	ListWhere(ctx context.Context, tx *sqlx.Tx, where string, params *filter.Params, logger util.Logger, values ...any) (Seq, error)
	Random(ctx context.Context, tx *sqlx.Tx, logger util.Logger) (Mdl, error)
	Create(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...Mdl) error
	CreateChunked(ctx context.Context, tx *sqlx.Tx, chunkSize int, logger util.Logger, models ...Mdl) error
	Update(ctx context.Context, tx *sqlx.Tx, model Mdl, logger util.Logger) error
	Save(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...Mdl) error
	SaveChunked(ctx context.Context, tx *sqlx.Tx, chunkSize int, logger util.Logger, models ...Mdl) error
	DeleteWhere(ctx context.Context, tx *sqlx.Tx, wc string, expected int, logger util.Logger, values ...any) error
}

type Service[Mdl Model, Seq any] interface {
	ServiceBase[Mdl, Seq]
	Count(ctx context.Context, tx *sqlx.Tx, whereClause string, logger util.Logger, args ...any) (int, error)
	List(ctx context.Context, tx *sqlx.Tx, params *filter.Params, logger util.Logger) (Seq, error)
}

type ServiceSearch[Seq any] interface {
	Search(ctx context.Context, query string, tx *sqlx.Tx, params *filter.Params, logger util.Logger) (Seq, error)
}

type ServiceID[Mdl Model, Seq any, ID any] interface {
	Service[Mdl, Seq]
	Get(ctx context.Context, tx *sqlx.Tx, id ID, logger util.Logger) (Mdl, error)
	GetMultiple(ctx context.Context, tx *sqlx.Tx, params *filter.Params, logger util.Logger, ids ...ID) (Seq, error)
}

type ServiceSoftDelete[Mdl Model, Seq any] interface {
	ServiceBase[Mdl, Seq]
	Count(ctx context.Context, tx *sqlx.Tx, whereClause string, includeDeleted bool, logger util.Logger, args ...any) (int, error)
	List(ctx context.Context, tx *sqlx.Tx, params *filter.Params, includeDeleted bool, logger util.Logger) (Seq, error)
}

type ServiceSoftDeleteID[Mdl Model, Seq any, ID any] interface {
	ServiceSoftDelete[Mdl, Seq]
	Get(ctx context.Context, tx *sqlx.Tx, id ID, includeDeleted bool, logger util.Logger) (Mdl, error)
	GetMultiple(ctx context.Context, tx *sqlx.Tx, params *filter.Params, includeDeleted bool, logger util.Logger, ids ...ID) (Seq, error)
}
