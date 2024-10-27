package gomodel

import (
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func modelFieldDescs(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.WF("var %sFieldDescs = util.FieldDescs{", m.Proper())
	for _, c := range m.Columns.NotDerived() {
		ret.WF("\t{Key: %q, Title: %q, Description: %q, Type: %q},", c.CamelNoReplace(), c.Title(), c.HelpString, c.Type.String())
	}
	ret.W("}")
	return ret, nil
}
