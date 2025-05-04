package util

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

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
	return StringJoin(lo.Map(k, func(x *KeyVal[T], _ int) string {
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

func (k *KeyTypeDesc) String() string {
	return fmt.Sprintf("%s (%s)", k.Key, k.Type)
}

func (k *KeyTypeDesc) Array(key string) []string {
	return []string{strings.ReplaceAll("`"+k.Key+"`", "{key}", key), k.Type, strings.ReplaceAll(k.Description, "{key}", key)}
}

func (k *KeyTypeDesc) Matches(x *KeyTypeDesc) bool {
	return k.Key == x.Key
}

type KeyTypeDescs []*KeyTypeDesc

func (k KeyTypeDescs) Strings() []string {
	return lo.Map(k, func(x *KeyTypeDesc, _ int) string {
		return x.String()
	})
}

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
