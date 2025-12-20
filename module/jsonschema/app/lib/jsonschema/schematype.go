package jsonschema

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

//nolint:lll
var (
	SchemaTypeObject  = SchemaType{Key: "object", Name: "Object", Icon: "cog"}
	SchemaTypeEnum    = SchemaType{Key: "enum", Name: "Enum", Icon: "cog"}
	SchemaTypeArray   = SchemaType{Key: "array", Name: "Array", Icon: "cog"}
	SchemaTypeRef     = SchemaType{Key: "ref", Name: "Ref", Icon: "cog"}
	SchemaTypeEmpty   = SchemaType{Key: "empty", Name: "Empty", Icon: "cog"}
	SchemaTypeNot     = SchemaType{Key: "not", Name: "Not", Icon: "cog"}
	SchemaTypeString  = SchemaType{Key: "string", Name: "String", Icon: "cog"}
	SchemaTypeBoolean = SchemaType{Key: "boolean", Name: "Boolean", Icon: "cog"}
	SchemaTypeInteger = SchemaType{Key: "integer", Name: "Integer", Icon: "cog"}
	SchemaTypeNumber  = SchemaType{Key: "number", Name: "Number", Icon: "cog"}
	SchemaTypeDate    = SchemaType{Key: "date", Name: "Date", Icon: "cog"}
	SchemaTypeNull    = SchemaType{Key: "null", Name: "Null", Icon: "cog"}
	SchemaTypeUnknown = SchemaType{Key: "unknown", Name: "Unknown", Icon: "cog"}

	AllSchemaTypes = SchemaTypes{SchemaTypeObject, SchemaTypeEnum, SchemaTypeArray, SchemaTypeRef, SchemaTypeEmpty, SchemaTypeNot, SchemaTypeString, SchemaTypeBoolean, SchemaTypeInteger, SchemaTypeNumber, SchemaTypeDate, SchemaTypeNull, SchemaTypeUnknown}
)

type SchemaType struct {
	Key         string
	Name        string
	Description string
	Icon        string
}

func (s SchemaType) String() string {
	return s.Key
}

func (s SchemaType) NameSafe() string {
	if s.Name != "" {
		return s.Name
	}
	return s.String()
}

func (s SchemaType) Matches(xx SchemaType) bool {
	return s.Key == xx.Key
}

func (s SchemaType) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(s.Key, false), nil
}

func (s *SchemaType) UnmarshalJSON(data []byte) error {
	key, err := util.FromJSONString(data)
	if err != nil {
		return err
	}
	*s = AllSchemaTypes.Get(key, nil)
	return nil
}

func (s SchemaType) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.EncodeElement(s.Key, start)
}

func (s *SchemaType) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var key string
	if err := dec.DecodeElement(&key, &start); err != nil {
		return err
	}
	*s = AllSchemaTypes.Get(key, nil)
	return nil
}

func (s SchemaType) Value() (driver.Value, error) {
	return s.Key, nil
}

func (s *SchemaType) Scan(value any) error {
	if value == nil {
		return nil
	}
	if converted, err := driver.String.ConvertValue(value); err == nil {
		if str, err := util.Cast[string](converted); err == nil {
			*s = AllSchemaTypes.Get(str, nil)
			return nil
		}
	}
	return errors.Errorf("failed to scan SchemaType enum from value [%v]", value)
}

func SchemaTypeParse(logger util.Logger, keys ...string) SchemaTypes {
	if len(keys) == 0 {
		return nil
	}
	return lo.Map(keys, func(x string, _ int) SchemaType {
		return AllSchemaTypes.Get(x, logger)
	})
}

type SchemaTypes []SchemaType

func (s SchemaTypes) Keys() []string {
	return lo.Map(s, func(x SchemaType, _ int) string {
		return x.Key
	})
}

func (s SchemaTypes) Strings() []string {
	return lo.Map(s, func(x SchemaType, _ int) string {
		return x.String()
	})
}

func (s SchemaTypes) NamesSafe() []string {
	return lo.Map(s, func(x SchemaType, _ int) string {
		return x.NameSafe()
	})
}

func (s SchemaTypes) Help() string {
	return "Available schema type options: [" + util.StringJoin(s.Strings(), ", ") + "]"
}

func (s SchemaTypes) Get(key string, logger util.Logger) SchemaType {
	for _, value := range s {
		if strings.EqualFold(value.Key, key) {
			return value
		}
	}
	if key == "" {
		return SchemaTypeUnknown
	}
	msg := fmt.Sprintf("unable to find [SchemaType] with key [%s]", key)
	if logger != nil {
		logger.Warn(msg)
	}
	return SchemaTypeUnknown
}

func (s SchemaTypes) GetByName(name string, logger util.Logger) SchemaType {
	for _, value := range s {
		if strings.EqualFold(value.Name, name) {
			return value
		}
	}
	if name == "" {
		return SchemaTypeUnknown
	}
	msg := fmt.Sprintf("unable to find [SchemaType] with name [%s]", name)
	if logger != nil {
		logger.Warn(msg)
	}
	return SchemaTypeUnknown
}

func (s SchemaTypes) Random() SchemaType {
	return util.RandomElement(s)
}
