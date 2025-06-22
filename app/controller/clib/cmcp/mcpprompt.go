package cmcp

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/mcpserver"
	"projectforge.dev/projectforge/views/vmcp"
)

func MCPPrompt(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.prompt", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, prompt, err := mcpPrompt(r, as, ps)
		if err != nil {
			return "", err
		}
		var ret string
		if len(prompt.Args) == 0 {
			ret, err = prompt.Fn(ps.Context, as, mcp.GetPromptRequest{}, nil, ps.Logger)
			if err != nil {
				return "", err
			}
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Prompt [%s]", prompt.Name), prompt)
		return controller.Render(r, as, &vmcp.PromptDetail{Server: mcpx, Prompt: prompt, Result: ret}, ps, mcpBreadcrumb, prompt.Name)
	})
}

func MCPPromptRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("mcp.prompt.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mcpx, prompt, err := mcpPrompt(r, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		ret, err := prompt.Fn(ps.Context, as, mcp.GetPromptRequest{}, frm, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("MCP Prompt [%s] Result", prompt.Name), ret)
		return controller.Render(r, as, &vmcp.PromptDetail{Server: mcpx, Prompt: prompt, Args: frm, Result: ret}, ps, mcpBreadcrumb, prompt.Name)
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
