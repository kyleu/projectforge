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

const mcpBreadcrumb = "MCP||/admin/mcp**graph"

func MCPIndex(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.index", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, _, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("MCP", mcpx)
		return controller.Render(r, as, &vadmin.MCP{Tools: mcpx.Tools}, ps, keyAdmin, mcpBreadcrumb)
	})
}

func MCPTask(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.task", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, tool, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Tool [%s]", tool.Name), tool)
		return controller.Render(r, as, &vadmin.MCP{Tools: mcpx.Tools, Tool: tool}, ps, keyAdmin, mcpBreadcrumb, tool.Name)
	})
}

func MCPTaskRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.task", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, tool, err := mcpTool(r, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		ret, err := tool.Fn(ps.Context, frm, mcp.CallToolRequest{})
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Tool [%s] Result", tool.Name), ret)
		return controller.Render(r, as, &vadmin.MCP{Tools: mcpx.Tools, Tool: tool, Args: frm, Result: ret}, ps, keyAdmin, mcpBreadcrumb, tool.Name)
	})
}

func mcpTool(r *http.Request, as *app.State, ps *cutil.PageState) (*mcpserver.Server, *mcpserver.Tool, error) {
	mcpserver.InitMCP(as.BuildInfo, as.Debug)
	toolKey, _ := cutil.PathString(r, "tool", false)
	mcpx, err := mcpserver.NewServer(ps.Context, as.BuildInfo.Version)
	if err != nil {
		return nil, nil, err
	}
	var tool *mcpserver.Tool
	if toolKey != "" {
		tool = mcpx.Tools[toolKey]
		if tool == nil {
			return nil, nil, errors.Errorf("unable to find tool [%s]", toolKey)
		}
	}
	return mcpx, tool, nil
}
