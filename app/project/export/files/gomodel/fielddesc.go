package gomodel

import (
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func modelFieldDescs(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("var FieldDescs = util.FieldDescs{")
	for _, c := range m.Columns {
		ret.W("\t{Key: %q, Title: %q, Description: %q, Type: %q},", c.Camel(), c.Title(), c.HelpString, c.Type.String())
	}
	ret.W("}")
	return ret, nil
}
