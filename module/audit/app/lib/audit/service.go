package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

type Service struct {
	db     *database.Service
	logger util.Logger
}

func NewService(db *database.Service, logger util.Logger) *Service {
	logger = logger.With("svc", "audit")
	filter.AllowedColumns["audit"] = columns
	filter.AllowedColumns["audit_record"] = recordColumns
	return &Service{db: db, logger: logger}
}

func (s *Service) Apply(ctx context.Context, a *Audit, logger util.Logger, r ...*Record) (*Audit, Records, error) {
	tx, err := s.db.StartTransaction(logger)
	defer func() { _ = tx.Rollback() }()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to start transaction")
	}
	err = s.Create(ctx, tx, logger, a)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to insert audit")
	}
	err = s.CreateRecords(ctx, tx, logger, r...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to insert audit records")
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to commit transaction")
	}
	return a, r, nil
}

func (s *Service) ApplyObj(ctx context.Context, a *Audit, l any, r any, md util.ValueMap, logger util.Logger) (*Audit, Records, error) {
	o := r
	if o == nil {
		o = l
	}
	d := util.DiffObjects(l, r, "")
	rec := NewRecord(a.ID, fmt.Sprintf("%T", o), fmt.Sprint(o), d, md)
	return s.Apply(ctx, a, logger, rec)
}

func (s *Service) ApplyObjSimple(ctx context.Context, act string, msg string, l any, r any, md util.ValueMap, logger util.Logger) (*Audit, Records, error) {
	a := New(act, "", "", "", md, msg)
	a.Completed = time.Now()
	return s.ApplyObj(ctx, a, l, r, md, logger)
}
