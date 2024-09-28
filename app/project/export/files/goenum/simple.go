package goenum

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const timePointer = "*time.Time"

func structSimple(e *enum.Enum) []*golang.Block {
	tBlock := golang.NewBlock(e.Proper(), "typealias")
	tBlock.W("type (")
	tBlock.WF("\t%s string", util.StringPad(e.Proper(), len(e.ProperPlural())))
	tBlock.WF("\t%s []%s", e.ProperPlural(), e.Proper())
	tBlock.W(")")

	cBlock := golang.NewBlock(e.Proper(), "constvar")
	cBlock.W("const (")
	maxCount := util.StringArrayMaxLength(e.ValuesCamel())
	pl := len(e.Proper())
	maxColLength := maxCount + pl
	lo.ForEach(e.Values, func(v *enum.Value, _ int) {
		cBlock.WF("\t%s %s = %q", util.StringPad(e.Proper()+v.Proper(), maxColLength), e.Proper(), v.Key)
	})
	cBlock.W(")")
	return []*golang.Block{tBlock, cBlock}
}
