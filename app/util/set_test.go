package util

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNewSet(t *testing.T) {
	s := NewSet[int](5)
	if len(s.Map) != 0 {
		t.Errorf("Expected empty set, got %v", s.Map)
	}
}

func TestSet(t *testing.T) {
	s := NewSet[int]()
	s.Set(1)
	s.Set(2)
	s.Set(1)
	if len(s.Map) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(s.Map))
	}
}

func TestContains(t *testing.T) {
	s := NewSet[string]()
	s.Set("a")
	if !s.Contains("a") {
		t.Error("Expected set to contain 'a'")
	}
	if s.Contains("b") {
		t.Error("Expected set to not contain 'b'")
	}
}

func TestRemove(t *testing.T) {
	s := NewSet[int]()
	s.Set(1)
	s.Set(2)
	s.Remove(1)
	if s.Contains(1) {
		t.Error("Expected 1 to be removed")
	}
	if !s.Contains(2) {
		t.Error("Expected 2 to still be present")
	}
}

func TestEntries(t *testing.T) {
	s := NewSet[int]()
	s.Set(3)
	s.Set(1)
	s.Set(2)
	entries := s.Entries()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(entries, expected) {
		t.Errorf("Expected %v, got %v", expected, entries)
	}
}

func TestClone(t *testing.T) {
	s1 := NewSet[int]()
	s1.Set(1)
	s1.Set(2)
	s2 := s1.Clone()
	if !reflect.DeepEqual(s1.Map, s2.Map) {
		t.Error("Cloned set not equal to original")
	}
	s2.Set(3)
	if s1.Contains(3) {
		t.Error("Original set should not be affected by changes to clone")
	}
}

func TestMarshalYAML(t *testing.T) {
	s := NewSet[string]()
	s.Set("c")
	s.Set("a")
	s.Set("b")
	yaml, err := s.MarshalYAML()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(yaml, expected) {
		t.Errorf("Expected %v, got %v", expected, yaml)
	}
}

func TestMarshalXML(t *testing.T) {
	s := NewSet[int]()
	s.Set(1)
	s.Set(2)
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	err := s.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "set"}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := `<set><int>1</int><int>2</int></set>`
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}

func TestUnmarshalJSON(t *testing.T) {
	jsonData := `[1, 2, 3]`
	s := NewSet[int]()
	err := s.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(s.Entries(), expected) {
		t.Errorf("Expected %v, got %v", expected, s.Entries())
	}
}

func TestMarshalJSON(t *testing.T) {
	s := NewSet[string]()
	s.Set("a")
	s.Set("b")
	s.Set("c")
	js, err := s.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := ToJSON([]string{"a", "b", "c"})
	if string(js) != expected {
		t.Errorf("Expected %s, got %s", expected, string(js))
	}
}
