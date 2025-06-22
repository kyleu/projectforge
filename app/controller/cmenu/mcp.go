package cmenu

import (
	"projectforge.dev/projectforge/app/lib/mcpserver"
	"projectforge.dev/projectforge/app/lib/menu"
)

func mcpMenu() *menu.Item {
	ret := &menu.Item{Key: "mcp", Title: "Model Context Protocol", Icon: "robot", Route: "/mcp"}
	if x := mcpserver.CurrentDefaultServer(); x != nil {
		for _, r := range x.Resources {
			ret.Children = append(ret.Children, &menu.Item{Key: r.Name, Title: r.Name, Description: r.Description, Icon: r.IconSafe(), Route: "/mcp/resource/" + r.Name})
		}
		for _, t := range x.Tools {
			ret.Children = append(ret.Children, &menu.Item{Key: t.Name, Title: t.Name, Description: t.Description, Icon: t.IconSafe(), Route: "/mcp/tool/" + t.Name})
		}
		for _, p := range x.Prompts {
			ret.Children = append(ret.Children, &menu.Item{Key: p.Name, Title: p.Name, Description: p.Description, Icon: p.IconSafe(), Route: "/mcp/prompt/" + p.Name})
		}
	}
	return ret
}
