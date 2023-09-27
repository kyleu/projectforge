// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Argument struct {
	Key  string   `json:"key"`
	Type *Wrapped `json:"type"`
}

func (a Argument) String() string {
	return fmt.Sprintf("%s %s", a.Key, a.Type.String())
}

type Arguments []Argument

const KeyMethod = "method"

type Method struct {
	Args Arguments `json:"args,omitempty"`
	Ret  *Wrapped  `json:"ret,omitempty"`
}

var _ Type = (*Method)(nil)

func (x *Method) Key() string {
	return KeyMethod
}

func (x *Method) Sortable() bool {
	for _, a := range x.Args {
		if !a.Type.Sortable() {
			return false
		}
	}
	return x.Ret.Sortable()
}

func (x *Method) Scalar() bool {
	return false
}

func (x *Method) String() string {
	argStrings := lo.Map(x.Args, func(arg Argument, _ int) string {
		return arg.String()
	})
	return fmt.Sprintf("fn(%s) %s", strings.Join(argStrings, ", "), x.Ret.String())
}

func (x *Method) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Method) Default(key string) any {
	return key + "()"
}

func NewMethod(ret *Wrapped) *Wrapped {
	return Wrap(&Method{Ret: ret})
}
