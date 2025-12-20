package jsonschema_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/util"
)

func TestSchemaDataBoolAndObject(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()
		s, err := util.FromJSONObj[*jsonschema.Schema]([]byte(" \n\ttrue\t "))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s == nil {
			t.Fatalf("expected trueSchemaData, got %p", s)
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		s, err := util.FromJSONObj[*jsonschema.Schema]([]byte("false"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s == nil {
			t.Fatalf("expected falseSchemaData, got %p", s)
		}
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		s, err := util.FromJSONObj[*jsonschema.Schema]([]byte(`{"title":"ok"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s == nil {
			t.Fatalf("expected a new schemaData, got shared pointer %p", s)
		}
		if s.Title != "ok" {
			t.Fatalf("expected title [ok], got [%s]", s.Title)
		}
	})

	t.Run("invalid_array", func(t *testing.T) {
		t.Parallel()
		_, err := util.FromJSONObj[*jsonschema.Schema]([]byte(`[]`))
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}
