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
	b.W("    private const string BaseRoute = %q;", "/"+m.CamelPlural())
	b.W("")
	controllerList(m, b)
	if len(m.PKs()) == 1 {
		b.W("")
		controllerDetail(m, m.PKs()[0], b)
		b.W("")
		controllerDelete(m, m.PKs()[0], b)
	}
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}

func controllerList(m *model.Model, b *csharp.Block) {
	b.W("    [Route(BaseRoute)]")
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> List()")
	b.W("    {")
	b.W("        var ret = await svc.List();")
	b.W("        return Result(%q, ret, %q);", m.TitlePlural(), m.ProperPlural())
	b.W("    }")
}

func controllerDetail(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    [Route(BaseRoute + %q)]", "/{"+pk.Camel()+"}")
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> Get(%s %s)", ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        var ret = await svc.Get(%s);", pk.Camel())
	b.W("        if (ret == null)")
	b.W("        {")
	b.W("            throw new ArgumentNullException(\"no %s available with id [\" + id + \"]\");", m.TitleLower())
	b.W("        }")
	b.W("")
	b.W("        return Result(%q, ret, ret.ToString());", m.Proper())
	b.W("    }")
}

func controllerDelete(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    [Route(BaseRoute + %q)]", "/{"+pk.Camel()+"}/delete")
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> Delete(%s %s)", ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        var ret = await svc.Delete(%s);", pk.Camel())
	b.W("        return Redirect(BaseRoute);")
	b.W("    }")
}
