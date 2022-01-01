package files

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func controllerArgFor(col *model.Column, b *golang.Block) {
	switch col.Type.Key {
	case model.TypeInt.Key:
		b.W("\t%sArgStr, err := rcRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn nil, errors.Wrap(err, \"must provide [%s] as an argument\")", col.Camel())
		b.W("\t}")
		b.W("\t%sArg, err := strconv.Atoi(%sArgStr)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn nil, errors.Wrap(err, \"must provide [%s] as an argument\")", col.Camel())
		b.W("\t}")
	case model.TypeString.Key:
		b.W("\t%sArg, err := rcRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn nil, errors.Wrap(err, \"must provide [%s] as an argument\")", col.Camel())
		b.W("\t}")
	case model.TypeUUID.Key:
		b.W("\t%sArgStr, err := rcRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn nil, errors.Wrap(err, \"must provide [%s] as an argument\")", col.Camel())
		b.W("\t}")
		b.W("\t%sArgP := util.UUIDFromString(%sArgStr)", col.Camel(), col.Camel())
		b.W("\tif %sArgP == nil {", col.Camel())
		b.W("\t\treturn nil, errors.Errorf(\"argument [%s] (%%%%s) is not a valid UUID\", %sArgStr)", col.Camel(), col.Camel())
		b.W("\t}")
		b.W("\t%sArg := *%sArgP", col.Camel(), col.Camel())
	default:
		b.W("\tERROR: unhandled controller arg type [%s]", col.Type.String())
	}
}

func grpcArgFor(col *model.Column, b *golang.Block, zeroVals string) {
	switch col.Type.Key {
	case model.TypeInt.Key:
		b.W("\t%s, err := provider.GetRequestInt(r, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
	case model.TypeString.Key:
		b.W("\t%s, err := provider.GetRequestString(r, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
	default:
		b.W("\tERROR: unhandled arg type [%s]", col.Type.String())
	}
}
