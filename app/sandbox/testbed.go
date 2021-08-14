// Package sandbox $PF_IGNORE$
package sandbox

import (
	"context"

	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app"
)

var testbed = &Sandbox{Key: "testbed", Title: "Testbed", Icon: "star", Run: onTestbed}

func onTestbed(ctx context.Context, st *app.State, logger *zap.SugaredLogger) (interface{}, error) {
	ret := util.ValueMap{"status": "ok"}
	return ret, nil
}
