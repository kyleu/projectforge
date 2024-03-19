// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"net/http"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	x := util.ValueMap{"status": "OK"}
	_, _ = cutil.RespondJSON(w, "", x)
}
