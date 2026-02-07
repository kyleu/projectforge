package types_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

type wrappedCheck func(t *testing.T, w *types.Wrapped)

type wrappedDefaultTest struct {
	name  string
	input string
	check wrappedCheck
}

var wrappedDefaultTests = []wrappedDefaultTest{
	{name: "boolean-alias", input: `"boolean"`, check: checkWrappedBooleanAlias},
	{name: "map-defaults", input: `"map"`, check: checkWrappedMapDefaults},
	{name: "list-defaults", input: `"list"`, check: checkWrappedListDefaults},
	{name: "orderedmap-defaults", input: `"orderedMap"`, check: checkWrappedOrderedMapDefaults},
	{name: "option-defaults", input: `"option"`, check: checkWrappedOptionDefaults},
	{name: "range-defaults", input: `"range"`, check: checkWrappedRangeDefaults},
	{name: "set-defaults", input: `"set"`, check: checkWrappedSetDefaults},
	{name: "method-defaults", input: `"method"`, check: checkWrappedMethodDefaults},
}

func TestWrappedMarshalJSONDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value *types.Wrapped
		want  string
	}{
		{name: "int-default", value: types.NewInt(0), want: `"int"`},
		{name: "string-default", value: types.NewString(), want: `"string"`},
		{name: "json-default", value: types.NewJSON(), want: `"json"`},
		{name: "map-default", value: types.NewStringKeyedMap(), want: `"map"`},
		{name: "list-default", value: types.NewList(types.NewAny()), want: `"list"`},
		{name: "orderedmap-default", value: types.NewOrderedMap(types.NewAny()), want: `"orderedMap"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := util.ToJSONCompact(tt.value)
			if got != tt.want {
				t.Errorf("MarshalJSON = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestWrappedMarshalJSONWithArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value *types.Wrapped
		want  string
	}{
		{name: "int-bits", value: types.NewInt(32), want: `{"k":"int","t":{"bits":32}}`},
		{name: "json-object", value: types.NewJSONArgs(true, false), want: `{"k":"json","t":{"obj":true}}`},
		{name: "map-string-string", value: types.NewMap(types.NewString(), types.NewString()), want: `{"k":"map","t":{"k":"string","v":"string"}}`},
		{name: "list-string", value: types.NewList(types.NewString()), want: `{"k":"list","t":{"v":"string"}}`},
		{name: "orderedmap-string", value: types.NewOrderedMap(types.NewString()), want: `{"k":"orderedMap","t":{"v":"string"}}`},
		{name: "option-string", value: types.NewOption(types.NewString()), want: `{"k":"option","t":{"v":"string"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := util.ToJSONCompact(tt.value)
			if got != tt.want {
				t.Errorf("MarshalJSON = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestWrappedUnmarshalJSONDefaults(t *testing.T) {
	t.Parallel()

	for _, tt := range wrappedDefaultTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var w types.Wrapped
			if err := util.FromJSON([]byte(tt.input), &w); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tt.check(t, &w)
		})
	}
}

func checkWrappedBooleanAlias(t *testing.T, w *types.Wrapped) {
	t.Helper()
	if w.Key() != types.KeyBool {
		t.Errorf("Key() = %s, want %s", w.Key(), types.KeyBool)
	}
	if types.TypeAs[*types.Bool](w) == nil {
		t.Fatalf("expected Bool type, got %T", w.T)
	}
}

func checkWrappedMapDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	m := types.TypeAs[*types.Map](w)
	if m == nil {
		t.Fatalf("expected Map type, got %T", w.T)
	}
	if m.K == nil || m.V == nil {
		t.Fatalf("expected map key/value to be set, got K=%v V=%v", m.K, m.V)
	}
	if m.K.Key() != types.KeyString {
		t.Errorf("map key = %s, want %s", m.K.Key(), types.KeyString)
	}
	if m.V.Key() != types.KeyAny {
		t.Errorf("map value = %s, want %s", m.V.Key(), types.KeyAny)
	}
}

func checkWrappedListDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	l := types.TypeAs[*types.List](w)
	if l == nil {
		t.Fatalf("expected List type, got %T", w.T)
	}
	if l.V == nil {
		t.Fatalf("expected list value to be set")
	}
	if l.V.Key() != types.KeyAny {
		t.Errorf("list value = %s, want %s", l.V.Key(), types.KeyAny)
	}
}

func checkWrappedOrderedMapDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	m := types.TypeAs[*types.OrderedMap](w)
	if m == nil {
		t.Fatalf("expected OrderedMap type, got %T", w.T)
	}
	if m.V == nil {
		t.Fatalf("expected ordered map value to be set")
	}
	if m.V.Key() != types.KeyAny {
		t.Errorf("ordered map value = %s, want %s", m.V.Key(), types.KeyAny)
	}
}

func checkWrappedOptionDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	o := types.TypeAs[*types.Option](w)
	if o == nil {
		t.Fatalf("expected Option type, got %T", w.T)
	}
	if o.V == nil {
		t.Fatalf("expected option value to be set")
	}
	if o.V.Key() != types.KeyAny {
		t.Errorf("option value = %s, want %s", o.V.Key(), types.KeyAny)
	}
}

func checkWrappedRangeDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	r := types.TypeAs[*types.Range](w)
	if r == nil {
		t.Fatalf("expected Range type, got %T", w.T)
	}
	if r.V == nil {
		t.Fatalf("expected range value to be set")
	}
	if r.V.Key() != types.KeyAny {
		t.Errorf("range value = %s, want %s", r.V.Key(), types.KeyAny)
	}
}

func checkWrappedSetDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	s := types.TypeAs[*types.Set](w)
	if s == nil {
		t.Fatalf("expected Set type, got %T", w.T)
	}
	if s.V == nil {
		t.Fatalf("expected set value to be set")
	}
	if s.V.Key() != types.KeyAny {
		t.Errorf("set value = %s, want %s", s.V.Key(), types.KeyAny)
	}
}

func checkWrappedMethodDefaults(t *testing.T, w *types.Wrapped) {
	t.Helper()
	m := types.TypeAs[*types.Method](w)
	if m == nil {
		t.Fatalf("expected Method type, got %T", w.T)
	}
	if m.Ret == nil {
		t.Fatalf("expected method return type to be set")
	}
	if m.Ret.Key() != types.KeyAny {
		t.Errorf("method return = %s, want %s", m.Ret.Key(), types.KeyAny)
	}
}

func TestWrappedJSONRoundTrip(t *testing.T) {
	t.Parallel()

	src := types.NewMap(types.NewString(), types.NewList(types.NewString()))
	var dst types.Wrapped
	if err := util.CycleJSON(src, &dst); err != nil {
		t.Fatalf("cycle error: %v", err)
	}
	if !dst.Equals(src) {
		t.Errorf("round-trip mismatch: %s != %s", dst.String(), src.String())
	}
}
