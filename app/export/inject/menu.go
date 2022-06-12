package inject

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/menu"
)

const msg = `Key: %q, Title: %q, Description: %q, Icon: %q, Route: "/%s"`

func Menu(f *file.File, args *model.Args) error {
	if len(args.Models) == 0 && len(args.Groups) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models)+len(args.Groups))

	var adminGroups, nonAdminGroups model.Groups
	for _, g := range args.Groups {
		if g.HasTag("public") {
			nonAdminGroups = append(nonAdminGroups, g)
		} else {
			adminGroups = append(adminGroups, g)
		}
	}

	var adminModels, nonAdminModels model.Models
	for _, m := range args.Models {
		if m.HasTag("public") {
			nonAdminModels = append(nonAdminModels, m)
		} else {
			adminModels = append(adminModels, m)
		}
	}

	if len(nonAdminGroups) > 0 || len(nonAdminModels) > 0 {
		out = append(out, "ret = append(ret,")
		for _, x := range menuItemsFor(nonAdminGroups, nonAdminModels) {
			out = append(out, menuSerialize(x, "\t\t")...)
		}
		out = append(out, ")")
	}

	if len(adminGroups) > 0 || len(adminModels) > 0 {
		out = append(out, "if isAdmin {", "\tret = append(ret,")
		for _, x := range menuItemsFor(adminGroups, adminModels) {
			out = append(out, menuSerialize(x, "\t\t")...)
		}
		out = append(out, "\t)", "}")
	}

	content := map[string]string{"codegen": "\n\t" + strings.Join(out, "\n\t") + "\n\t// "}
	return file.Inject(f, content)
}

func menuItemsFor(groups model.Groups, models model.Models) menu.Items {
	ret := make(menu.Items, 0, len(groups)+len(models))
	for _, g := range groups {
		ret = append(ret, menuItemForGroup(g))
	}
	for _, m := range models {
		ret = append(ret, menuItemForModel(m))
	}
	return ret
}

func menuItemForGroup(g *model.Group) *menu.Item {
	ret := &menu.Item{Key: g.Key, Title: g.TitleSafe(), Description: g.Description, Icon: g.IconSafe()}
	for _, child := range g.Children {
		ret.Children = append(ret.Children, menuItemForGroup(child))
	}
	return ret
}

func menuItemForModel(m *model.Model) *menu.Item {
	ret := &menu.Item{Key: m.Package, Title: m.TitlePlural(), Description: m.Description, Icon: m.Icon, Route: m.Route()}
	if len(m.GroupedColumns()) > 0 {
		for _, g := range m.GroupedColumns() {
			desc := fmt.Sprintf("%s from %s", g.ProperPlural(), m.Plural())
			kid := &menu.Item{Key: g.Camel(), Title: g.ProperPlural(), Description: desc, Icon: m.Icon, Route: m.Route() + "/" + g.Camel()}
			ret.Children = append(ret.Children, kid)
		}
	}
	return ret
}

func menuSerialize(m *menu.Item, prefix string) []string {
	var out []string
	var rt string
	if m.Route != "" {
		rt = fmt.Sprintf(", Route: %q", "/"+m.Route)
	}
	var desc string
	if m.Description != "" {
		desc = fmt.Sprintf("Description: %q, ", m.Description)
	}
	args := fmt.Sprintf("Key: %q, Title: %q, %sIcon: %q%s", m.Key, m.Title, desc, m.Icon, rt)
	if len(m.Children) == 0 {
		out = append(out, prefix+"&menu.Item{"+args+"},")
	} else {
		out = append(out, prefix+"&menu.Item{"+args+", Children: menu.Items{")
		for _, kid := range m.Children {
			out = append(out, menuSerialize(kid, prefix+"\t")...)
		}
		out = append(out, prefix+"}},")
	}
	return out
}
