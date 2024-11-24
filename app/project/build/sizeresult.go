package build

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type SizeResult struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size"`
}

func (s *SizeResult) String() string {
	return fmt.Sprintf("[%s:%d] %s", s.Type, s.Size, s.Name)
}

type SizeResults []*SizeResult

func (r SizeResults) TotalSize() int {
	return lo.SumBy(r, func(x *SizeResult) int {
		return x.Size
	})
}

type SizeResultMap map[string]SizeResults

func (s SizeResultMap) Strings() []string {
	ret := make([]string, 0, len(s))
	for k, v := range s {
		ret = append(ret, lo.Map(v, func(x *SizeResult, _ int) string {
			return k + ": " + x.String()
		})...)
	}
	return util.ArraySorted(ret)
}

func (s SizeResultMap) TotalStrings() []string {
	ret := make([]string, 0, len(s))
	counts := s.TotalCount()
	for k, v := range s.TotalSizes() {
		ret = append(ret, fmt.Sprintf("%s: [%d] bytes across [%d] files", k, v, counts[k]))
	}
	return util.ArraySorted(ret)
}

func (s SizeResultMap) Add(key string, v ...*SizeResult) {
	s[key] = append(s[key], v...)
}

func (s SizeResultMap) TotalSizes() map[string]int {
	return lo.MapValues(s, func(x SizeResults, _ string) int {
		return x.TotalSize()
	})
}

func (s SizeResultMap) TotalCount() map[string]int {
	return lo.MapValues(s, func(x SizeResults, _ string) int {
		return len(x)
	})
}

func (s SizeResultMap) Flatten() SizeResultMap {
	ret := SizeResultMap{}
	for k, v := range s {
		for _, x := range v {
			slashIdx := strings.LastIndex(k, "/")
			if slashIdx > -1 {
				pk := k[:slashIdx]
				xk := k[slashIdx+1:]
				if strings.Contains(xk, ".") {
					nk, nn := util.StringSplit(xk, '.', true)
					nk = clean(pk + "/" + nk)
					x.Name = nn + "." + x.Name
					ret.Add(nk, x)
				}
			} else {
				if strings.Contains(k, ".") {
					nk, nn := util.StringSplit(k, '.', true)
					nk = clean(nk)
					x.Name = nn + "." + x.Name
					ret.Add(nk, x)
				}
			}
		}
	}
	return ret
}

func clean(k string) string {
	if bracketIdx := strings.Index(k, "["); bracketIdx > -1 {
		k = k[:bracketIdx]
	}
	return k
}
