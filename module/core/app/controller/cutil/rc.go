package cutil

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

func RequestCtxToMap(r *http.Request, {{{ if .HasModule "help" }}}as{{{ else }}}_{{{ end }}} *app.State, ps *PageState) util.ValueMap {
	req := util.ValueMap{
		"url": r.URL.String(), "protocol": r.URL.Scheme,
		"host": r.URL.Host, "path": r.URL.Path,
		"queryString": r.URL.RawQuery, "headers": r.Header,
	}{{{ if .HasModule "help" }}}
	hasHelp := as.Services != nil && as.Services.Help != nil && as.Services.Help.Contains(ps.Action){{{ end }}}
	action := util.ValueMap{
		"action": ps.Action, "admin": ps.Admin, "authed": ps.Authed,
		"redirect": ps.ForceRedirect, "flashes": ps.Flashes, "breadcrumbs": ps.Breadcrumbs,
		"browser": ps.Browser, "browserVersion": ps.BrowserVersion, "os": ps.OS, "osVersion": ps.OSVersion, "platform": ps.Platform,
		"description": ps.Description, "title": ps.Title, "started": ps.Started,{{{ if .HasModule "help" }}} "help": hasHelp,{{{ end }}}
	}
	ret := util.ValueMap{"action": action, "data": ps.Data, "request": req}
	return ret
}

func PathString(r *http.Request, key string, allowEmpty bool) (string, error) {
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

func PathBool(r *http.Request, key string) (bool, error) {
	ret, err := PathString(r, key, true)
	if err != nil {
		return false, err
	}
	return ret == util.BoolTrue, nil
}

func PathInt(r *http.Request, key string) (int, error) {
	s, err := PathString(r, key, true)
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseInt(s, 10, 32)
	return int(ret), err
}

func PathUUID(r *http.Request, key string) (*uuid.UUID, error) {
	ret, err := PathString(r, key, true)
	if err != nil {
		return nil, err
	}
	return util.UUIDFromString(ret), nil
}

func PathArray(r *http.Request, key string) ([]string, error) {
	ret, err := PathString(r, key, true)
	if err != nil {
		return nil, err
	}
	return util.StringSplitAndTrim(ret, ","), nil
}

func QueryStringBool(r *http.Request, key string) bool {
	x := r.URL.Query().Get(key)
	return x == util.BoolTrue || x == "t" || x == "True" || x == "TRUE"
}

func QueryArgsMap(uri *url.URL) util.ValueMap {
	ret := make(util.ValueMap, len(uri.Query()))
	for k, v := range uri.Query() {
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
