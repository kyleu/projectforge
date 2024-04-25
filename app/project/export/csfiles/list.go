package csfiles

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func cshtmlList(m *model.Model) (*file.File, error) {
	f := csharp.NewTemplate([]string{"Views", m.Title()}, m.TitlePlural()+".cshtml")
	b := csharp.NewBlock(m.Title()+":List", "cshtml")
	b.W("@model IEnumerable<%s>", m.Namespace())
	b.W("<div class=\"card\">")
	b.W("    <h3>%s</h3>", m.TitlePlural())
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
	b.W("                @foreach (var item in Model) {")
	b.W("                <tr>")
	for _, col := range m.Columns {
		b.W("                    <td>@Html.DisplayFor(_ => item.%s)</td>", col.Proper())
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
