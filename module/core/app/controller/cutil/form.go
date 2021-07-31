package cutil

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/util"
)

func ParseForm(ctx *fasthttp.RequestCtx) (util.ValueMap, error) {
	if ct := GetContentType(ctx); IsContentTypeJSON(ct) {
		return parseJSONForm(ctx)
	}
	return parseHTTPForm(ctx)
}

func ParseFormAsChanges(ctx *fasthttp.RequestCtx) (util.ValueMap, error) {
	ret, err := ParseForm(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse form")
	}
	return ret.AsChanges()
}

func parseJSONForm(ctx *fasthttp.RequestCtx) (util.ValueMap, error) {
	ret := util.ValueMap{}
	err := util.FromJSON(ctx.Request.Body(), &ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse JSON body")
	}
	return ret, nil
}

func parseHTTPForm(ctx *fasthttp.RequestCtx) (util.ValueMap, error) {
	f := ctx.PostArgs()
	ret := make(util.ValueMap, f.Len())
	f.VisitAll(func(key []byte, value []byte) {
		k := string(key)
		xs := f.PeekMulti(k)
		v := make([]string, 0, len(xs))
		for _, x := range xs {
			v = append(v, string(x))
		}
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			ret[k] = v
		}
	})
	return ret, nil
}
