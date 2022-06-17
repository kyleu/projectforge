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
