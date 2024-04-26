package csfiles

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func controller(m *model.Model, p *project.Project) (*file.File, error) {
	f := csharp.NewFile(p.Package+".Controllers", []string{"Controllers"}, m.Title()+"Controller.cs")
	f.AddImport(p.Package+".Services."+m.Proper(), ImpMVC, ImpControllers)
	b := csharp.NewBlock("Controller", "class")
	b.W("public class %sController(%sService svc) : BaseController", m.Title(), m.Proper())
	b.W("{")
	b.W("    [Route(%q)]", "/"+m.CamelPlural())
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> %s()", m.TitlePlural())
	b.W("    {")
	b.W("        var ret = await svc.%s();", m.TitlePlural())
	b.W("        ViewData[%q] = %q;", "Title", m.TitlePlural())
	b.W("        return Result(ret);")
	b.W("    }")
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}
