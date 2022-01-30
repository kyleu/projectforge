// Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
	"strings"
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
	argStrings := make([]string, 0, len(x.Args))
	for _, arg := range x.Args {
		argStrings = append(argStrings, arg.String())
	}
	return fmt.Sprintf("fn(%s) %s", strings.Join(argStrings, ", "), x.Ret.String())
}

func (x *Method) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *Method) Default(key string) interface{} {
	return key + "()"
}

func NewMethod(ret *Wrapped) *Wrapped {
	return Wrap(&Method{Ret: ret})
}
