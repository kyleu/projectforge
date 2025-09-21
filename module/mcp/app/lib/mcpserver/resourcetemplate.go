package mcpserver

import (
	"context"
	"encoding/base64"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type (
	ResourceReq             mcp.ReadResourceRequest
	ResourceTemplateHandler func(ctx context.Context, as *app.State, req ResourceReq, args util.ValueMap, logger util.Logger) (string, string, any, error)
	ResourceTemplate        struct {
		Name        string                  `json:"name"`
		Description string                  `json:"description,omitzero"`
		Icon        string                  `json:"icon,omitzero"`
		URI         string                  `json:"uri"`
		MIMEType    string                  `json:"mimeType"`
		Args        util.FieldDescs         `json:"args,omitempty"`
		Fn          ResourceTemplateHandler `json:"-"`
	}
)

func NewResourceTemplate(name string, description string, uri string, mimeType string, content string) *Resource {
	if mimeType == "" {
		mimeType = util.MIMETypeJSON
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
		u, mt, content, err := r.Fn(ctx, as, ResourceReq(req), util.ValueMapFrom(args), logger)
		if err != nil {
			return nil, errors.Errorf("error running resource template [%s] with arguments %s: %+v", r.Name, util.ToJSONCompact(args), err)
		}
		var rc mcp.ResourceContents
		u = util.Choose(u == "", r.URI, u)
		mt = util.Choose(mt == "", r.MIMEType, mt)
		switch t := content.(type) {
		case string:
			rc = mcp.TextResourceContents{URI: u, MIMEType: mt, Text: t}
		case []byte:
			rc = mcp.BlobResourceContents{URI: u, MIMEType: mt, Blob: base64.StdEncoding.EncodeToString(t)}
		default:
			rc = mcp.TextResourceContents{URI: u, MIMEType: mt, Text: util.ToJSON(t)}
		}
		ret = append(ret, rc)
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
