// Content managed by Project Forge, see [projectforge.md] for details.
//go:build test_all || !func_test
// +build test_all !func_test

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

func TestPadRight(t *testing.T) {
	t.Parallel()

	type Args struct {
		Input string
		Size  int
		Chr   rune
	}

	cases := []struct {
		Name     string
		Args     Args
		Expected string
	}{
		{
			Name: "normal args",
			Args: Args{
				Input: "foo",
				Size:  5,
				Chr:   '界',
			},
			Expected: "foo界界",
		},
		{
			Name: "nothing to pad, size is equal to input len",
			Args: Args{
				Input: "привет",
				Size:  6,
				Chr:   '界',
			},
			Expected: "привет",
		},
		{
			Name: "nothing to pad, size is smaller than input len",
			Args: Args{
				Input: "привет",
				Size:  4,
				Chr:   '界',
			},
			Expected: "привет",
		},
		{
			Name: "negative size, input is left intact",
			Args: Args{
				Input: "привет",
				Size:  -10,
				Chr:   '界',
			},
			Expected: "привет",
		},
		{
			Name: "empty input",
			Args: Args{
				Input: "",
				Size:  3,
				Chr:   '界',
			},
			Expected: "界界界",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			res := util.StringPadRight(tCase.Args.Input, tCase.Args.Size, tCase.Args.Chr)
			if res != tCase.Expected {
				t.Errorf("unexpected result %s, expected %s", res, tCase.Expected)
			}
		})
	}
}

func BenchmarkStringPadRight(b *testing.B) {
	padSize := 1_000_000
	inputString := "test"

	b.SetBytes(int64(len(inputString) + padSize))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var res string
		res = util.StringPadRight(inputString, padSize, ' ')
		_ = res
	}

	/*
		goos: darwin
		goarch: arm64
		pkg: projectforge.dev/projectforge/app/util
		BenchmarkStringPadRight
									   run times(b.N) avg time/op       throughput       (?)                memory allocations/op
		# StringPadRight with strings.Builder
		BenchmarkStringPadRight-8   	     445	   2607051 ns/op	 383.58 MB/s	 5241614 B/op	      33 allocs/op
		-------------------------------------------------------------------------------------------------------------------------
									   run times(b.N) avg time/op       throughput       (?)                memory allocations/op
		# StringPadRight without strings.Builder
		BenchmarkStringPadRight-8   	       1	37883665875 ns/op	   0.03 MB/s	503999148496 B/op	 2015246 allocs/op
	*/
}

func TestPadLeft(t *testing.T) {
	t.Parallel()

	type Args struct {
		Input string
		Size  int
		Chr   rune
	}

	cases := []struct {
		Name     string
		Args     Args
		Expected string
	}{
		{
			Name: "normal args",
			Args: Args{
				Input: "foo",
				Size:  5,
				Chr:   '界',
			},
			Expected: "界界foo",
		},
		{
			Name: "nothing to pad, size is equal to input len",
			Args: Args{
				Input: "привет",
				Size:  6,
				Chr:   '界',
			},
			Expected: "привет",
		},
		{
			Name: "nothing to pad, size is smaller than input len",
			Args: Args{
				Input: "привет",
				Size:  4,
				Chr:   '界',
			},
			Expected: "привет",
		},
		{
			Name: "negative size, input is left intact",
			Args: Args{
				Input: "привет",
				Size:  -10,
				Chr:   '界',
			},
			Expected: "привет",
		},
		{
			Name: "empty input",
			Args: Args{
				Input: "",
				Size:  3,
				Chr:   '界',
			},
			Expected: "界界界",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			res := util.StringPadLeft(tCase.Args.Input, tCase.Args.Size, tCase.Args.Chr)
			if res != tCase.Expected {
				t.Errorf("unexpected result %s, expected %s", res, tCase.Expected)
			}
		})
	}
}
