package metaschema

import (
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/util"
)

func ImportArgs(coll *jsonschema.Collection, args *metamodel.Args) (*metamodel.Args, error) {
	if coll == nil || len(coll.Schemas) == 0 {
		return nil, errors.New("empty collection provided for arguments")
	}
	ret := &metamodel.Args{}
	for _, sch := range coll.Schemas {
		switch sch.Type {
		case KeyObject:
			x, err := ImportModel(sch, coll, args)
			if err != nil {
				return nil, err
			}
			ret.Models = append(ret.Models, x)
		case KeyString:
			if len(sch.Enum) == 0 {
				return nil, errors.Errorf("unhandled type [string] without [enum] values for schema [%s]", sch.String())
			}
			x, err := ImportEnum(sch, coll, args)
			if err != nil {
				return nil, err
			}
			ret.Enums = append(ret.Enums, x)
		default:
			return nil, errors.Errorf("invalid type [%v] for schema [%s]", sch.Type, sch.String())
		}
	}
	return ret, nil
}

func ExportAny(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (any, error) {
	switch sch.Type {
	case KeyObject:
		return ImportModel(sch, coll, args)
	case KeyEnum:
		return ImportEnum(sch, coll, args)
	default:
		return nil, errors.Errorf("invalid type [%v] for schema [%s]", sch.Type, sch.String())
	}
}

func parseID(id string) (string, string, []string) {
	parts := util.StringSplitAndTrim(id, "/")
	var pid, pkg string
	var grp []string
	switch len(parts) {
	case 0:
		// noop
	case 1:
		pid = id
	case 2:
		pid = parts[1]
		pkg = parts[0]
	default:
		pid = parts[len(parts)-1]
		pkg = parts[len(parts)-2]
		grp = parts[:len(parts)-2]
	}
	pid = strings.TrimSuffix(strings.TrimSuffix(pid, ".json"), ".schema")
	return pid, pkg, grp
}
