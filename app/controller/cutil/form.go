// Package cutil - Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"projectforge.dev/projectforge/app/util"
)

func ParseForm(r *http.Request, b []byte) (util.ValueMap, error) {
	ct := GetContentType(r)
	if IsContentTypeJSON(ct) {
		return parseJSONForm(b)
	}
	if IsContentTypeXML(ct) {
		return parseXMLForm(b)
	}
	if IsContentTypeYAML(ct) {
		return parseYAMLForm(b)
	}
	return parseHTTPForm(r)
}

func ParseFormAsChanges(r *http.Request, b []byte) (util.ValueMap, error) {
	ret, err := ParseForm(r, b)
	if err != nil {
		return nil, err
	}
	return ret.AsChanges()
}

func parseJSONForm(b []byte) (util.ValueMap, error) {
	ret, err := util.FromJSONAny(b)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse JSON body")
	}
	switch t := ret.(type) {
	case map[string]any:
		return t, nil
	default:
		return util.ValueMap{"resultArray": t}, nil
	}
}

func parseXMLForm(b []byte) (util.ValueMap, error) {
	ret, err := util.FromXMLMap(b)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse XML body")
	}
	return ret, nil
}

func parseYAMLForm(b []byte) (util.ValueMap, error) {
	ret := util.ValueMap{}
	err := yaml.Unmarshal(b, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse YAML body")
	}
	return ret, nil
}

func parseHTTPForm(r *http.Request) (util.ValueMap, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	ret := make(util.ValueMap, len(r.PostForm))
	for k, v := range r.PostForm {
		if len(v) == 1 {
			ret[k] = v[0]
		} else {
			ret[k] = v
		}
	}
	return ret, nil
}

func CleanID(key string, id string) string {
	if id == "" {
		return key + "-" + util.RandomString(6)
	}
	return id
}
