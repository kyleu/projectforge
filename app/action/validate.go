package action

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
)

func onValidate(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeValidate, pm.Prj, pm.Cfg, pm.Logger)
	errs := project.Validate(pm.Prj, pm.MSvc.Deps())
	for _, err := range errs {
		ret = ret.WithError(errors.Errorf("%s: %s", err.Code, err.Message))
	}
	return ret
}
