package csfiles

import (
	"fmt"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/csharp"
	"projectforge.dev/projectforge/app/project/export/model"
)

func menu(ns string, models model.Models) (*file.File, error) {
	fn := "Menu"
	pkg := "Util"
	pth := []string{ns}
	if ns == "Shared" {
		fn += "Shared"
		pth = append(pth, "Controllers")
		pkg = "Controllers"
	} else {
		pth = append(pth, "Util")
	}
	f := csharp.NewFile(ns+"."+pkg, pth, fn+".cs")
	if ns != "Shared" {
		f.AddImport(ImpControllers)
	}
	b := csharp.NewBlock("Menu", "class")
	b.W("public static class %s", fn)
	b.W("{")
	b.W("    public static IEnumerable<MenuItem> GetItems()")
	b.W("    {")
	b.W("        return")
	b.W("        [")
	for _, m := range models {
		const msg = "            new MenuItem { Key = %q, Title = %q, Description = %q, Url = %q, Icon = %q },"
		desc := fmt.Sprintf("The %s in our system", m.TitlePluralLower())
		b.W(msg, m.Name, m.TitlePlural(), desc, m.CSRoute(), m.Name)
	}
	b.W("            new MenuItem()")
	b.W("        ];")
	b.W("    }")
	b.W("}")
	f.AddBlocks(b)
	return f.Render()
}
