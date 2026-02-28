package clib

import (
	"net/http"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
)

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	x := util.ValueMap{"status": util.KeyOK}
	_, _ = cutil.RespondJSON(cutil.NewWriteCounter(w), "", x)
}
