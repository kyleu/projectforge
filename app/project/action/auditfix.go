package action

import (
	"context"
	"github.com/pkg/errors"
)

func auditRemove(ctx context.Context, fn string, pm *PrjAndMods, ret *Result) error {
	if fn == "" {
		paths := ret.Modules.Paths(false)
		if len(paths) == 0 {
			x := &Result{}
			err := auditRun(pm, x)
			if err != nil {
				return err
			}
			paths = x.Modules.Paths(false)
		}
		for _, path := range paths {
			err := auditRemove(ctx, path, pm, ret)
			if err != nil {
				return errors.Wrapf(err, "can't fix audit of [%s]", path)
			}
		}
		return nil
	}
	ret.AddLog("removed [%s]", fn)
	return pm.FS.Remove(fn, pm.Logger)
}
