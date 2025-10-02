package gohelper

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func BlockFieldDescs(cols model.Columns, str metamodel.StringProvider) (*golang.Block, error) {
	ret := golang.NewBlock(str.Proper(), "struct")
	ret.WF("var %sFieldDescs = util.FieldDescs{", str.Proper())
	for _, c := range cols.NotDerived() {
		t := strings.TrimPrefix(c.Type.String(), "ref:")
		if idx := strings.LastIndex(t, "/"); idx > -1 {
			t = t[idx+1:]
		}
		if c.Example == "" {
			ret.WF("\t{Key: %q, Title: %q, Description: %q, Type: %q},", c.CamelNoReplace(), c.Title(), c.HelpString, t)
		} else {
			ret.WF("\t{Key: %q, Title: %q, Description: %q, Type: %q, Default: %q},", c.CamelNoReplace(), c.Title(), c.HelpString, t, c.Example)
		}
	}
	ret.W("}")
	return ret, nil
}
