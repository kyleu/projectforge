package audit

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

type Service struct {
	db     *database.Service
	logger *zap.SugaredLogger
}

func NewService(db *database.Service, logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "audit")
	filter.AllowedColumns["audit"] = columns
	filter.AllowedColumns["audit_record"] = recordColumns
	return &Service{db: db, logger: logger}
}

func (s *Service) ApplyObj(ctx context.Context, a *Audit, l any, r any, md util.ValueMap) (*Audit, Records, error) {
	o := r
	if o == nil {
		o = l
	}
	d := util.DiffObjects(l, r, "")
	rec := NewRecord(a.ID, fmt.Sprintf("%T", o), fmt.Sprint(o), d, md)
	return s.Apply(ctx, a, rec)
}

func (s *Service) Apply(ctx context.Context, a *Audit, r ...*Record) (*Audit, Records, error) {
	tx, err := s.db.StartTransaction(s.logger)
	defer func() { _ = tx.Rollback() }()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to start transaction")
	}
	err = s.Create(ctx, tx, a)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to insert audit")
	}
	err = s.CreateRecords(ctx, tx, r...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to insert audit records")
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to commit transaction")
	}
	return a, r, nil
}
