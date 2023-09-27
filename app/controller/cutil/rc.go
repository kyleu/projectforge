// Package cutil - Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/util"
)

func RequestCtxToMap(rc *fasthttp.RequestCtx, data any) util.ValueMap {
	req := util.ValueMap{
		"id":          rc.ID(),
		"url":         rc.URI().String(),
		"protocol":    string(rc.Request.URI().Scheme()),
		"host":        string(rc.Request.URI().Host()),
		"path":        string(rc.Request.URI().Path()),
		"queryString": string(rc.Request.URI().QueryString()),
		"headers":     RequestHeadersMap(rc),
		"bodySize":    len(rc.Request.Body()),
		"string":      rc.Request.String(),
	}
	rsp := util.ValueMap{
		"code":     rc.Response.StatusCode(),
		"bodySize": len(rc.Response.Body()),
		"headers":  ResponseHeadersMap(rc),
		"string":   rc.Response.String(),
	}
	return util.ValueMap{"data": data, "request": req, "response": rsp}
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
