//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"{{{ .Package }}}/app/util"
)

type xmlThing struct {
	XMLName xml.Name `xml:"thing"`
	Name    string   `xml:"name"`
}

func TestTOMLAndXML(t *testing.T) {
	src := map[string]any{"a": 1, "b": "two"}
	msg := util.ToTOML(src)
	if msg == "" {
		t.Fatalf("expected toml output")
	}
	var out util.ValueMap
	if err := util.FromTOML([]byte(msg), &out); err != nil || out["b"] != "two" {
		t.Fatalf("unexpected toml parse: %v %v", out, err)
	}

	x := xmlThing{Name: "hi"}
	xmlStr, err := util.ToXML(x)
	if err != nil || !strings.Contains(xmlStr, "<thing>") {
		t.Fatalf("unexpected xml output: %v %q", err, xmlStr)
	}
	var decoded xmlThing
	if err := util.FromXML([]byte(xmlStr), &decoded); err != nil || decoded.Name != "hi" {
		t.Fatalf("unexpected xml decode: %v %v", decoded, err)
	}

	compact, err := util.ToXMLCompact([]string{"a", "b"})
	if err != nil || !strings.Contains(compact, "<string>") {
		t.Fatalf("unexpected compact xml: %v %q", err, compact)
	}

	if err := util.CycleXML(xmlThing{Name: "ok"}, &decoded); err != nil || decoded.Name != "ok" {
		t.Fatalf("unexpected CycleXML result: %v %v", decoded, err)
	}

	if err := util.FromXMLStrict([]byte(xmlStr), &decoded); err != nil {
		t.Fatalf("unexpected FromXMLStrict error: %v", err)
	}

	buf := bytes.NewBufferString(xmlStr)
	if err := util.FromXMLReader(buf, &decoded); err != nil {
		t.Fatalf("unexpected FromXMLReader error: %v", err)
	}
}
