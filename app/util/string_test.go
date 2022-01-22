// Content managed by Project Forge, see [projectforge.md] for details.
package util_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/kyleu/projectforge/app/util"
)

type ToTitleTest struct {
	TestValue string
	Expected  string
}

func (t *ToTitleTest) Test() error {
	if res := util.StringToTitle(t.TestValue); res != t.Expected {
		return errors.Errorf("StringToTitle returned [%s], not expected [%s]", res, t.Expected)
	}
	return nil
}

var titleTests = []*ToTitleTest{
	{TestValue: "SimpleCamelCase", Expected: "Simple Camel Case"},
	{TestValue: "CSVFilesAreCoolButTXTRules", Expected: "CSV Files Are Cool But TXT Rules"},
	{TestValue: "MediaTypes", Expected: "Media Types"},
}

func TestToTitle(t *testing.T) {
	t.Parallel()
	for _, test := range titleTests {
		err := test.Test()
		if err != nil {
			t.Error(err)
		}
	}
}
