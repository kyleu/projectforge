// Package sandbox - $PF_IGNORE$
package sandbox

import (
	"context"

	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

var testbed = &Sandbox{Key: "testbed", Title: "Testbed", Icon: "star", Run: onTestbed}

func onTestbed(_ context.Context, _ *app.State, _ *zap.SugaredLogger) (any, error) {
	ret := util.ValueMap{"status": "ok"}
	return ret, nil
}
