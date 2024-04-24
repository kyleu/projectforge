package csfiles

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func CSAll(p *project.Project, addHeader bool) (file.Files, error) {
	if p.ExportArgs == nil {
		return nil, errors.New("export arguments aren't loaded")
	}
	args := p.ExportArgs
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	ret := make(file.Files, 0, (len(args.Models)*10)+len(args.Enums))

	for _, m := range args.Models {
		calls, err := CSModelAll(m, p, args, addHeader)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, calls...)
	}

	return ret, nil
}

func CSModelAll(m *model.Model, p *project.Project, args *model.Args, addHeader bool) (file.Files, error) {
	var ret file.Files

	list, err := cshtmlList(m, args, addHeader)
	if err != nil {
		return nil, err
	}
	ret = append(ret, list)

	return ret, nil
}

func cshtmlList(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := csharp.NewTemplate([]string{"Views", m.Title()}, m.TitlePlural()+".cshtml")
	b := csharp.NewBlock(m.Title()+":List", "cshtml")
	b.W("@model IEnumerable<%s>", m.Title())
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
		b.W("                    <td>@Html.DisplayFor(modelItem => item.%s)</td>", col.Proper())
	}
	b.W("                </tr>")
	b.W("                }")
	b.W("            </tbody>")
	b.W("        </table>")
	b.W("    </div>")
	b.W("</div>")

	/*





	               <tr>
	                   <th>ID</th>
	                   <th>Email</th>
	                   <th>Full Name</th>
	                   <th>Phone</th>
	                   <th>DoB</th>
	                   <th>Last Login</th>
	                   <th>Status</th>
	                   <th>Created</th>
	               </tr>

	               @foreach (var item in Model) {
	               <tr>
	                   <td>@Html.DisplayFor(modelItem => item.ID)</td>
	                   <td>@Html.DisplayFor(modelItem => item.Email)</td>
	                   <td>@Html.DisplayFor(modelItem => item.FullName)</td>
	                   <td>@Html.DisplayFor(modelItem => item.Phone)</td>
	                   <td>@Html.DisplayFor(modelItem => item.DateOfBirth)</td>
	                   <td>@Html.DisplayFor(modelItem => item.LastLoggedIn)</td>
	                   <td>@Html.DisplayFor(modelItem => item.Status)</td>
	                   <td>@Html.DisplayFor(modelItem => item.DateCreated)</td>
	               </tr>
	               }
	           </tbody>

	   </div>

	*/
	g.AddBlocks(b)
	return g.Render(addHeader)
}
