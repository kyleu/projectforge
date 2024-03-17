package build

import (
	"slices"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const sepKey, firstKey, thirdKey, selfKey = "sep", "1st", "3rd", "self"

func check(imports []string, orig []string) ([]string, []string, error) {
	var err error
	var state int
	var lastLine string
	var observed []string
	var imps []string
	var lastSep bool
	first, third, self := util.ValueMap{}, util.ValueMap{}, util.ValueMap{}
	var firsts, thirds, selfs []string

	observe := func(key string, i string) {
		lo.ForEach(observed, func(ob string, _ int) {
			if ob > i {
				err = errors.Errorf("%s sorting", key)
			}
		})
		observed = append(observed, i)
	}
	clr := func() {
		observed = []string{}
	}
	chk := func(tgt int, msg string) {
		if state != tgt {
			if !lastSep && len(first) > 0 {
				err = errors.New("missing separator")
			}
			state = tgt
			clr()
		}
		lastSep = false
		if state > tgt {
			err = errors.New(msg)
		}
	}

	for idx, imp := range imports {
		i, l := util.StringSplitLast(imp, ':', true)
		switch l {
		case sepKey:
			if state != 0 && lastLine != "" {
				state++
				clr()
			}
			lastSep = true
		case firstKey:
			if state > 1 {
				err = errors.New("invalid ordering")
			}
			if state != 1 {
				state = 1
				clr()
			}
			lastSep = false
			if state > 1 {
				err = errors.New("1st party")
			}
			chk(1, "1st party")
			first[i] = orig[idx]
			firsts = append(firsts, i)
			imps = append(imps, imp)
			observe(i, "1st party")
		case thirdKey:
			chk(2, "3rd party")
			third[i] = orig[idx]
			thirds = append(thirds, i)
			imps = append(imps, imp)
			observe(i, "3rd party")
		case selfKey:
			chk(3, "self")
			self[i] = orig[idx]
			selfs = append(selfs, i)
			imps = append(imps, imp)
			observe(i, selfKey)
		default:
			return nil, nil, errors.New("invalid type")
		}
		lastLine = l
	}
	if !slices.IsSorted(firsts) {
		err = errors.New("first-party imports are not sorted")
	}
	if !slices.IsSorted(thirds) {
		err = errors.New("third-party imports are not sorted")
	}
	if !slices.IsSorted(selfs) {
		err = errors.New("this app's imports are not sorted")
	}
	return imps, makeResult(first, third, self), err
}

func makeResult(first util.ValueMap, third util.ValueMap, self util.ValueMap) []string {
	ret := util.NewStringSlice(make([]string, 0, len(first)+len(third)+len(self)))
	lo.ForEach(first.Keys(), func(k string, _ int) {
		s, _ := first[k].(string)
		ret.Push(s)
	})
	if len(first) > 0 && (len(third) > 0 || len(self) > 0) {
		ret.Push("")
	}
	lo.ForEach(third.Keys(), func(k string, _ int) {
		s, _ := third[k].(string)
		ret.Push(s)
	})
	if len(third) > 0 && len(self) > 0 {
		ret.Push("")
	}
	lo.ForEach(self.Keys(), func(k string, _ int) {
		s, _ := self[k].(string)
		ret.Push(s)
	})
	return ret.Slice
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
