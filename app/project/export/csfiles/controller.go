package csfiles

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func controller(m *model.Model, p *project.Project) (*file.File, error) {
	f := csharp.NewFile(p.Package+".Controllers", []string{"Controllers"}, m.Title()+"Controller.cs")
	f.AddImport(p.Package+".Entities", ImpMVC, ImpEF, ImpControllers)
	b := csharp.NewBlock("Controller", "class")
	b.W("public class %sController(Database db) : BaseController()", m.Title())
	b.W("{")
	b.W("    [Route(%q)]", "/"+m.CamelPlural())
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> %s()", m.TitlePlural())
	b.W("    {")
	b.W("        var ret = await db.%s.AsQueryable().ToListAsync();", m.TitlePlural())
	b.W("        ViewData[%q] = %q;", "Title", m.TitlePlural())
	b.W("        return Result(ret);")
	b.W("    }")
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}
