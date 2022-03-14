// Content managed by Project Forge, see [projectforge.md] for details.
package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

var (
	testArr = []interface{}{"x", 1, true}
	testObj = util.ValueMap{"x": util.ValueMap{"y": 1}}
)

var diffTests = []struct {
	k string
	l interface{}
	r interface{}
	d util.Diffs
}{
	{k: "bool.same", l: true, r: true, d: util.Diffs{}},
	{k: "bool.different", l: true, r: false, d: util.Diffs{
		&util.Diff{Old: "true", New: "false"}},
	},
	{k: "int.same", l: 1, r: 1, d: util.Diffs{}},
	{k: "int.different", l: 1, r: 10, d: util.Diffs{
		&util.Diff{Old: "1", New: "10"}},
	},
	{k: "string.same", l: "x", r: "x", d: util.Diffs{}},
	{k: "string.different", l: "x", r: "y", d: util.Diffs{
		&util.Diff{Old: "x", New: "y"}},
	},
	{k: "map.same", l: testObj, r: testObj, d: util.Diffs{}},
	{k: "map.different", l: testObj, r: util.ValueMap{"x": util.ValueMap{"y": 2}}, d: util.Diffs{
		&util.Diff{Path: ".x.y", Old: "1", New: "2"}},
	},
	{k: "map.missing", l: testObj, r: util.ValueMap{"y": true}, d: util.Diffs{
		&util.Diff{Path: ".x", Old: "map[y:1]", New: ""},
		&util.Diff{Path: ".y", Old: "", New: "true"}},
	},
	{k: "map.extra", l: testObj, r: util.ValueMap{"x": util.ValueMap{"y": 1, "z": true}}, d: util.Diffs{
		&util.Diff{Path: ".x.z", Old: "", New: "true"}},
	},
	{k: "array.same", l: testArr, r: testArr, d: util.Diffs{}},
	{k: "array.different", l: testArr, r: []interface{}{"y", 2, true}, d: util.Diffs{
		&util.Diff{Path: ".0", Old: "x", New: "y"},
		&util.Diff{Path: ".1", Old: "1", New: "2"}},
	},
	{k: "array.missing", l: testArr, r: []interface{}{"x", 1}, d: util.Diffs{
		&util.Diff{Path: ".2", Old: "true", New: ""}},
	},
	{k: "array.extra", l: testArr, r: []interface{}{"x", 1, true, "new"}, d: util.Diffs{
		&util.Diff{Path: ".3", Old: "", New: "new"}},
	},
}

func TestDiffs(t *testing.T) {
	t.Parallel()

	for _, tt := range diffTests {
		diffs := util.DiffObjects(tt.l, tt.r, "")
		diffJSON := util.ToJSONCompact(diffs)
		if len(diffs) == len(tt.d) {
			for i := 0; i < len(diffs); i++ {
				expected := tt.d[i]
				observed := diffs[i]
				if expected.Path != observed.Path {
					t.Errorf("%s: diff [%d] has path [%s], expected path [%s]: %s", tt.k, i, observed.Path, expected.Path, diffJSON)
				}
				if expected.Old != observed.Old {
					t.Errorf("%s: diff [%d] has old value [%s], expected old value [%s]: %s", tt.k, i, observed.Old, expected.Old, diffJSON)
				}
				if expected.New != observed.New {
					t.Errorf("%s: diff [%d] has new value [%s], expected new value [%s]: %s", tt.k, i, observed.New, expected.New, diffJSON)
				}
			}
		} else {
			t.Errorf("%s: found [%d] diffs, expected [%d]: %s", tt.k, len(diffs), len(tt.d), diffJSON)
		}
	}
}
