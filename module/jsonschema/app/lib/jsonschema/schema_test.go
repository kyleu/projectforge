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
		{name: "ref", sch: jsonschema.NewRefSchema("test"), typ: "ref", json: `{"$ref":"test"}`},
		{name: "true", sch: jsonschema.NewTrueSchema(), typ: "empty", json: "true"},
		{name: "false", sch: jsonschema.NewFalseSchema(), typ: "not", json: "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if st := tt.sch.SchemaType(); st.Key != tt.typ {
				t.Errorf("SchemaType() = %v, want %v", st.Key, tt.typ)
			}
			if js := util.ToJSONCompact(tt.sch); js != tt.json {
				t.Errorf("ToJSONCompact() = %v, want %v", js, tt.json)
			}
		})
	}
}
