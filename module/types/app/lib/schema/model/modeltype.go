package model

import (
	"{{{ .Package }}}/app/util"
)

type Type struct {
	Key    string
	Title  string
	Plural string
	Icon   string
}

var (
	TypeEnum         = Type{Key: "enum", Title: "Enum", Plural: "Enums", Icon: "list"}
	TypeStruct       = Type{Key: "struct", Title: "Struct", Plural: "Structs", Icon: "struct"}
	TypeInterface    = Type{Key: "interface", Title: "Interface", Plural: "Interfaces", Icon: "list"}
	TypeUnion        = Type{Key: "union", Title: "Union", Plural: "Unions", Icon: "world"}
	TypeIntersection = Type{Key: "intersection", Title: "Intersection", Plural: "Intersections", Icon: "world"}
	TypeUnknown      = Type{Key: "unknown", Title: "Unknown", Plural: "Unknowns", Icon: "world"}
)

var AllModelTypes = []Type{TypeEnum, TypeStruct, TypeInterface, TypeUnion, TypeIntersection, TypeUnknown}

func modelTypeFromString(s string) (Type, error) {
	for _, t := range AllModelTypes {
		if t.Key == s {
			return t, nil
		}
	}
	return TypeUnknown, nil
}

func (t *Type) String() string {
	return t.Key
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(t.Key, false), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := util.FromJSON(data, &s); err != nil {
		return err
	}
	x, err := modelTypeFromString(s)
	if err != nil {
		return err
	}
	*t = x
	return nil
}
