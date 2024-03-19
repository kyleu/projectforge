// Package cutil - Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

func RequestCtxToMap(w http.ResponseWriter, r *http.Request, as *app.State, ps *PageState) util.ValueMap {
	req := util.ValueMap{
		"url": r.URL.String(), "protocol": r.URL.Scheme,
		"host": r.URL.Host, "path": r.URL.Path,
		"queryString": r.URL.RawQuery, "headers": r.Header,
	}
	hasHelp := as.Services != nil && as.Services.Help != nil && as.Services.Help.Contains(ps.Action)
	action := util.ValueMap{
		"action": ps.Action, "admin": ps.Admin, "authed": ps.Authed,
		"redirect": ps.ForceRedirect, "flashes": ps.Flashes, "breadcrumbs": ps.Breadcrumbs,
		"browser": ps.Browser, "browserVersion": ps.BrowserVersion, "os": ps.OS, "osVersion": ps.OSVersion, "platform": ps.Platform,
		"description": ps.Description, "title": ps.Title, "started": ps.Started, "help": hasHelp,
	}
	ret := util.ValueMap{"action": action, "data": ps.Data, "request": req}
	// $PF_SECTION_START(debugstuff)$
	// $PF_SECTION_END(debugstuff)$
	return ret
}

func RCRequiredString(r *http.Request, key string, allowEmpty bool) (string, error) {
	v := mux.Vars(r)[key]
	if (!allowEmpty) && v == "" {
		return v, errors.Errorf("must provide [%s] in path", key)
	}
	v, err := url.QueryUnescape(v)
	if err != nil {
		return "", err
	}
	return v, nil
}

func RCRequiredBool(r *http.Request, key string) (bool, error) {
	ret, err := RCRequiredString(r, key, true)
	if err != nil {
		return false, err
	}
	return ret == util.BoolTrue, nil
}

func RCRequiredInt(r *http.Request, key string) (int, error) {
	s, err := RCRequiredString(r, key, true)
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseInt(s, 10, 32)
	return int(ret), err
}

func RCRequiredUUID(r *http.Request, key string) (*uuid.UUID, error) {
	ret, err := RCRequiredString(r, key, true)
	if err != nil {
		return nil, err
	}
	return util.UUIDFromString(ret), nil
}

func RCRequiredArray(r *http.Request, key string) ([]string, error) {
	ret, err := RCRequiredString(r, key, true)
	if err != nil {
		return nil, err
	}
	return util.StringSplitAndTrim(ret, ","), nil
}

func QueryStringBool(r *http.Request, key string) bool {
	x := r.URL.Query().Get(key)
	return x == util.BoolTrue || x == "t" || x == "True" || x == "TRUE"
}

func QueryArgsMap(r *http.Request) util.ValueMap {
	ret := make(util.ValueMap, len(r.Header))
	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			ret[k] = v
		}
	}
	return ret
}

func RequestHeadersMap(r *http.Request) util.ValueMap {
	ret := make(util.ValueMap, len(r.Header))
	for k, v := range r.Header {
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			ret[k] = v
		}
	}
	return ret
}

func ResponseHeadersMap(w http.ResponseWriter) util.ValueMap {
	ret := make(util.ValueMap, len(w.Header()))
	for k, v := range w.Header() {
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			ret[k] = v
		}
	}
	return ret
}
