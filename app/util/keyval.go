// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type KeyValInt struct {
	Key   string `json:"key" db:"key"`
	Count int    `json:"val" db:"val"`
}

func (k KeyValInt) String() string {
	return fmt.Sprintf("%s: %d", k.Key, k.Count)
}

type KeyValInts []*KeyValInt

func (k KeyValInts) ToMap() map[string]int {
	ret := make(map[string]int, len(k))
	for _, x := range k {
		ret[x.Key] = x.Count
	}
	return ret
}

func (k KeyValInts) String() string {
	ret := make([]string, 0, len(k))
	for _, x := range k {
		ret = append(ret, x.String())
	}
	return strings.Join(ret, ", ")
}

type KeyTypeDesc struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (k *KeyTypeDesc) Array(key string) []string {
	return []string{strings.ReplaceAll("`"+k.Key+"`", "{key}", key), k.Type, strings.ReplaceAll(k.Description, "{key}", key)}
}

type KeyTypeDescs []*KeyTypeDesc

func (k KeyTypeDescs) Sort() {
	slices.SortFunc(k, func(l *KeyTypeDesc, r *KeyTypeDesc) bool {
		return strings.ToLower(l.Key) < strings.ToLower(r.Key)
	})
}

func (k KeyTypeDescs) Array(key string) [][]string {
	k.Sort()
	ret := make([][]string, 0, len(k))
	for _, x := range k {
		ret = append(ret, x.Array(key))
	}
	return ret
}
