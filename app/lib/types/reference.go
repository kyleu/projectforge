package types

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/util"
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
	return "ref:" + util.StringJoin(x.Path(), "/")
}

func (x *Reference) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Reference) Path() util.Pkg {
	ret := util.ArrayCopy(x.Pkg)
	return append(ret, strings.TrimPrefix(x.K, "*"))
}

func (x *Reference) LastRef(includePkg bool) string {
	if len(x.Pkg) == 0 || !includePkg {
		return x.K
	}
	prefix := util.Choose(strings.HasPrefix(x.K, "*"), "*", "")
	return fmt.Sprintf("%s%s.%s", prefix, x.Pkg.Last(), strings.TrimPrefix(x.K, "*"))
}

func (x *Reference) LastAddr(includePkg bool) string {
	if len(x.Pkg) == 0 || !includePkg {
		return strings.ReplaceAll(x.K, "*", "&")
	}
	prefix := util.Choose(strings.HasPrefix(x.K, "*"), "&", "")
	return fmt.Sprintf("%s%s.%s", prefix, x.Pkg.Last(), strings.TrimPrefix(x.K, "*"))
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

func NewReferencePath(pth string, isPointer bool) *Wrapped {
	pth = strings.TrimPrefix(pth, "ref:")
	pth = strings.TrimSuffix(pth, ".json")
	pth = strings.TrimSuffix(pth, ".schema")
	parts := strings.Split(pth, "/")
	pkg := parts[:len(parts)-1]
	ref := parts[len(parts)-1]
	if isPointer && !strings.HasPrefix(ref, "*") {
		ref = "*" + ref
	}
	return NewReferenceArgs(pkg, ref)
}
