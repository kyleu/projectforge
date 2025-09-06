package typescript

import (
	"fmt"
	"slices"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func tsModelConstructor(m *model.Model, enums enum.Enums, ret *golang.Block) {
	s := &util.StringSlice{}
	for _, col := range m.Columns {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), " | undefined", "")
		s.Pushf("%s: %s%s", col.Camel(), tsType(col.Type, enums), optional)
	}
	args := strings.Join(s.Slice, ", ")
	if len(args) < 106 {
		ret.WF("  constructor(%s) {", args)
	} else {
		ret.WF("  constructor(")
		for idx, col := range s.Slice {
			comma := util.Choose(idx+1 < len(s.Slice), ",", "")
			ret.WF("    %s%s", col, comma)
		}
		ret.W("  ) {")
	}
	for _, col := range m.Columns {
		ret.WF("    this.%s = %s;", col.Camel(), col.Camel())
	}
	ret.W("  }")
}

func tsFromObject(m *model.Model, enums enum.Enums, ret *golang.Block) {
	ret.WB()
	ret.WF("  static fromObject(obj: { [_: string]: unknown }): %s {", m.Proper())
	for _, col := range m.Columns {
		switch col.Type.Key() {
		case types.KeyString, types.KeyUUID:
			dflt := util.Choose(col.Nullable || col.HasTag("optional-json"), "undefined", `""`)
			ret.WF(`    const %s = typeof obj.%s === "string" ? obj.%s : %s;`, col.Camel(), col.Camel(), col.Camel(), dflt)
		case types.KeyMap, types.KeyOrderedMap:
			ret.WF(`    const %s = (obj.%s as { [key: string]: unknown }) || {};`, col.Camel(), col.Camel())
		case types.KeyTimestamp, types.KeyTimestampZoned:
			ret.WF(`    const %s = typeof obj.%s === "string" ? new Date(obj.%s) : undefined;`, col.Camel(), col.Camel(), col.Camel())
		default:
			ret.WF("    const %s = obj.%s as %s;", col.Camel(), col.Camel(), tsType(col.Type, enums))
		}
	}

	ret.WB()
	args := m.Columns.CamelNames()
	argsJoined := strings.Join(args, ", ")
	if len(argsJoined) < 100 {
		ret.WF("    return new %s(%s);", m.Proper(), argsJoined)
	} else {
		ret.WF("    return new %s(", m.Proper())
		for idx, c := range m.Columns.CamelNames() {
			comma := util.Choose(idx+1 < len(m.Columns.CamelNames()), ",", "")
			ret.WF("      %s%s", c, comma)
		}
		ret.W("    );")
	}
	ret.W("  }")
}

func tsModel(m *model.Model, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("ts-model-"+m.Name, "ts")
	ret.WF("export class %s {", m.Proper())
	for _, col := range m.Columns {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), " | undefined", "")
		ret.WF("  %s: %s%s;", col.Camel(), tsType(col.Type, enums), optional)
	}
	ret.WB()
	tsModelConstructor(m, enums, ret)
	tsFromObject(m, enums, ret)
	ret.W("}")
	ret.WB()
	ret.WF("export type %s = Array<%s>;", m.ProperPlural(), m.Proper())
	return ret
}

func tsModelContent(imps []string, m *model.Model, enums enum.Enums) golang.Blocks {
	var ret golang.Blocks
	if len(imps) > 0 {
		b := golang.NewBlock("imports", "ts")
		for _, l := range imps {
			b.W(l)
		}
		ret = append(ret, b)
	}
	ret = append(ret, tsModel(m, enums))
	return ret
}

func tsModelImports(allEnums enum.Enums, allModels model.Models, extraTypes model.Models, m *model.Model) ([]string, error) {
	var ret []string
	add := func(s string, args ...any) {
		s = fmt.Sprintf(s, args...)
		if !slices.Contains(ret, s) {
			ret = append(ret, s)
		}
	}
	for _, col := range m.Columns {
		if e, _ := model.AsEnum(col.Type); e != nil {
			if en := allEnums.Get(e.Ref); en != nil {
				if en.PackageWithGroup("") != m.PackageWithGroup("") {
					add(`import type { %s } from "../%s/%s";`, en.Proper(), en.Camel(), en.Camel())
				} else {
					add(`import type { %s } from "./%s";`, en.Proper(), en.Camel())
				}
			}
		}
		r, rm, _ := helper.LoadRef(col, allModels, extraTypes)
		if rm == nil {
			if col.Metadata != nil {
				if tsImport := col.Metadata.GetStringOpt("tsImport"); tsImport != "" {
					add(`import type { %s } from "%s";`, r.K, tsImport)
				}
			}
		} else {
			if tsImport := rm.Config.GetStringOpt("tsImport"); tsImport != "" {
				add(`import type { %s } from "%s";`, r.K, tsImport)
			} else if rm.PackageWithGroup("") != m.PackageWithGroup("") {
				add(`import type { %s } from "../%s/%s";`, r.K, rm.Camel(), rm.Camel())
			} else {
				add(`import type { %s } from "./%s";`, r.K, rm.Camel())
			}
		}
		if col.Type.Key() == types.KeyNumeric {
			relPath := util.StringRepeat("../", len(m.Group)+1)
			add(`import type { Numeric } from "%snumeric/numeric";`, relPath)
		}
	}
	return ret, nil
}

func ModelContent(m *model.Model, allEnums enum.Enums, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, m.PackageWithGroup(""))
	filename := m.CamelLower()
	g := golang.NewGoTemplate(dir, filename+".ts")
	g.AddBlocks(tsModelContent(imps, m, allEnums)...)
	return g.Render(linebreak)
}
