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
	f.AddImport(p.Package+".Entities", ImpEF)
	b := csharp.NewBlock("Service", "class")
	b.W("public class %sService(Database db) : BaseService(db)", m.Title())
	b.W("{")
	b.W("    protected internal readonly DbSet<Entities.%s> _%s = db.%s;", m.Proper(), m.CamelPlural(), m.ProperPlural())
	b.W("")
	serviceList(m, b)
	if len(m.PKs()) == 1 {
		b.W("")
		serviceDetail(m, m.PKs()[0], b)
	}
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}

func serviceList(m *model.Model, b *csharp.Block) {
	b.W("    public async Task<List<Entities.%s>> %s()", m.Proper(), m.ProperPlural())
	b.W("    {")
	b.W("        return await _%s.Take(100).AsQueryable().ToListAsync();", m.CamelPlural())
	b.W("    }")
}

func serviceDetail(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    public async Task<Entities.%s?> %sBy%s(%s %s)", m.Proper(), m.Proper(), pk.Proper(), ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        return await _%s.FirstAsync(x => x.%s == %s);", m.CamelPlural(), pk.Proper(), pk.Camel())
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
