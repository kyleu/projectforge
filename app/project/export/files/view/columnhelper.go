package view

import (
	"fmt"

	"projectforge.dev/projectforge/app/project/export/model"
)

func colRow(ind string, col *model.Column, u string, viewString string, link bool) string {
	switch {
	case (col.PK || col.HasTag("link")) && link:
		return fmt.Sprintf(ind+"<td><a href=%q>%s</a></td>", u, viewString)
	case col.HasTag("title"):
		return ind + "<td><strong>" + viewString + "</strong></td>"
	default:
		return ind + "<td>" + viewString + "</td>"
	}
}
