package gomodel

import (
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"strings"
)

func modelFieldDescs(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.WF("var %sFieldDescs = util.FieldDescs{", m.Proper())
	for _, c := range m.Columns.NotDerived() {
		t := c.Type.String()
		if idx := strings.LastIndex(t, "/"); idx > -1 {
			t = t[idx+1:]
		}
		ret.WF("\t{Key: %q, Title: %q, Description: %q, Type: %q},", c.CamelNoReplace(), c.Title(), c.HelpString, t)
	}
	ret.W("}")
	return ret, nil
}
