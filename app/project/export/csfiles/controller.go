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
	controllerList(m, b)
	if len(m.PKs()) == 1 {
		b.W("")
		controllerDetail(m, m.PKs()[0], b)
	}
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}

func controllerList(m *model.Model, b *csharp.Block) {
	b.W("    [Route(%q)]", "/"+m.CamelPlural())
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> %s()", m.TitlePlural())
	b.W("    {")
	b.W("        var ret = await svc.%s();", m.TitlePlural())
	b.W("        ViewData[%q] = %q;", "Title", m.TitlePlural())
	b.W("        return Result(ret);")
	b.W("    }")
}

func controllerDetail(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    [Route(%q)]", "/"+m.CamelPlural()+"/{"+pk.Camel()+"}")
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> %s(%s %s)", m.Proper(), ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        var ret = await svc.%sBy%s(%s);", m.Proper(), pk.Proper(), pk.Camel())
	b.W("        if (ret == null)")
	b.W("        {")
	b.W("            throw new ArgumentNullException(\"no %s available with id [\" + id + \"]\");", m.TitleLower())
	b.W("        }")
	b.W("        ViewData[\"Title\"] = ret.ToString();")
	b.W("        return Result(ret);")
	b.W("    }")
}
