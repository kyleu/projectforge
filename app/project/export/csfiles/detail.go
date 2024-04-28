package csfiles

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func cshtmlDetail(m *model.Model) (*file.File, error) {
	f := csharp.NewTemplate([]string{"Views", m.Proper()}, m.Proper()+".cshtml")
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
		b.W("                    <td>@Html.DisplayFor(_ => Model.%s)</td>", col.Proper())
		b.W("                </tr>")
	}
	b.W("            </tbody>")
	b.W("        </table>")
	b.W("    </div>")
	b.W("</div>")
	f.AddBlocks(b)
	return f.Render()
}
