package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

var ExampleTool = &Tool{
	Name:        "example_mcp_server",
	Description: "Returns a random integer",
	Args:        util.FieldDescs{{Key: "max", Description: "Maximum possible random int (exclusive), defaults to 100", Default: "100"}},
	Fn:          exampleHandler,
}

func exampleHandler(ctx context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (any, error) {
	mx, err := args.GetInt("max", true)
	if err != nil {
		return "", errors.Errorf("argument [max] must be an integer")
	}
	if mx == 0 {
		mx = 100
	}
	return util.RandomInt(mx), nil
}
