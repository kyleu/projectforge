package action

import (
	"bytes"
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema"

	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/assets"
)

var schemata = map[string]*jsonschema.Schema{}

func schemaCheck(pm *PrjAndMods) ([]string, error) {
	if len(schemata) == 0 {
		err := loadSchemata()
		if err != nil {
			return nil, err
		}
	}
	ret := &util.StringSlice{}
	prj, err := schemaProject(schemata["project"], pm.File)
	if err != nil {
		return nil, err
	}
	ret.Push(prj...)

	if pm.Prj.ExportArgs != nil {
		args, err := schemaConfig(schemata["config"], pm.Prj.ExportArgs.ConfigFile)
		if err != nil {
			return nil, err
		}
		ret.Push(args...)

		for k, v := range pm.Prj.ExportArgs.EnumFiles {
			enum, err := schemaEnum(schemata["enum"], k, v)
			if err != nil {
				return nil, err
			}
			ret.Push(enum...)
		}

		for k, v := range pm.Prj.ExportArgs.ModelFiles {
			model, err := schemaModel(schemata["model"], k, v)
			if err != nil {
				return nil, err
			}
			ret.Push(model...)
		}
	}
	return ret.Slice, nil
}

func schemaProject(sch *jsonschema.Schema, f json.RawMessage) ([]string, error) {
	ret := &util.StringSlice{}
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret.Pushf("[project.json]: %v", err)
	}
	return ret.Slice, nil
}

func schemaConfig(sch *jsonschema.Schema, f json.RawMessage) ([]string, error) {
	if len(f) == 0 {
		return nil, nil
	}
	ret := &util.StringSlice{}
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret.Pushf("[config.json]: %v", err)
	}
	return ret.Slice, nil
}

func schemaEnum(sch *jsonschema.Schema, k string, f json.RawMessage) ([]string, error) {
	ret := &util.StringSlice{}
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret.Pushf("[enum/%s.json]: %v", k, err)
	}
	return ret.Slice, nil
}

func schemaModel(sch *jsonschema.Schema, k string, f json.RawMessage) ([]string, error) {
	ret := &util.StringSlice{}
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret.Pushf("[model/%s.json]: %v", k, err)
	}
	return ret.Slice, nil
}

func loadSchemata() error {
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft7

	x := func(k string) error {
		e, err := assets.Embed("schema/" + k + ".schema.json")
		if err != nil {
			return err
		}
		err = c.AddResource(k, bytes.NewReader(e.Bytes))
		if err != nil {
			return err
		}
		sch, err := c.Compile(k)
		if err != nil {
			return err
		}
		schemata[k] = sch
		return nil
	}
	err := x("project")
	if err != nil {
		return err
	}
	err = x("config")
	if err != nil {
		return err
	}
	err = x("enum")
	if err != nil {
		return err
	}
	err = x("model")
	if err != nil {
		return err
	}
	return nil
}
