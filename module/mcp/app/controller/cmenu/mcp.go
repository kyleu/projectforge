package cmenu

import (
	"{{{ .Package }}}/app/lib/mcpserver"
	"{{{ .Package }}}/app/lib/menu"
)

func mcpMenu() *menu.Item {
	ret := &menu.Item{Key: "mcp", Title: "Model Context Protocol", Icon: "robot", Route: "/mcp"}
	if x := mcpserver.CurrentDefaultServer(); x != nil {
		for _, t := range x.Tools {
			ret.Children = append(ret.Children, &menu.Item{Key: t.Name, Title: t.Name, Description: t.Description, Icon: t.IconSafe(), Route: "/mcp/tool/" + t.Name})
		}
	}
	return ret
}
