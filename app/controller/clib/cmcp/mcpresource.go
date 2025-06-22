package cmcp

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/mcpserver"
	"projectforge.dev/projectforge/views/vmcp"
)

func MCPResource(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.resource", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, resource, err := mcpResource(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Resource [%s]", resource.Name), resource)
		return controller.Render(r, as, &vmcp.ResourceDetail{Server: mcpx, Resource: resource}, ps, mcpBreadcrumb, resource.Name)
	})
}

func mcpResource(r *http.Request, as *app.State, ps *cutil.PageState) (*mcpserver.Server, *mcpserver.Resource, error) {
	resourceKey, _ := cutil.PathString(r, "resource", false)
	mcpx, err := mcpserver.GetDefaultServer(ps.Context, as, ps.Logger)
	if err != nil {
		return nil, nil, err
	}
	var resource *mcpserver.Resource
	if resourceKey != "" {
		resource = mcpx.Resources.Get(resourceKey)
		if resource == nil {
			return nil, nil, errors.Errorf("unable to find resource [%s]", resourceKey)
		}
	}
	return mcpx, resource, nil
}
