package export

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

type Type struct {
	Keys    []string
	Go      string
	Imports GoImports
}

func newType(keys []string, g string, imps ...*GoImport) *Type {
	return &Type{Keys: keys, Go: g, Imports: imps}
}

var (
	TypeInt       = newType([]string{"int"}, "int")
	TypeString    = newType([]string{"string", "text"}, "string")
	TypeTimestamp = newType([]string{"timestamp", "datetime"}, "time.Time", &GoImport{Type: ImportTypeInternal, Value: "time"})
	TypeUUID      = newType([]string{"uuid"}, "uuid.UUID", &GoImport{Type: ImportTypeExternal, Value: "github.com/google/uuid"})
)

var (
	AllTypes = []*Type{TypeInt, TypeString, TypeTimestamp, TypeUUID}
)

type Types []*Type

func (t Types) GoKeys() []string {
	ret := make([]string, 0, len(t))
	for _, x := range t {
		ret = append(ret, x.Go)
	}
	return ret
}

func (t Types) MaxGoKeyLength() int {
	return util.StringArrayMaxLength(t.GoKeys())
}

func (t Types) Imports() GoImports {
	ret := GoImports{}
	for _, x := range t {
		ret = ret.Add(x.Imports...)
	}
	return ret
}

func TypeFromString(s string) (*Type, error) {
	for _, t := range AllTypes {
		for _, k := range t.Keys {
			if k == s {
				return t, nil
			}
		}
	}
	return nil, errors.New("No export type available with key [" + s + "]")
}

func (t *Type) String() string {
	return t.Keys[0]
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(t.Keys[0], false), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := util.FromJSON(data, &s); err != nil {
		return err
	}
	x, err := TypeFromString(s)
	if err != nil {
		return errors.Wrapf(err, "no export type available with key [%s]", s)
	}

	*t = *x
	return nil
}
