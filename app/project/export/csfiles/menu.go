package csfiles

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func menu(models model.Models, p *project.Project) (*file.File, error) {
	f := csharp.NewFile(p.Package+".Util", []string{"Util"}, "Menu.cs")
	f.AddImport(ImpControllers)
	b := csharp.NewBlock("Menu", "class")
	b.W("public static class Menu")
	b.W("{")
	b.W("    public static List<MenuItem> GetItems()")
	b.W("    {")
	b.W("        return [")
	for _, m := range models {
		b.W("            new MenuItem { Key = %q, Title = %q, Url = %q, Icon = %q },", m.Name, m.TitlePlural(), "/"+m.Plural(), m.Name)
	}
	b.W("        ];")
	b.W("    }")
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}
