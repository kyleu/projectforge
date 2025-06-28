package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/samber/lo"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type PromptHandler func(ctx context.Context, as *app.State, req mcp.GetPromptRequest, args util.ValueMap, logger util.Logger) (string, error)

type Prompt struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Content     string `json:"content,omitempty"`
}

func NewStaticPrompt(name string, description string, content string) *Prompt {
	return &Prompt{Name: name, Description: description, Content: content}
}

func (p *Prompt) ToMCP() (mcp.Prompt, error) {
	return mcp.NewPrompt(p.Name, mcp.WithPromptDescription(p.Description)), nil
}

func (p *Prompt) IconSafe() string {
	return util.Choose(p.Icon == "", "gift", p.Icon)
}

func (p *Prompt) Handler(as *app.State, logger util.Logger) server.PromptHandlerFunc {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := util.ValueMap{}
		for k, v := range req.Params.Arguments {
			args[k] = v
		}
		return &mcp.GetPromptResult{
			Description: p.Description,
			Messages:    []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: p.Content}}},
		}, nil
	}
}

type Prompts []*Prompt

func (r Prompts) Get(n string) *Prompt {
	return lo.FindOrElse(r, nil, func(x *Prompt) bool {
		return x.Name == n
	})
}

func (s *Server) AddPrompts(as *app.State, logger util.Logger, prompts ...*Prompt) error {
	for _, r := range prompts {
		s.Prompts = append(s.Prompts, r)
		m, err := r.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddPrompt(m, r.Handler(as, logger))
	}
	return nil
}
