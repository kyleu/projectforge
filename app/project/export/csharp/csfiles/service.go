package csfiles

import (
	"fmt"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func service(ns string, m *model.Model, args *model.Args) (*file.File, error) {
	f := csharp.NewFile(ns+".Services."+m.Proper(), []string{ns, "Services", m.Proper()}, m.Proper()+"Service.cs")
	cls := "Database"
	if ns == "Shared" {
		f.AddImport(ns+".Database", ImpEF)
		cls = "DatabaseContext"
	} else {
		f.AddImport(ns+".Entities", ImpEF, ImpServices)
	}
	b := csharp.NewBlock("Service", "class")
	b.W("public class %sService(%s db) : BaseService<Entities.%s>(db, db.%s)", m.Proper(), cls, m.Proper(), m.ProperPlural())
	b.W("{")
	if ns != "Shared" {
		b.W("    protected internal readonly Database Database = db;")
		b.W("")
	}
	serviceList(m, b)
	if len(m.PKs()) == 1 {
		b.W("")
		serviceGet(m, m.PKs()[0], b, args)
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

func serviceGet(m *model.Model, pk *model.Column, b *csharp.Block, args *model.Args) {
	b.W("    public async Task<Entities.%s?> Get(%s %s)", m.Proper(), ToCSharpType(pk), pk.Camel())
	b.W("    {")
	if len(m.Relations) > 0 {
		r := m.Relations[0]
		tgt := args.Models.Get(r.Table)
		inc := fmt.Sprintf(".Include(x => x.%s)", tgt.Proper())
		b.W("        return await Models%s.FirstOrDefaultAsync(x => x.%s == %s);", inc, pk.Proper(), pk.Camel())
	} else {
		b.W("        return await Models.FirstOrDefaultAsync(x => x.%s == %s);", pk.Proper(), pk.Camel())
	}
	b.W("    }")
}

func serviceDelete(m *model.Model, pk *model.Column, b *csharp.Block) {
	b.W("    public async Task<Entities.%s?> Delete(%s %s)", m.Proper(), ToCSharpType(pk), pk.Camel())
	b.W("    {")
	b.W("        var ret = await Get(%s);", pk.Camel())
	b.W("        if (ret != null)")
	b.W("        {")
	b.W("            Models.Remove(ret);")
	b.W("        }")
	b.W("")
	b.W("        return ret;")
	b.W("    }")
}
