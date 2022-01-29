package audit

import (
	"go.uber.org/zap"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/filter"
)

type Service struct {
	db     *database.Service
	logger *zap.SugaredLogger
}

func NewService(db *database.Service, logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "audit")
	filter.AllowedColumns["audit"] = columns
	return &Service{db: db, logger: logger}
}

func filters(orig *filter.Params) *filter.Params {
	return orig.Sanitize("audit", &filter.Ordering{Column: "started"})
}
