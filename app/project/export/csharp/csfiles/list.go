package csfiles

import (
	"fmt"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/csharp"
)

func cshtmlList(ns string, m *model.Model, args *model.Args) (*file.File, error) {
	f := csharp.NewTemplate([]string{ns, "Views", m.Proper()}, m.ProperPlural()+".cshtml")
	b := csharp.NewBlock(m.Proper()+":List", "cshtml")
	b.W("@model IEnumerable<%s>", m.Proper())
	b.W("@inject IconRegistryService IconRegistry")
	b.W("<div class=\"card\">")
	b.W("    <h3>@IconRegistry.Ref(%q, 20, 20, %q)%s</h3>", m.Icon, "icon", m.TitlePlural())
	b.W("    <div class=\"overflow full-width\">")
	b.W("        <table class=\"mt min-200 expanded\">")
	b.W("            <thead>")
	b.W("                <tr>")
	for _, col := range m.Columns {
		b.W("                    <th>%s</th>", col.Title())
	}
	b.W("                </tr>")
	b.W("            </thead>")
	b.W("            <tbody>")
	b.W("                @foreach (var item in Model)")
	b.W("                {")
	b.W("                <tr>")
	for _, col := range m.Columns {
		if col.PK {
			pth := fmt.Sprintf("%s/@item.%s", m.CSRoute(), col.Proper())
			b.W("                    <td><a href=%q>@Html.DisplayFor(_ => item.%s)</a></td>", pth, col.Proper())
		} else {
			r := m.RelationsFor(col)
			if len(r) > 0 {
				tgt := args.Models.Get(r[0].Table)
				b.W("                    <td><a href=\"/%s/@item.%s\">@Html.DisplayFor(_ => item.%s)</a></td>", tgt.CamelLower(), col.Proper(), col.Proper())
			} else {
				b.W("                    <td>@Html.DisplayFor(_ => item.%s)</td>", col.Proper())
			}
		}
	}
	b.W("                </tr>")
	b.W("                }")
	b.W("            </tbody>")
	b.W("        </table>")
	b.W("    </div>")
	b.W("</div>")
	f.AddBlocks(b)
	return f.Render()
}
