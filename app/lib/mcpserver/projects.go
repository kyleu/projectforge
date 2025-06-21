package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var ListProjectsTool = &Tool{
	Name:        "list_projects",
	Description: "List the projects managed by " + util.AppName,
	Fn:          projectsHandler,
}

var GetProjectTool = &Tool{
	Name:        "get_project",
	Description: "Get details of a specific project managed by " + util.AppName,
	Args: util.FieldDescs{
		{Key: "id", Description: "Optional project id. If omitted, all projects will be returned"},
	},
	Fn: projectsHandler,
}

func projectsHandler(_ context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (string, error) {
	id, _ := args.GetString("id", true)
	svc := as.Services.Projects
	ret, err := svc.Refresh(logger)
	if err != nil {
		return "", err
	}
	if id != "" {
		ret = project.Projects{ret.Get(id)}
	}
	return util.ToJSON(ret), nil
}
