package view

import (
	"fmt"

	"projectforge.dev/projectforge/app/project/export/model"
)

const endfunc = "{%% endfunc %%}"

func colRow(ind string, col *model.Column, u string, viewString string, link bool) string {
	ret := viewString
	if col.HasTag("title") {
		ret = "<strong>" + ret + "</strong>"
	}
	if (col.PK || col.HasTag("link")) && link {
		ret = fmt.Sprintf("<a href=%q>%s</a>", u, ret)
	}
	return ind + "<td>" + ret + "</td>"
}
