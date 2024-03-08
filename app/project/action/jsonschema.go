package action

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema"

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
	var ret []string
	prj, err := schemaProject(schemata["project"], pm.File)
	if err != nil {
		return nil, err
	}
	ret = append(ret, prj...)

	if pm.EArgs != nil {
		args, err := schemaConfig(schemata["config"], pm.EArgs.ConfigFile)
		if err != nil {
			return nil, err
		}
		ret = append(ret, args...)

		for k, v := range pm.EArgs.EnumFiles {
			enum, err := schemaEnum(schemata["enum"], k, v)
			if err != nil {
				return nil, err
			}
			ret = append(ret, enum...)
		}

		for k, v := range pm.EArgs.ModelFiles {
			model, err := schemaModel(schemata["model"], k, v)
			if err != nil {
				return nil, err
			}
			ret = append(ret, model...)
		}
	}
	return ret, nil
}

func schemaProject(sch *jsonschema.Schema, f json.RawMessage) ([]string, error) {
	var ret []string
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret = append(ret, fmt.Sprintf("[project.json]: %v", err))
	}
	return ret, nil
}

func schemaConfig(sch *jsonschema.Schema, f json.RawMessage) ([]string, error) {
	if len(f) == 0 {
		return nil, nil
	}
	var ret []string
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret = append(ret, fmt.Sprintf("[config.json]: %v", err))
	}
	return ret, nil
}

func schemaEnum(sch *jsonschema.Schema, k string, f json.RawMessage) ([]string, error) {
	var ret []string
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret = append(ret, fmt.Sprintf("[enum/%s.json]: %v", k, err))
	}
	return ret, nil
}

func schemaModel(sch *jsonschema.Schema, k string, f json.RawMessage) ([]string, error) {
	var ret []string
	err := sch.Validate(bytes.NewReader(f))
	if err != nil {
		ret = append(ret, fmt.Sprintf("[model/%s.json]: %v", k, err))
	}
	return ret, nil
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
