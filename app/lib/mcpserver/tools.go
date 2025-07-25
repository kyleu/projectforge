package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

var ListProjectsTool = &Tool{
	Name:        "list_projects",
	Description: "List the projects managed by " + util.AppName,
	Fn:          listProjectsHandler,
}

var GetProjectTool = &Tool{
	Name:        "get_project",
	Description: "Get details of a specific project managed by " + util.AppName,
	Args: util.FieldDescs{
		{Key: "id", Description: "Optional project id. If omitted, all projects will be returned"},
	},
	Fn: getProjectHandler,
}

func listProjectsHandler(_ context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (any, error) {
	return as.Services.Projects.Projects(), nil
}

func getProjectHandler(_ context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (any, error) {
	id, _ := args.GetString("id", true)
	return as.Services.Projects.Projects().Get(id), nil
}
