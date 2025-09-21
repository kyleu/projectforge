package icons

import (
	"fmt"
	"strings"
	"sync"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Library struct {
	Icons map[string]*Icon  `json:"icons"`
	Keys  map[string]string `json:"keys,omitzero"`
}

func NewLibrary(icons ...*Icon) *Library {
	ret := &Library{Icons: map[string]*Icon{}, Keys: map[string]string{}}
	lo.ForEach(icons, func(x *Icon, _ int) {
		ret.AddIcon(x)
	})
	return ret
}

func (l *Library) AddIcon(bi *Icon) {
	l.Icons[bi.Key] = bi
	l.Keys[bi.Title] = bi.Key
	for _, x := range bi.Aliases {
		l.Keys[x] = bi.Key
	}
}

func (l *Library) SortedKeys() []string {
	return util.MapKeysSorted(l.Icons)
}

func (l *Library) HTML(key string) string {
	key = strings.TrimPrefix(key, "brand-")
	ret, ok := l.Icons[key]
	if !ok {
		return fmt.Sprintf("<!-- unknown brand icon [%s] -->", key)
	}
	return ret.HTML("brand-")
}

var (
	brandLibCache *Library
	brandLibMu    sync.Mutex
)

func BrandLibrary() *Library {
	if brandLibCache == nil {
		brandLibMu.Lock()
		defer brandLibMu.Unlock()
		if brandLibCache == nil {
			brandLibCache = NewLibrary(BrandIcons...)
		}
	}
	return brandLibCache
}
