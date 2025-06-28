package cmcp

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/mcpserver"
	"{{{ .Package }}}/views/vmcp"
)

func MCPResourceTemplate(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.resourcetemplate", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, rt, err := mcpResourceTemplate(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP ResourceTemplate [%s]", rt.Name), rt)
		return controller.Render(r, as, &vmcp.ResourceTemplateDetail{Server: mcpx, ResourceTemplate: rt}, ps, mcpBreadcrumb, "resourcetemplate", rt.Name)
	})
}

func MCPResourceTemplateRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.resourcetemplate.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, rt, err := mcpResourceTemplate(r, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		params := mcp.ReadResourceParams{URI: rt.URI, Arguments: frm}
		u, mt, ret, err := rt.Fn(ps.Context, as, mcp.ReadResourceRequest{Params: params}, frm, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Resource Template [%s] Result", rt.Name), ret)
		page := &vmcp.ResourceTemplateDetail{Server: mcpx, ResourceTemplate: rt, Args: frm, URI: u, MIMEType: mt, Result: ret}
		return controller.Render(r, as, page, ps, mcpBreadcrumb, "resourcetemplate", rt.Name)
	})
}

func mcpResourceTemplate(r *http.Request, as *app.State, ps *cutil.PageState) (*mcpserver.Server, *mcpserver.ResourceTemplate, error) {
	rtKey, _ := cutil.PathString(r, "rt", false)
	mcpx, err := mcpserver.GetDefaultServer(ps.Context, as, ps.Logger)
	if err != nil {
		return nil, nil, err
	}
	var rt *mcpserver.ResourceTemplate
	if rtKey != "" {
		rt = mcpx.ResourceTemplates.Get(rtKey)
		if rt == nil {
			return nil, nil, errors.Errorf("unable to find resource template [%s]", rtKey)
		}
	}
	return mcpx, rt, nil
}
