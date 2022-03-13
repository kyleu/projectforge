package util

import (
	"fmt"
	"sort"
	"strings"
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

func (k *KeyTypeDesc) String() string {
	return fmt.Sprintf("%s (%s): %s", k.Key, k.Type, k.Description)
}

type KeyTypeDescs []*KeyTypeDesc

func (k KeyTypeDescs) Sort() {
	sort.Slice(k, func(i, j int) bool {
		l, r := k[i], k[j]
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

