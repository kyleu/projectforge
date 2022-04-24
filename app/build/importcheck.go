package build

import (
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/util"
)

func check(imports []string, orig []string) ([]string, error) {
	var err error
	var state int
	var lastLine string
	var observed []string
	first, third, self := util.ValueMap{}, util.ValueMap{}, util.ValueMap{}
	observe := func(key string, i string) {
		for _, ob := range observed {
			if ob > i {
				err = errors.Errorf("%s sorting", key)
			}
		}
		observed = append(observed, i)
	}
	clear := func() {
		observed = []string{}
	}

	for idx, imp := range imports {
		i, l := util.StringSplitLast(imp, ':', true)
		switch l {
		case "sep":
			if state != 0 && lastLine != "" {
				state++
				clear()
			}
		case "1st":
			if state > 1 {
				err = errors.New("1st party")
			}
			if state != 1 {
				state = 1
				clear()
			}
			first[i] = orig[idx]
			observe(i, "1st party")
		case "3rd":
			if state > 2 {
				err = errors.New("3rd party")
			}
			if state != 2 {
				state = 2
				clear()
			}
			third[i] = orig[idx]
			observe(i, "3rd party")
		case "self":
			if state > 3 {
				err = errors.New("self")
			}
			if state != 3 {
				state = 3
				clear()
			}
			self[i] = orig[idx]
			observe(i, "self")
		default:
			return nil, errors.New("invalid type")
		}
		lastLine = l
	}
	var ret []string
	for _, k := range first.Keys() {
		ret = append(ret, first[k].(string))
	}
	if len(first) > 0 && (len(third) > 0 || len(self) > 0) {
		ret = append(ret, "")
	}
	for _, k := range third.Keys() {
		ret = append(ret, third[k].(string))
	}
	if len(third) > 0 && len(self) > 0 {
		ret = append(ret, "")
	}
	for _, k := range self.Keys() {
		ret = append(ret, self[k].(string))
	}
	return ret, err
}

func importType(i string, self string) string {
	if i == "" {
		return "sep"
	}
	if strings.HasPrefix(i, self) || strings.HasPrefix(i, "{{{ .Package }}}") {
		return "self"
	}
	if strings.Contains(i, ".") {
		return "3rd"
	}
	return "1st"
}
