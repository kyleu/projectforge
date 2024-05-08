package csfiles

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/csharp"
)

func cshtmlDetail(ns string, m *model.Model, args *model.Args) (*file.File, error) {
	f := csharp.NewTemplate([]string{ns, "Views", m.Proper()}, m.Proper()+".cshtml")
	b := csharp.NewBlock(m.Proper()+":Detail", "cshtml")
	b.W("@model %s", m.Proper())
	b.W("@inject IconRegistryService IconRegistry")
	b.W("<div class=\"card\">")
	b.W("    <h3>@IconRegistry.Ref(%q, 20, 20, %q)@Model.ToString()</h3>", m.Icon, "icon")
	b.W("    <em>%s</em>", m.Title())
	b.W("    <div class=\"overflow full-width\">")
	b.W("        <table class=\"mt min-200 expanded\">")
	b.W("            <tbody>")
	for _, col := range m.Columns {
		b.W("                <tr>")
		b.W("                    <th class=\"shrink\">%s</th>", col.Title())
		r := m.RelationsFor(col)
		if len(r) > 0 {
			msg := "                    <td><a href=\"%s/@Model.%s\">@Html.DisplayFor(_ => Model.%s)</a> (@Model.%s?.ToString())</td>"
			tgt := args.Models.Get(r[0].Table)
			b.W(msg, tgt.CSRoute(), col.Proper(), col.Proper(), tgt.Proper())
		} else {
			b.W("                    <td>@Html.DisplayFor(_ => Model.%s)</td>", col.Proper())
		}
		b.W("                </tr>")
	}
	b.W("            </tbody>")
	b.W("        </table>")
	b.W("    </div>")
	b.W("</div>")
	f.AddBlocks(b)
	return f.Render()
}
