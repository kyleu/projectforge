// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"strings"

	"golang.org/x/exp/slices"
)

type KeyValInt struct {
	Key   string `json:"key" db:"key"`
	Count int    `json:"val" db:"val"`
}

type KeyTypeDesc struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (k *KeyTypeDesc) Array() []string {
	return []string{strings.ReplaceAll("`"+k.Key+"`", "{key}", AppKey), k.Type, strings.ReplaceAll(k.Description, "{key}", AppKey)}
}

type KeyTypeDescs []*KeyTypeDesc

func (k KeyTypeDescs) Sort() {
	slices.SortFunc(k, func(l *KeyTypeDesc, r *KeyTypeDesc) bool {
		return strings.ToLower(l.Key) < strings.ToLower(r.Key)
	})
}

func (k KeyTypeDescs) Array() [][]string {
	k.Sort()
	ret := make([][]string, 0, len(k))
	for _, x := range k {
		ret = append(ret, x.Array())
	}
	return ret
}
