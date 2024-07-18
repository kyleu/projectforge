package util

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	om := NewOrderedMap[int](true, 5)
	if !om.Lexical {
		t.Error("Expected Lexical to be true")
	}
	if cap(om.Order) != 5 {
		t.Errorf("expected Order capacity to be 5, got %d", cap(om.Order))
	}
	if len(om.Map) != 0 {
		t.Errorf("expected Map to be empty, got length %d", len(om.Map))
	}
}

func TestNewOMap(t *testing.T) {
	om := NewOMap[string]()
	if om.Lexical {
		t.Error("Expected Lexical to be false")
	}
	if len(om.Order) != 0 {
		t.Errorf("expected Order to be empty, got length %d", len(om.Order))
	}
	if len(om.Map) != 0 {
		t.Errorf("expected Map to be empty, got length %d", len(om.Map))
	}
}

func TestOrderedMap_Append(t *testing.T) {
	om := NewOMap[int]()
	om.Append("a", 1)
	om.Append("b", 2)

	if len(om.Order) != 2 || om.Order[0] != "a" || om.Order[1] != "b" {
		t.Errorf("unexpected Order: %v", om.Order)
	}
	if om.Map["a"] != 1 || om.Map["b"] != 2 {
		t.Errorf("unexpected Map: %v", om.Map)
	}
}

func TestOrderedMap_Set(t *testing.T) {
	om := NewOMap[int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("a", 3)

	if len(om.Order) != 2 || om.Order[0] != "a" || om.Order[1] != "b" {
		t.Errorf("unexpected Order: %v", om.Order)
	}
	if om.Map["a"] != 3 || om.Map["b"] != 2 {
		t.Errorf("unexpected Map: %v", om.Map)
	}
}

func TestOrderedMap_HasKey(t *testing.T) {
	om := NewOMap[int]()
	om.Set("a", 1)

	if !om.HasKey("a") {
		t.Error("Expected HasKey('a') to be true")
	}
	if om.HasKey("b") {
		t.Error("Expected HasKey('b') to be false")
	}
}

func TestOrderedMap_IndexOf(t *testing.T) {
	om := NewOMap[int]()
	om.Append("a", 1)
	om.Append("b", 2)

	if om.IndexOf("a") != 0 {
		t.Errorf("expected IndexOf('a') to be 0, got %d", om.IndexOf("a"))
	}
	if om.IndexOf("c") != -1 {
		t.Errorf("expected IndexOf('c') to be -1, got %d", om.IndexOf("c"))
	}
}

func TestOrderedMap_Get(t *testing.T) {
	om := NewOMap[int]()
	om.Set("a", 1)

	v, ok := om.Get("a")
	if !ok || v != 1 {
		t.Errorf("unexpected result for Get('a'): %v, %v", v, ok)
	}

	v, ok = om.Get("b")
	if ok || v != 0 {
		t.Errorf("unexpected result for Get('b'): %v, %v", v, ok)
	}
}

func TestOrderedMap_GetSimple(t *testing.T) {
	om := NewOMap[int]()
	om.Set("a", 1)

	if om.GetSimple("a") != 1 {
		t.Errorf("unexpected result for GetSimple('a'): %v", om.GetSimple("a"))
	}
}

func TestOrderedMap_Pairs(t *testing.T) {
	om := NewOMap[int]()
	om.Append("a", 1)
	om.Append("b", 2)

	pairs := om.Pairs()
	expected := []*OrderedPair[int]{
		{K: "a", V: 1},
		{K: "b", V: 2},
	}

	if !reflect.DeepEqual(pairs, expected) {
		t.Errorf("unexpected Pairs result: %v", pairs)
	}
}

func TestOrderedMap_Remove(t *testing.T) {
	om := NewOMap[int]()
	om.Append("a", 1)
	om.Append("b", 2)
	om.Remove("a")

	if len(om.Order) != 1 || om.Order[0] != "b" {
		t.Errorf("unexpected Order after Remove: %v", om.Order)
	}
	if len(om.Map) != 1 || om.Map["b"] != 2 {
		t.Errorf("unexpected Map after Remove: %v", om.Map)
	}
}

func TestOrderedMap_MarshalJSON(t *testing.T) {
	om := NewOMap[int]()
	om.Append("b", 2)
	om.Append("a", 1)

	b, err := json.Marshal(om)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `{"b":2,"a":1}`
	if string(b) != expected {
		t.Errorf("unexpected JSON: %s", string(b))
	}
}

func TestOrderedMap_UnmarshalJSON(t *testing.T) {
	input := `{"b":2,"a":1}`
	om := NewOMap[int]()

	err := json.Unmarshal([]byte(input), om)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(om.Order) != 2 || om.Order[0] != "b" || om.Order[1] != "a" {
		t.Errorf("unexpected Order: %v", om.Order)
	}
	if om.Map["b"] != 2 || om.Map["a"] != 1 {
		t.Errorf("unexpected Map: %v", om.Map)
	}
}

func TestOrderedMap_MarshalXML(t *testing.T) {
	om := NewOMap[int]()
	om.Append("b", 2)
	om.Append("a", 1)

	b, err := xml.Marshal(om)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := `<OrderedMap><b>2</b><a>1</a></OrderedMap>`
	if string(b) != expected {
		t.Errorf("unexpected XML: %s", string(b))
	}
}
