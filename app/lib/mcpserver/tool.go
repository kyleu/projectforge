package mcpserver

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

type ToolHandler func(ctx context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (any, error)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Icon        string          `json:"icon,omitempty"`
	Args        util.FieldDescs `json:"args,omitempty"`
	Fn          ToolHandler     `json:"-"`
}

func (t *Tool) ToMCP() (mcp.Tool, error) {
	opts := []mcp.ToolOption{mcp.WithDescription(t.Description)}
	for _, x := range t.Args {
		argOpts := []mcp.PropertyOption{mcp.Description(x.Description)}
		if x.Default == "" {
			argOpts = append(argOpts, mcp.Required())
		} else {
			argOpts = append(argOpts, mcp.DefaultString(x.Default))
		}
		switch x.Type {
		case "string", "":
			opts = append(opts, mcp.WithString(x.Key, argOpts...))
		case "float", "float64", "int", "int64", "number":
			opts = append(opts, mcp.WithNumber(x.Key, argOpts...))
		case "bool", "boolean":
			opts = append(opts, mcp.WithBoolean(x.Key, argOpts...))
		default:
			return mcp.Tool{}, errors.Errorf("unable to parse tool argument [%s] as type [%s]", x.Key, x.Type)
		}
	}
	return mcp.NewTool(t.Name, opts...), nil
}

func (t *Tool) Handler(as *app.State, logger util.Logger) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := util.ValueMapFrom(req.GetArguments())
		if t.Fn == nil {
			return nil, errors.Errorf("no handler for tool [%s]", t.Name)
		}
		ret, err := t.Fn(ctx, as, req, args, logger)
		if err != nil {
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("errors running tool [%s] with arguments %s", t.Name, util.ToJSONCompact(args)), err), nil
		}
		return mcp.NewToolResultText(valToText(ret)), nil
	}
}

func (t *Tool) IconSafe() string {
	return util.Choose(t.Icon == "", "cog", t.Icon)
}

type Tools []*Tool

func (t Tools) Get(n string) *Tool {
	return lo.FindOrElse(t, nil, func(x *Tool) bool {
		return x.Name == n
	})
}

func (s *Server) AddTools(as *app.State, logger util.Logger, tools ...*Tool) error {
	for _, tl := range tools {
		s.Tools = append(s.Tools, tl)
		m, err := tl.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddTool(m, tl.Handler(as, logger))
	}
	return nil
}
