package cmcp

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/mcpserver"
	"{{{ .Package }}}/views/vmcp"
)

func MCPPrompt(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.prompt", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, prompt, err := mcpPrompt(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Prompt [%s]", prompt.Name), prompt)
		return controller.Render(r, as, &vmcp.PromptDetail{Server: mcpx, Prompt: prompt}, ps, mcpBreadcrumb, "prompt", prompt.Name)
	})
}

func mcpPrompt(r *http.Request, as *app.State, ps *cutil.PageState) (*mcpserver.Server, *mcpserver.Prompt, error) {
	promptKey, _ := cutil.PathString(r, "prompt", false)
	mcpx, err := mcpserver.GetDefaultServer(ps.Context, as, ps.Logger)
	if err != nil {
		return nil, nil, err
	}
	var prompt *mcpserver.Prompt
	if promptKey != "" {
		prompt = mcpx.Prompts.Get(promptKey)
		if prompt == nil {
			return nil, nil, errors.Errorf("unable to find prompt [%s]", promptKey)
		}
	}
	return mcpx, prompt, nil
}
