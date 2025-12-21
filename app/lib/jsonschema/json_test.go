package jsonschema_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/util"
)

func TestSchemaJSONTrue(t *testing.T) {
	t.Parallel()
	s, err := util.FromJSONObj[*jsonschema.Schema]([]byte(" \n\ttrue\t "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatalf("expected trueSchemaData, got %p", s)
	}
	if !s.IsEmpty() {
		t.Fatalf("expected [true] schema, got %v", s)
	}
}

func TestSchemaJSONFalse(t *testing.T) {
	t.Parallel()
	s, err := util.FromJSONObj[*jsonschema.Schema]([]byte("false"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatalf("expected falseSchemaData, got %p", s)
	}
	if !s.IsEmptyExceptNot() {
		t.Fatalf("expected [false] schema, got %v", s)
	}
}

func TestSchemaJSONSimpleObject(t *testing.T) {
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
}

func TestSchemaJSONArray(t *testing.T) {
	t.Parallel()
	s, err := util.FromJSONObj[*jsonschema.Schema]([]byte("\n[]\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatalf("expected trueSchemaData, got %p", s)
	}
	if !s.IsEmpty() {
		t.Fatalf("expected empty schema, got %v", s)
	}
}
