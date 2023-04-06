package types

import (
	"strings"

	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/util"
)

const KeyReference = "reference"

type Reference struct {
	Pkg util.Pkg `json:"pkg,omitempty"`
	K   string   `json:"k"`
}

func (x *Reference) Key() string {
	return KeyReference
}

func (x *Reference) Sortable() bool {
	return true
}

func (x *Reference) Scalar() bool {
	return false
}

func (x *Reference) String() string {
	return "ref:" + strings.Join(x.Path(), ".")
}

func (x *Reference) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Reference) Path() util.Pkg {
	ret := slices.Clone(x.Pkg)
	return append(ret, x.K)
}

func (x *Reference) Default(string) any {
	return ""
}

func NewReference() *Wrapped {
	return Wrap(&Reference{})
}

func NewReferenceArgs(pkg util.Pkg, k string) *Wrapped {
	return Wrap(&Reference{Pkg: pkg, K: k})
}
