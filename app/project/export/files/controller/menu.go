package controller

import (
	"fmt"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func Menu(args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile("cmenu", []string{"app", "controller", "cmenu"}, "generated")
	g.AddImport(helper.ImpAppMenu)
	g.AddBlocks(menuBlock(args))
	return g.Render(addHeader, linebreak)
}

func menuBlock(args *model.Args) *golang.Block {
	ret := golang.NewBlock("menu", "func")
	ret.W("func generatedMenu() menu.Items {")
	rct := menuContent(args)
	lo.ForEach(rct, func(x string, _ int) {
		ret.W(x)
	})
	if len(rct) == 0 {
		ret.W("\treturn nil")
	}
	ret.W("}")
	return ret
}

func menuContent(args *model.Args) []string {
	if len(args.Models) == 0 && len(args.Groups) == 0 {
		return nil
	}
	out := make([]string, 0, len(args.Models)+len(args.Groups))

	if len(args.Groups) == 0 && len(args.Models) == 0 {
		out = append(out, "\treturn menu.Items{}")
	} else {
		out = append(out, "\treturn menu.Items{")
		lo.ForEach(menuItemsFor(args.Groups, args.Models), func(x *menu.Item, _ int) {
			out = append(out, menuSerialize(x, "\t\t")...)
		})
		out = append(out, "\t}")
	}
	return out
}

func menuItemsFor(groups model.Groups, models model.Models) menu.Items {
	ret := make(menu.Items, 0, len(groups)+len(models))
	lo.ForEach(groups, func(g *model.Group, _ int) {
		ret = append(ret, menuItemForGroup(g, models))
	})
	lo.ForEach(models, func(m *model.Model, _ int) {
		if len(m.Group) == 0 {
			ret = append(ret, menuItemForModel(m, models))
		}
	})
	return ret
}

func menuItemForGroup(g *model.Group, models model.Models, pth ...string) *menu.Item {
	np := append(slices.Clone(pth), g.Key)
	ret := &menu.Item{Key: g.Key, Title: g.TitleSafe(), Description: g.Description, Icon: g.IconSafe()}
	lo.ForEach(g.Children, func(child *model.Group, _ int) {
		ret.Children = append(ret.Children, menuItemForGroup(child, models, np...))
	})
	matches := models.ForGroup(np...)
	lo.ForEach(matches, func(m *model.Model, _ int) {
		ret.Children = append(ret.Children, menuItemForModel(m, models))
	})
	return ret
}

func menuItemForModel(m *model.Model, models model.Models) *menu.Item {
	ret := &menu.Item{Key: m.Package, Title: m.TitlePlural(), Description: m.Description, Icon: m.Icon, Route: m.Route()}
	if len(m.GroupedColumns()) > 0 {
		lo.ForEach(m.GroupedColumns(), func(g *model.Column, _ int) {
			desc := fmt.Sprintf("%s from %s", g.ProperPlural(), m.Plural())
			kid := &menu.Item{Key: g.Camel(), Title: g.ProperPlural(), Description: desc, Icon: m.Icon, Route: m.Route() + "/" + g.Camel()}
			ret.Children = append(ret.Children, kid)
		})
	}
	lo.ForEach(models.ForGroup(append(slices.Clone(m.Group), m.Name)...), func(x *model.Model, _ int) {
		kid := menuItemForModel(x, models)
		ret.Children = append(ret.Children, kid)
	})
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
		lo.ForEach(m.Children, func(kid *menu.Item, _ int) {
			out = append(out, menuSerialize(kid, prefix+"\t")...)
		})
		out = append(out, prefix+"}},")
	}
	return out
}
