package build

import (
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/util"
)

const sepKey, firstKey, thirdKey, selfKey = "sep", "1st", "3rd", "self"

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
	ss := func(expected int) {
		if state != expected {
			state = expected
			clear()
		}
	}

	for idx, imp := range imports {
		i, l := util.StringSplitLast(imp, ':', true)
		switch l {
		case sepKey:
			if state != 0 && lastLine != "" {
				state++
				clear()
			}
		case firstKey:
			if state > 1 {
				err = errors.New("1st party")
			}
			ss(1)
			first[i] = orig[idx]
			observe(i, "1st party")
		case thirdKey:
			if state > 2 {
				err = errors.New("3rd party")
			}
			ss(2)
			third[i] = orig[idx]
			observe(i, "3rd party")
		case selfKey:
			if state > 3 {
				err = errors.New("self")
			}
			ss(3)
			self[i] = orig[idx]
			observe(i, selfKey)
		default:
			return nil, errors.New("invalid type")
		}
		lastLine = l
	}
	return makeResult(first, third, self), err
}

func makeResult(first util.ValueMap, third util.ValueMap, self util.ValueMap) []string {
	ret := make([]string, 0, len(first)+len(third)+len(self))
	for _, k := range first.Keys() {
		ret = append(ret, first.GetStringOpt(k))
	}
	if len(first) > 0 && (len(third) > 0 || len(self) > 0) {
		ret = append(ret, "")
	}
	for _, k := range third.Keys() {
		ret = append(ret, third.GetStringOpt(k))
	}
	if len(third) > 0 && len(self) > 0 {
		ret = append(ret, "")
	}
	for _, k := range self.Keys() {
		ret = append(ret, self.GetStringOpt(k))
	}
	return ret
}

func importType(i string, self string) string {
	if i == "" {
		return sepKey
	}
	if strings.HasPrefix(i, self) || strings.HasPrefix(i, "{{{ .Package }}}") {
		return selfKey
	}
	if strings.Contains(i, ".") {
		return thirdKey
	}
	return firstKey
}
