// $PF_GENERATE_ONCE$
// Change this file to include your app's settings
package settings

import (
	"fmt"

	"{{{ .Package }}}/app/util"
)

type Settings struct {
	ExampleBool   bool   `json:"exampleBool,omitzero"`
	ExampleString string `json:"exampleString,omitzero"`
}

func NewSettings(exampleBool bool, exampleString string) *Settings {
	return &Settings{ExampleBool: exampleBool, ExampleString: exampleString}
}

func (s *Settings) Clone() *Settings {
	return &Settings{s.ExampleBool, s.ExampleString}
}

func (s *Settings) ToMap() util.ValueMap {
	return util.ValueMap{"exampleBool": s.ExampleBool, "exampleString": s.ExampleString}
}

func SettingsFromMap(m util.ValueMap, setPK bool) (*Settings, util.ValueMap, error) {
	ret := &Settings{}
	extra := util.ValueMap{}
	for k, v := range m {
		var err error
		switch k {
		case "exampleBool":
			ret.ExampleBool, err = m.ParseBool(k, true, true)
		case "exampleString":
			ret.ExampleString, err = m.ParseString(k, true, true)
		default:
			extra[k] = v
		}
		if err != nil {
			return nil, nil, err
		}
	}
	// $PF_SECTION_START(extrachecks)$
	// $PF_SECTION_END(extrachecks)$
	return ret, extra, nil
}

func (s *Settings) ToOrderedMap() *util.OrderedMap[any] {
	pairs := util.OrderedPairs[any]{{K: "exampleBool", V: s.ExampleBool}, {K: "exampleString", V: s.ExampleString}}
	return util.NewOrderedMap(false, 4, pairs...)
}

func (s *Settings) Diff(sx *Settings) util.Diffs {
	var diffs util.Diffs
	if s.ExampleBool != sx.ExampleBool {
		diffs = append(diffs, util.NewDiff("exampleBool", fmt.Sprint(s.ExampleBool), fmt.Sprint(sx.ExampleBool)))
	}
	if s.ExampleString != sx.ExampleString {
		diffs = append(diffs, util.NewDiff("exampleString", s.ExampleString, sx.ExampleString))
	}
	return diffs
}

func (s *Settings) Strings() []string {
	return []string{fmt.Sprint(s.ExampleBool), s.ExampleString}
}

func (s *Settings) ToCSV() ([]string, [][]string) {
	return SettingsFieldDescs.Keys(), [][]string{s.Strings()}
}

func (s *Settings) ToData() []any {
	return []any{s.ExampleBool, s.ExampleString}
}

func RandomSettings() *Settings {
	return &Settings{
		ExampleBool:   util.RandomBool(),
		ExampleString: util.RandomID(),
	}
}

var SettingsFieldDescs = util.FieldDescs{
	{Key: "exampleBool", Title: "Example Bool", Description: "", Type: "bool"},
	{Key: "exampleString", Title: "Example String", Description: "", Type: "string"},
}
