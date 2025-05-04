// $PF_GENERATE_ONCE$
package user

import (
	"github.com/google/uuid"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

type Service struct {
	files  filesystem.FileLoader
	logger util.Logger
}

func NewService(f filesystem.FileLoader, logger util.Logger) *Service {
	return &Service{files: f, logger: logger}
}

func filters(orig *filter.Params) *filter.Params {
	return orig.Sanitize("user", &filter.Ordering{Column: "created"})
}

func dirFor(userID uuid.UUID) string {
	return util.StringPath("users", userID.String())
}
