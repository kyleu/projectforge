package csfiles

import (
	"fmt"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/csharp"
)

func controller(ns string, m *model.Model) (*file.File, error) {
	f := csharp.NewFile(ns+".Controllers", []string{ns, "Controllers"}, m.Proper()+"Controller.cs")
	if ns != "Shared" {
		f.AddImport(ImpControllers)
	}
	f.AddImport(ns+".Entities", ns+".Services."+m.Proper(), ImpMVC)
	b := csharp.NewBlock("Controller", "class")
	b.W("public class %sController(%sService svc) : BaseController", m.Proper(), m.Proper())
	b.W("{")
	b.W("    private const string BaseRoute = %q;", "/"+m.CamelLower())
	b.W("")
	controllerList(m, b)
	b.W("")
	controllerCreate(m, b)
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
	b.W("        return Ok(ret);")
	b.W("    }")
}

func controllerCreate(m *model.Model, b *csharp.Block) {
	b.W("    [Route(BaseRoute)]")
	b.W("    [HttpPost]")
	b.W("    public async Task<IActionResult> Create([FromBody] %s mdl)", m.Proper())
	b.W("    {")
	b.W("        svc.Models.Add(mdl);")
	b.W("        await svc.Flush();")
	b.W("        return Ok(mdl);")
	b.W("    }")
}

func controllerDetail(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    [Route(BaseRoute + %q)]", fmt.Sprintf("/{%s%s}", pk.Camel(), ToCSharpViewType(pk)))
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> Get(%s %s)", ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        var ret = await svc.Get(%s);", pk.Camel())
	b.W("        if (ret == null)")
	b.W("        {")
	b.W("            throw new ArgumentNullException(\"no %s available with id [\" + id + \"]\");", m.TitleLower())
	b.W("        }")
	b.W("")
	b.W("        return Ok(ret);")
	b.W("    }")
}

func controllerDelete(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    [Route(BaseRoute + %q)]", fmt.Sprintf("/{%s%s}/delete", pk.Camel(), ToCSharpViewType(pk)))
	b.W("    [HttpGet]")
	b.W("    public async Task<IActionResult> Delete(%s %s)", ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        await svc.Delete(%s);", pk.Camel())
	b.W("        return Redirect(BaseRoute);")
	b.W("    }")
}
