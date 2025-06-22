package clib

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/mcpserver"
	"projectforge.dev/projectforge/views/vadmin"
)

const mcpBreadcrumb = "mcp"

func MCPIndex(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, _, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("MCP", mcpx)
		return controller.Render(r, as, &vadmin.MCPList{Server: mcpx}, ps, mcpBreadcrumb)
	})
}

func MCPServe(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.serve.streamable", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, _, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		mcpx.ServeHTTP(ps.Context, w, r, ps.Logger)
		return "", nil
	})
}

func MCPTool(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.tool", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, tool, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Tool [%s]", tool.Name), tool)
		return controller.Render(r, as, &vadmin.MCPDetail{Server: mcpx, Tool: tool}, ps, mcpBreadcrumb, tool.Name)
	})
}

func MCPToolRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.tool.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, tool, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		ret, err := tool.Fn(ps.Context, as, mcp.CallToolRequest{}, frm, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Tool [%s] Result", tool.Name), ret)
		return controller.Render(r, as, &vadmin.MCPDetail{Server: mcpx, Tool: tool, Args: frm, Result: ret}, ps, mcpBreadcrumb, tool.Name)
	})
}

func mcpTool(r *http.Request, as *app.State, ps *cutil.PageState) (*mcpserver.Server, *mcpserver.Tool, error) {
	toolKey, _ := cutil.PathString(r, "tool", false)
	mcpx, err := mcpserver.GetDefaultServer(ps.Context, as, ps.Logger)
	if err != nil {
		return nil, nil, err
	}
	var tool *mcpserver.Tool
	if toolKey != "" {
		tool = mcpx.Tools.Get(toolKey)
		if tool == nil {
			return nil, nil, errors.Errorf("unable to find tool [%s]", toolKey)
		}
	}
	return mcpx, tool, nil
}
