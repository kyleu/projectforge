package jsonschema_test

import (
	"testing"

	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/util"
)

func TestSchemaType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		sch  *jsonschema.Schema
		typ  string
		json string
	}{
		{name: "empty", sch: &jsonschema.Schema{}, typ: "empty", json: "true"},
		{name: "ref", sch: jsonschema.NewRefSchema("test"), typ: "ref", json: `{"$ref":"ref:test"}`},
		{name: "true", sch: jsonschema.NewTrueSchema(), typ: "empty", json: "true"},
		{name: "false", sch: jsonschema.NewFalseSchema(), typ: "not", json: "false"},
		{name: "string", sch: jsonschema.NewSchema(jsonschema.SchemaTypeString), typ: "string", json: `{"type":"string"}`},
		{name: "number", sch: jsonschema.NewSchema(jsonschema.SchemaTypeNumber), typ: "number", json: `{"type":"number"}`},
		{name: "integer", sch: jsonschema.NewSchema(jsonschema.SchemaTypeInteger), typ: "integer", json: `{"type":"integer"}`},
		{name: "boolean", sch: jsonschema.NewSchema(jsonschema.SchemaTypeBoolean), typ: "boolean", json: `{"type":"boolean"}`},
		{name: "array", sch: jsonschema.NewSchema(jsonschema.SchemaTypeArray), typ: "array", json: `{"type":"array"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if st := tt.sch.DetectSchemaType(); st.Key != tt.typ {
				t.Errorf("DetectSchemaType() = %v, want %v", st.Key, tt.typ)
			}
			if js := util.ToJSONCompact(tt.sch); js != tt.json {
				t.Errorf("ToJSONCompact() = %v, want %v", js, tt.json)
			}
		})
	}
}
