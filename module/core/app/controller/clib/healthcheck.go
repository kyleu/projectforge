package clib

import (
	"net/http"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	x := util.ValueMap{"status": util.KeyOK}
	_, _ = cutil.RespondJSON(cutil.NewWriteCounter(w), "", x)
}
