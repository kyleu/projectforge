package clib

import (
	"net/http"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	x := util.ValueMap{"status": "OK"}
	_, _ = cutil.RespondJSON(w, "", x)
}
