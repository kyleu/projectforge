package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

var ProjectForgeResource = &Resource{
	Name:        "projectforge",
	Description: "Description, purpose, and basic usage instructions for Project Forge (this MCP server)",
	URI:         "project://projectforge",
	MIMEType:    "text/markdown",
	Content:     "# Project Forge\n\n[Project Forge](https://projectforge.dev) is an application that allows you to generate, manage, and grow web applications built using the Go programming language.",
}

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
		return util.ToJSON(ret.Get(id)), nil
	}
	return util.ToJSON(ret), nil
}

var ProjectPrompt = &Prompt{
	Name:        "project_usage",
	Description: "A simple prompt that helps the system build, test, and work with an application managed by Project Forge" + util.AppName,
	Fn: func(ctx context.Context, as *app.State, req mcp.GetPromptRequest, args util.ValueMap, logger util.Logger) (string, error) {
		return `This application is written using Go, and is a web application managed by Project Forge.`, nil
	},
}
