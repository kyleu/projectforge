//go:build test_all || !func_test

// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util_test

import (
	"testing"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

type StringTest struct {
	TestValue string
	Expected  string
	Type      string
}

func (t *StringTest) Test() error {
	var res string
	switch t.Type {
	case "camel":
		res = util.StringToCamel(t.TestValue)
	case "snake":
		res = util.StringToSnake(t.TestValue)
	default:
		res = util.StringToTitle(t.TestValue)
	}
	if res != t.Expected {
		return errors.Errorf("StringToTitle returned [%s], not expected [%s]", res, t.Expected)
	}
	return nil
}

var titleTests = []*StringTest{
	{TestValue: "SimpleCamelCase", Expected: "Simple Camel Case"},
	{TestValue: "CSVFilesAreCoolButTXTRules", Expected: "CSV Files Are Cool But TXT Rules"},
	{TestValue: "MediaTypes", Expected: "Media Types"},
	{TestValue: "ID", Expected: "ID"},
	{TestValue: "Id", Expected: "ID"},
	{TestValue: "id", Expected: "ID"},
	{TestValue: "bigXMLBlob", Expected: "Big XML Blob"},
	{TestValue: "bigXMLBlob", Expected: "BigXMLBlob", Type: "camel"},
	{TestValue: "bigXMLBlob", Expected: "big_xml_blob", Type: "snake"},
	// {TestValue: "bigBlobXMLs", Expected: "big_blob_xmls", Type: "snake"},
	// {TestValue: "SetOfIDs", Expected: "Set of IDs"},
	// {TestValue: "Set Of IDs", Expected: "set_of_ids", Type: "snake"},
}

func TestToTitle(t *testing.T) {
	t.Parallel()
	_ = util.InitAcronyms()
	for _, test := range titleTests {
		err := test.Test()
		if err != nil {
			t.Error(errors.Wrapf(err, "test [%s] failed [%s] check", test.TestValue, test.Type))
		}
	}
}

func TestSubstringBetween(t *testing.T) {
	t.Parallel()
	test := func(s string, l string, r string, expected string) {
		if res := util.StringSubstringBetween(s, l, r); res != expected {
			t.Errorf("invalid result for substring [%s] between [%s] and [%s]: %s", s, l, r, res)
		}
	}
	test("abc", "a", "c", "b")
	test("looooooongassstring", "looooooong", "string", "ass")
	test("thishasstuffinthemiddle", "has", "the", "stuffin")
	test("missingprefix", "invalid", "prefix", "")
	test("missingsuffix", "missing", "invalid", "suffix")
}

func TestReplaceBetween(t *testing.T) {
	t.Parallel()
	test := func(s string, l string, r string, replacement string, expected string) {
		res, err := util.StringReplaceBetween(s, l, r, replacement)
		if err != nil {
			t.Error(err)
		}
		if res != expected {
			t.Errorf("invalid result for substring [%s] between [%s] and [%s]: %s", s, l, r, res)
		}
	}
	test("abc", "a", "c", "x", "axc")
	test("ApplePearKiwi", "Apple", "Kiwi", "Strawberry", "AppleStrawberryKiwi")
	test("thishasstuffinthemiddle", "has", "the", "thingsin", "thishasthingsinthemiddle")
}

type StringArgs struct {
	Input string
	Size  int
	Chr   rune
}

func TestPadRight(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Name     string
		Args     StringArgs
		Expected string
	}{
		{Name: "normal args", Args: StringArgs{Input: "foo", Size: 5, Chr: '界'}, Expected: "foo界界"},
		{Name: "nothing to pad, size is equal to input len", Args: StringArgs{Input: "привет", Size: 6, Chr: '界'}, Expected: "привет"},
		{Name: "nothing to pad, size is smaller than input len", Args: StringArgs{Input: "привет", Size: 4, Chr: '界'}, Expected: "привет"},
		{Name: "negative size, input is left intact", Args: StringArgs{Input: "привет", Size: -10, Chr: '界'}, Expected: "привет"},
		{Name: "empty input", Args: StringArgs{Input: "", Size: 3, Chr: '界'}, Expected: "界界界"},
	}
	for _, tCase := range cases {
		c := tCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			res := util.StringPadRight(c.Args.Input, c.Args.Size, c.Args.Chr)
			if res != c.Expected {
				t.Errorf("unexpected result %s, expected %s", res, c.Expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Name     string
		Args     StringArgs
		Expected string
	}{
		{Name: "normal args", Args: StringArgs{Input: "foo", Size: 5, Chr: '界'}, Expected: "界界foo"},
		{Name: "nothing to pad, size is equal to input len", Args: StringArgs{Input: "привет", Size: 6, Chr: '界'}, Expected: "привет"},
		{Name: "nothing to pad, size is smaller than input len", Args: StringArgs{Input: "привет", Size: 4, Chr: '界'}, Expected: "привет"},
		{Name: "negative size, input is left intact", Args: StringArgs{Input: "привет", Size: -10, Chr: '界'}, Expected: "привет"},
		{Name: "empty input", Args: StringArgs{Input: "", Size: 3, Chr: '界'}, Expected: "界界界"},
	}
	for _, tCase := range cases {
		c := tCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			res := util.StringPadLeft(c.Args.Input, c.Args.Size, c.Args.Chr)
			if res != c.Expected {
				t.Errorf("unexpected result %s, expected %s", res, c.Expected)
			}
		})
	}
}
