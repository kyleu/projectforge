// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
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
	return lo.Associate(k, func(x *KeyValInt) (string, int) {
		return x.Key, x.Count
	})
}

func (k KeyValInts) String() string {
	return strings.Join(lo.Map(k, func(x *KeyValInt, _ int) string {
		return x.String()
	}), ", ")
}

type KeyVal[T any] struct {
	Key string `json:"key" db:"key"`
	Val T      `json:"val" db:"val"`
}

func (k KeyVal[T]) String() string {
	return fmt.Sprintf("%s: %v", k.Key, k.Val)
}

type KeyVals[T any] []*KeyVal[T]

func (k KeyVals[T]) ToMap() map[string]T {
	return lo.Associate(k, func(x *KeyVal[T]) (string, T) {
		return x.Key, x.Val
	})
}

func (k KeyVals[T]) String() string {
	return strings.Join(lo.Map(k, func(x *KeyVal[T], _ int) string {
		return x.String()
	}), ", ")
}

func (k KeyVals[T]) Values() []T {
	return lo.Map(k, func(x *KeyVal[T], _ int) T {
		return x.Val
	})
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

func (k KeyTypeDescs) Sort() KeyTypeDescs {
	slices.SortFunc(k, func(l *KeyTypeDesc, r *KeyTypeDesc) int {
		return cmp.Compare(strings.ToLower(l.Key), strings.ToLower(r.Key))
	})
	return k
}

func (k KeyTypeDescs) Array(key string) [][]string {
	return lo.Map(k.Sort(), func(x *KeyTypeDesc, _ int) []string {
		return x.Array(key)
	})
}
