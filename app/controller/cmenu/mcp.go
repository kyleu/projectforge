package cmenu

import (
	"projectforge.dev/projectforge/app/lib/mcpserver"
	"projectforge.dev/projectforge/app/lib/menu"
)

func MCPMenu() *menu.Item {
	ret := &menu.Item{Key: "mcp", Title: "Model Context Protocol", Icon: "robot", Route: "/mcp"}
	if x := mcpserver.CurrentDefaultServer(); x != nil {
		if len(x.Resources) > 0 {
			kid := &menu.Item{Key: "resource", Title: "Resources", Description: "Static resources for your AI to use", Icon: "file"}
			for _, r := range x.Resources {
				i := &menu.Item{Key: r.Name, Title: r.Name, Description: r.Description, Icon: r.IconSafe(), Route: "/mcp/resource/" + r.Name}
				kid.Children = append(kid.Children, i)
			}
			ret.Children = append(ret.Children, kid)
		}
		if len(x.ResourceTemplates) > 0 {
			kid := &menu.Item{Key: "resourcetemplate", Title: "Resource Templates", Description: "Dynamic resources for your AI use", Icon: "folder"}
			for _, rt := range x.ResourceTemplates {
				i := &menu.Item{Key: rt.Name, Title: rt.Name, Description: rt.Description, Icon: rt.IconSafe(), Route: "/mcp/resourcetemplate/" + rt.Name}
				kid.Children = append(kid.Children, i)
			}
			ret.Children = append(ret.Children, kid)
		}
		if len(x.Tools) > 0 {
			kid := &menu.Item{Key: "tool", Title: "Tools", Description: "Tools your AI can call", Icon: "cog"}
			for _, t := range x.Tools {
				i := &menu.Item{Key: t.Name, Title: t.Name, Description: t.Description, Icon: t.IconSafe(), Route: "/mcp/tool/" + t.Name}
				kid.Children = append(kid.Children, i)
			}
			ret.Children = append(ret.Children, kid)
		}
		if len(x.Prompts) > 0 {
			kid := &menu.Item{Key: "prompt", Title: "Prompts", Description: "Prompts for your AI to use", Icon: "gift"}
			for _, p := range x.Prompts {
				i := &menu.Item{Key: p.Name, Title: p.Name, Description: p.Description, Icon: p.IconSafe(), Route: "/mcp/prompt/" + p.Name}
				kid.Children = append(kid.Children, i)
			}
			ret.Children = append(ret.Children, kid)
		}
	}
	return ret
}
