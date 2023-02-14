// Package sandbox - $PF_GENERATE_ONCE$
package sandbox

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

var testbed = &Sandbox{Key: "testbed", Title: "Testbed", Icon: "star", Run: onTestbed}

func onTestbed(_ context.Context, _ *app.State, _ util.Logger) (any, error) {
	ret := util.ValueMap{"status": "ok"}
	return ret, nil
}
