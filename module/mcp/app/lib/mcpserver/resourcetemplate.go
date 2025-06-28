package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type ResourceTemplateHandler func(ctx context.Context, as *app.State, req mcp.ReadResourceRequest, args util.ValueMap, logger util.Logger) (string, string, any, error)

var ResourceTemplateArgs = util.FieldDescs{{Key: "uri", Description: "URI to request", Type: "string"}}

type ResourceTemplate struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Icon        string                  `json:"icon,omitempty"`
	URI         string                  `json:"uri"`
	MIMEType    string                  `json:"mimeType"`
	Args        util.FieldDescs         `json:"args,omitempty"`
	Fn          ResourceTemplateHandler `json:"-"`
}

func NewResourceTemplate(name string, description string, uri string, mimeType string, content string) *Resource {
	if mimeType == "" {
		mimeType = "application/json"
	}
	return &Resource{Name: name, Description: description, URI: uri, MIMEType: mimeType, Content: content}
}

func (r *ResourceTemplate) ToMCP() (mcp.ResourceTemplate, error) {
	return mcp.NewResourceTemplate(r.URI, r.Name, mcp.WithTemplateDescription(r.Description)), nil
}

func (r *ResourceTemplate) IconSafe() string {
	return util.Choose(r.Icon == "", "folder", r.Icon)
}

func (r *ResourceTemplate) Handler(as *app.State, logger util.Logger) server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		var ret []mcp.ResourceContents
		args := util.ValueMapFrom(req.Params.Arguments)
		u, mt, content, err := r.Fn(ctx, as, req, util.ValueMapFrom(args), logger)
		if err != nil {
			return nil, errors.Errorf("error running resource template [%s] with arguments %s: %+v", r.Name, util.ToJSONCompact(args), err)
		}
		trc := mcp.TextResourceContents{URI: util.Choose(u == "", r.URI, u), MIMEType: util.Choose(mt == "", r.MIMEType, mt)}
		switch t := content.(type) {
		case string:
			trc.Text = t
		case []byte:
			trc.Text = string(t)
		default:
			trc.Text = util.ToJSON(t)
		}
		ret = append(ret, trc)
		return ret, nil
	}
}

type ResourceTemplates []*ResourceTemplate

func (r ResourceTemplates) Get(n string) *ResourceTemplate {
	return lo.FindOrElse(r, nil, func(x *ResourceTemplate) bool {
		return x.Name == n
	})
}

func (s *Server) AddResourceTemplates(as *app.State, logger util.Logger, resources ...*ResourceTemplate) error {
	for _, r := range resources {
		s.ResourceTemplates = append(s.ResourceTemplates, r)
		m, err := r.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddResourceTemplate(m, r.Handler(as, logger))
	}
	return nil
}
