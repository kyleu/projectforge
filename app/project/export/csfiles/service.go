package csfiles

import (
	"fmt"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func service(m *model.Model, p *project.Project) (*file.File, error) {
	f := csharp.NewFile(p.Package+".Services."+m.Proper(), []string{"Services", m.Proper()}, m.Title()+"Service.cs")
	f.AddImport(p.Package+".Entities", ImpEF, ImpServices)
	b := csharp.NewBlock("Service", "class")
	b.W("public class %sService(Database db) : BaseService<Entities.%s>(db, db.%s)", m.Title(), m.Proper(), m.ProperPlural())
	b.W("{")
	serviceList(m, b)
	if len(m.PKs()) == 1 {
		b.W("")
		serviceDetail(m, m.PKs()[0], b)
		b.W("")
		serviceDelete(m, m.PKs()[0], b)
	}
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}

func serviceList(m *model.Model, b *csharp.Block) {
	b.W("    public async Task<List<Entities.%s>> List()", m.Proper())
	b.W("    {")
	b.W("        return await Models.Take(100).AsQueryable().ToListAsync();")
	b.W("    }")
}

func serviceDetail(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    public async Task<Entities.%s?> Get(%s %s)", m.Proper(), ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        return await Models.FirstOrDefaultAsync(x => x.%s == %s);", pk.Proper(), pk.Camel())
	b.W("    }")
}

func serviceDelete(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    public async Task<Entities.%s?> Delete(%s %s)", m.Proper(), ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        var ret = await Models.FirstOrDefaultAsync(x => x.%s == %s);", pk.Proper(), pk.Camel())
	b.W("        if (ret != null)")
	b.W("        {")
	b.W("            Models.Remove(ret);")
	b.W("        }")
	b.W("        return ret;")
	b.W("    }")
}

func ToCSharpType(col *model.Column) any {
	switch col.Type.Key() {
	case "uuid":
		return "Guid"
	default:
		return fmt.Sprintf("unknown-type[%s]", col.Type.String())
	}
}
