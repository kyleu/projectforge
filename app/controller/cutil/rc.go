// Package cutil - Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

func RequestCtxToMap(rc *fasthttp.RequestCtx, as *app.State, ps *PageState) util.ValueMap {
	req := util.ValueMap{
		"url": rc.URI().String(), "protocol": string(rc.Request.URI().Scheme()),
		"host": string(rc.Request.URI().Host()), "path": string(rc.Request.URI().Path()),
		"queryString": string(rc.Request.URI().QueryString()), "headers": RequestHeadersMap(rc),
		"bodySize": len(rc.Request.Body()),
	}
	rsp := util.ValueMap{"code": rc.Response.StatusCode(), "bodySize": len(rc.Response.Body()), "headers": ResponseHeadersMap(rc)}
	hasHelp := as.Services != nil && as.Services.Help != nil && as.Services.Help.Contains(ps.Action)
	action := util.ValueMap{
		"action": ps.Action, "admin": ps.Admin, "authed": ps.Authed,
		"redirect": ps.ForceRedirect, "flashes": ps.Flashes, "breadcrumbs": ps.Breadcrumbs,
		"browser": ps.Browser, "browserVersion": ps.BrowserVersion, "os": ps.OS, "osVersion": ps.OSVersion, "platform": ps.Platform,
		"description": ps.Description, "title": ps.Title, "started": ps.Started, "help": hasHelp,
	}
	ret := util.ValueMap{"action": action, "data": ps.Data, "request": req, "response": rsp}
	// $PF_SECTION_START(debugstuff)$
	// $PF_SECTION_END(debugstuff)$
	return ret
}

func RCRequiredString(rc *fasthttp.RequestCtx, key string, allowEmpty bool) (string, error) {
	v, ok := rc.UserValue(key).(string)
	if !ok || ((!allowEmpty) && v == "") {
		return v, errors.Errorf("must provide [%s] in path", key)
	}
	v, err := url.QueryUnescape(v)
	if err != nil {
		return "", err
	}
	return v, nil
}

func RCRequiredBool(rc *fasthttp.RequestCtx, key string) (bool, error) {
	ret, err := RCRequiredString(rc, key, true)
	if err != nil {
		return false, err
	}
	return ret == util.BoolTrue, nil
}

func RCRequiredInt(rc *fasthttp.RequestCtx, key string) (int, error) {
	s, err := RCRequiredString(rc, key, true)
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseInt(s, 10, 32)
	return int(ret), err
}

func RCRequiredUUID(rc *fasthttp.RequestCtx, key string) (*uuid.UUID, error) {
	ret, err := RCRequiredString(rc, key, true)
	if err != nil {
		return nil, err
	}
	return util.UUIDFromString(ret), nil
}

func RCRequiredArray(rc *fasthttp.RequestCtx, key string) ([]string, error) {
	ret, err := RCRequiredString(rc, key, true)
	if err != nil {
		return nil, err
	}
	return util.StringSplitAndTrim(ret, ","), nil
}

func QueryStringBool(rc *fasthttp.RequestCtx, key string) bool {
	x := string(rc.URI().QueryArgs().Peek(key))
	return x == util.BoolTrue || x == "t" || x == "True" || x == "TRUE"
}

func QueryArgsMap(rc *fasthttp.RequestCtx) util.ValueMap {
	ret := make(util.ValueMap, rc.QueryArgs().Len())
	rc.QueryArgs().VisitAll(func(k []byte, v []byte) {
		curr, _ := ret.GetStringArray(string(k), true)
		ret[string(k)] = append(curr, string(v))
	})
	return ret
}

func RequestHeadersMap(rc *fasthttp.RequestCtx) map[string]string {
	ret := make(map[string]string, rc.Request.Header.Len())
	rc.Request.Header.VisitAll(func(k, v []byte) {
		ret[string(k)] = string(v)
	})
	return ret
}

func ResponseHeadersMap(rc *fasthttp.RequestCtx) map[string]string {
	ret := make(map[string]string, rc.Request.Header.Len())
	rc.Response.Header.VisitAll(func(k, v []byte) {
		ret[string(k)] = string(v)
	})
	return ret
}
