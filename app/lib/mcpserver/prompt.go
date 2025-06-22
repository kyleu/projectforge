package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

type PromptHandler func(ctx context.Context, as *app.State, req mcp.GetPromptRequest, args util.ValueMap, logger util.Logger) (string, error)

type Prompt struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Icon        string          `json:"icon,omitempty"`
	Args        util.FieldDescs `json:"args,omitempty"`
	Fn          PromptHandler   `json:"-"`
}

func NewStaticPrompt(name string, description string, content string) *Prompt {
	fn := func(ctx context.Context, as *app.State, req mcp.GetPromptRequest, args util.ValueMap, logger util.Logger) (string, error) {
		return content, nil
	}
	return &Prompt{Name: name, Description: description, Fn: fn}
}

func (p *Prompt) ToMCP() (mcp.Prompt, error) {
	opts := []mcp.PromptOption{mcp.WithPromptDescription(p.Description)}
	for _, x := range p.Args {
		opts = append(opts, mcp.WithArgument(x.Key, mcp.RequiredArgument(), mcp.ArgumentDescription(x.Description)))
	}
	return mcp.NewPrompt(p.Name, opts...), nil
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
		if p.Fn == nil {
			return nil, errors.Errorf("no handler for prompt [%s]", p.Name)
		}
		ret, err := p.Fn(ctx, as, req, args, logger)
		if err != nil {
			return nil, errors.Errorf("errors running prompt [%s] with arguments %s: %+v", p.Name, util.ToJSONCompact(args), err)
		}
		return &mcp.GetPromptResult{
			Description: p.Description,
			Messages:    []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: ret}}},
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
