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
		if c.Type.Key() == "numericMap" {
			t = "map"
		}

		msg := "\t{Key: %q, Title: %q"
		args := []any{c.CamelNoReplace(), c.Title()}
		if c.Help != "" {
			msg += ", Description: %q"
			args = append(args, c.Help)
		}
		msg += ", Type: %q"
		args = append(args, t)

		if c.Example != "" {
			msg += ", Default: %q"
			args = append(args, c.Example)
		}
		msg += "},"

		ret.WF(msg, args...)
	}
	ret.W("}")
	return ret, nil
}
