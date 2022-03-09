// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"encoding/xml"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"

	"projectforge.dev/projectforge/app/util"
)

func ParseForm(rc *fasthttp.RequestCtx) (util.ValueMap, error) {
	ct := GetContentType(rc)
	if IsContentTypeJSON(ct) {
		return parseJSONForm(rc)
	}
	if IsContentTypeXML(ct) {
		return parseXMLForm(rc)
	}
	if IsContentTypeYAML(ct) {
		return parseYAMLForm(rc)
	}
	return parseHTTPForm(rc)
}

func ParseFormAsChanges(rc *fasthttp.RequestCtx) (util.ValueMap, error) {
	ret, err := ParseForm(rc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse form")
	}
	return ret.AsChanges()
}

func parseJSONForm(rc *fasthttp.RequestCtx) (util.ValueMap, error) {
	ret := util.ValueMap{}
	err := util.FromJSON(rc.Request.Body(), &ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse JSON body")
	}
	return ret, nil
}

func parseXMLForm(rc *fasthttp.RequestCtx) (util.ValueMap, error) {
	ret := util.ValueMap{}
	err := xml.Unmarshal(rc.Request.Body(), &ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse XML body")
	}
	return ret, nil
}

func parseYAMLForm(rc *fasthttp.RequestCtx) (util.ValueMap, error) {
	ret := util.ValueMap{}
	err := yaml.Unmarshal(rc.Request.Body(), &ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse YAML body")
	}
	return ret, nil
}

func parseHTTPForm(rc *fasthttp.RequestCtx) (util.ValueMap, error) {
	f := rc.PostArgs()
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
