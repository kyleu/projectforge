package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var ProjectsTool = &Tool{
	Name:        "projects",
	Description: "Returns the projects managed by " + util.AppName,
	Args: util.FieldDescs{
		{Key: "id", Description: "Optional project id. If omitted, all projects will be returned"},
	},
	Fn: projectsHandler,
}

func projectsHandler(_ context.Context, args util.ValueMap, _ mcp.CallToolRequest) (string, error) {
	id, _ := args.GetString("id", true)
	logger, _ := log.InitLogging(false)
	svc := project.NewService()
	ret, err := svc.Refresh(logger)
	if err != nil {
		return "", err
	}
	if id != "" {
		ret = project.Projects{ret.Get(id)}
	}
	return util.ToJSON(ret), nil
}
