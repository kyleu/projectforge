package gomodel

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/util"
)

func modelRandom(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Random", "struct")
	ret.W("func Random() *%s {", m.Proper())
	ret.W("\treturn &%s{", m.Proper())
	maxColLength := m.Columns.MaxCamelLength() + 1
	for _, col := range m.Columns {
		ret.W("\t\t%s %s,", util.StringPad(col.Proper()+":", maxColLength), randFor(col))
	}
	ret.W("\t}")
	ret.W("}")
	return ret
}

func randFor(col *model.Column) string {
	switch col.Type.Key {
	case model.TypeBool.Key:
		return "util.RandomBool()"
	case model.TypeInt.Key:
		return "util.RandomInt(10000)"
	case model.TypeInterface.Key:
		return "nil"
	case model.TypeMap.Key:
		return "util.RandomValueMap(4)"
	case model.TypeString.Key:
		return "util.RandomString(12)"
	case model.TypeTimestamp.Key:
		if col.HasTag("deleted") {
			return "nil"
		}
		if col.Nullable {
			return "util.NowPointer()"
		}
		return "time.Now()"
	case model.TypeUUID.Key:
		if col.Nullable {
			return "util.UUIDP()"
		}
		return "util.UUID()"
	default:
		return "TODO"
	}
}
