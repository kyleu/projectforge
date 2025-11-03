package typescript

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func tsFromObjectColumn(col *model.Column, enums enum.Enums, ret *golang.Block) error {
	opt := col.Nullable || col.HasTag("optional-json")
	switch col.Type.Key() {
	case types.KeyTimestamp, types.KeyTimestampZoned:
		if opt {
			ret.WF(`    const %s = Parse.dateOpt(obj.%s);`, col.Camel(), col.Camel())
		} else {
			ret.WF(`    const %s = Parse.date(obj.%s, () => new Date());`, col.Camel(), col.Camel())
		}
	case types.KeyInt:
		if opt {
			ret.WF(`    const %s = Parse.intOpt(obj.%s);`, col.Camel(), col.Camel())
		} else {
			ret.WF(`    const %s = Parse.int(obj.%s, () => 0);`, col.Camel(), col.Camel())
		}
	case types.KeyFloat:
		if opt {
			ret.WF(`    const %s = Parse.floatOpt(obj.%s);`, col.Camel(), col.Camel())
		} else {
			ret.WF(`    const %s = Parse.float(obj.%s, () => 0);`, col.Camel(), col.Camel())
		}
	case types.KeyString, types.KeyUUID:
		if opt {
			ret.WF(`    const %s = Parse.stringOpt(obj.%s);`, col.Camel(), col.Camel())
		} else {
			ret.WF(`    const %s = Parse.string(obj.%s, () => "");`, col.Camel(), col.Camel())
		}
	case types.KeyMap, types.KeyOrderedMap, types.KeyAny:
		if opt {
			ret.WF(`    const %s = Parse.objOpt(obj.%s);`, col.Camel(), col.Camel())
		} else {
			ret.WF(`    const %s = Parse.obj(obj.%s, () => ({}));`, col.Camel(), col.Camel())
		}
	case types.KeyNumeric:
		ret.WF(`    const %s = new Numeric(obj.%s as NumericSource);`, col.Camel(), col.Camel())
	case types.KeyNumericMap:
		if opt {
			ret.WF(`    const %s = NumericMap.parse(Parse.objOpt(obj.%s) ?? {});`, col.Camel(), col.Camel())
		} else {
			ret.WF(`    const %s = NumericMap.parse(Parse.obj(obj.%s, () => ({})));`, col.Camel(), col.Camel())
		}
	case types.KeyEnum:
		op := util.Choose(opt, "parse", "get")
		e, err := model.AsEnum(col.Type)
		if err != nil {
			return err
		}
		en := enums.Get(e.Ref)
		suffix := ""
		if dflt := en.Values.Default(); dflt != nil {
			suffix = fmt.Sprintf(", () => %q", dflt.Key)
		}
		ret.WF(`    const %s = %s%s(Parse.string(obj.%s%s));`, col.Camel(), op, en.Proper(), col.Camel(), suffix)
	default:
		ret.WF("    const %s = obj.%s as %s;", col.Camel(), col.Camel(), tsType(col.Type, enums))
	}
	return nil
}
