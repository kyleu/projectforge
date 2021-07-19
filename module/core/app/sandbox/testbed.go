// Package sandbox $PF_IGNORE$
package sandbox

import (
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/util"

	"go.uber.org/zap"
)

var testbed = &Sandbox{Key: "testbed", Title: "Testbed", Icon: "star", Run: onTestbed}

func onTestbed(st *app.State, logger *zap.SugaredLogger) (interface{}, error) {
	ret := util.ValueMap{"status": "ok"}
	return ret, nil
}
