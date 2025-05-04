package controller

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Menu(args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile("cmenu", []string{"app", "controller", "cmenu"}, "generated")
	g.AddImport(helper.ImpAppMenu)
	models := args.Models.WithRoutes()
	groups, names, orphans := sortModels(models, args)
	if vBlock := menuBlockV(models, args, groups, names); len(vBlock.Lines) > 0 {
		g.AddBlocks(vBlock)
	}
	g.AddBlocks(menuBlockGM(args, orphans))
	return g.Render(linebreak)
}

func sortModels(models model.Models, args *metamodel.Args) (map[string][]string, []string, []string) {
	groups := map[string][]string{}
	names := make([]string, 0, len(models)+len(args.Groups))
	orphans := make([]string, 0)
	lo.ForEach(models.SortedDisplay(), func(m *model.Model, _ int) {
		n := m.ProperWithGroup()
		names = append(names, n)
		if len(m.Group) == 0 {
			orphans = append(orphans, n)
		} else {
			gn := m.Group[len(m.Group)-1]
			curr := groups[gn]
			groups[gn] = append(curr, helper.TextMenuItem+n)
		}
	})
	return groups, names, orphans
}

func menuBlockV(models model.Models, args *metamodel.Args, groups map[string][]string, names []string) *golang.Block {
	nameLength := util.StringArrayMaxLength(names)
	lines := lo.Map(models, func(m *model.Model, _ int) string {
		n := util.StringPad(m.ProperWithGroup(), nameLength)
		i := menuSerialize(menuItemForModel(m, models, args.Acronyms), "", true)
		return fmt.Sprintf("\t%s%s = %s", helper.TextMenuItem, n, util.StringJoin(i, "\n"))
	})
	slices.Sort(lines)

	flatGroups := args.Groups.Flatten()
	maxGroupLength := util.StringArrayMaxLength(flatGroups.Strings(""))
	if len(lines) > 0 && len(flatGroups) > 0 {
		lines = append(lines, "")
	}
	for _, grp := range flatGroups {
		if grp.Provided {
			continue
		}
		k, g := grp.Key, groups[grp.Key]
		n := util.StringToProper(k, args.Acronyms...)
		msg := fmt.Sprintf("\tmenuGroup%s = &menu.Item{Key: %q, Title: %q", util.StringPad(n, maxGroupLength), grp.Key, grp.String())
		if grp.Icon != "" {
			msg += fmt.Sprintf(", Icon: %q", grp.Icon)
		}
		if grp.Route != "" {
			msg += fmt.Sprintf(", Route: %q", grp.Route)
		}
		if len(grp.Children) > 0 || len(g) > 0 {
			items := append(slices.Clone(g), grp.Children.Strings("menuGroup")...)
			msg += fmt.Sprintf(", Children: menu.Items{%s}", util.StringJoin(items, ", "))
		}
		msg += "}"
		lines = append(lines, msg)
	}

	v := golang.NewBlock("items", "var")
	switch len(lines) {
	case 0:
		// noop
	case 1:
		v.W("var " + strings.TrimSpace(lines[0]))
	default:
		v.W("var (")
		lo.ForEach(lines, func(l string, _ int) {
			v.W(l)
		})
		v.W(")")
	}

	return v
}

func menuBlockGM(args *metamodel.Args, orphans []string) *golang.Block {
	gm := golang.NewBlock("generatedMenu", "func")
	// gm.Lints = append(gm.Lints, "unused")
	gm.W("func generatedMenu() menu.Items {")
	if len(args.Groups) == 0 && len(orphans) == 0 {
		gm.W("\treturn menu.Items{}")
	} else {
		gm.W("\treturn menu.Items{")
		for _, g := range args.Groups {
			if !g.Provided {
				gm.WF("\t\tmenuGroup%s,", util.StringToProper(g.Proper(), args.Acronyms...))
			}
		}
		for _, o := range orphans {
			gm.WF("\t\t%s%s,", helper.TextMenuItem, o)
		}
		gm.W("\t}")
	}
	gm.W("}")

	return gm
}

func menuItemForModel(m *model.Model, models model.Models, acronyms []string) *menu.Item {
	w := m.ProperWithGroup()
	ret := &menu.Item{
		Key: m.Package, Title: m.TitlePlural(), Description: m.Description, Icon: m.Icon, Route: m.Route(), Hidden: m.HasTag("menu-hidden"), Warning: w,
	}
	lo.ForEach(models.ForGroup(append(slices.Clone(m.Group), m.Package)...), func(x *model.Model, _ int) {
		kid := menuItemForModel(x, models, acronyms)
		ret.Children = append(ret.Children, kid)
	})
	return ret
}

func menuSerialize(m *menu.Item, prefix string, top bool) []string {
	var ws string
	if !top {
		ws = "\t"
	}
	var out []string
	var rt string
	if m.Route != "" {
		rt += fmt.Sprintf(", Route: %q", "/"+m.Route)
	}
	if m.Hidden {
		rt += ", Hidden: true"
	}
	var desc string
	if m.Description != "" {
		desc = fmt.Sprintf("Description: %q, ", m.Description)
	}
	args := fmt.Sprintf("Key: %q, Title: %q, %sIcon: %q%s", m.Key, m.Title, desc, m.Icon, rt)
	if len(m.Children) == 0 {
		out = append(out, ws+prefix+fmt.Sprintf("&menu.Item{%s}", args))
	} else {
		kids := lo.Map(m.Children, func(kid *menu.Item, _ int) string {
			return helper.TextMenuItem + kid.Warning
		})
		out = append(out, ws+prefix+"&menu.Item{"+args+", Children: menu.Items{"+util.StringJoin(kids, ", ")+"}}")
	}
	return out
}
