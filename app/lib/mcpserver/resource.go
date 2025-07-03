package mcpserver

import (
	"context"
	"encoding/base64"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

type Resource struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	URI         string `json:"uri"`
	MIMEType    string `json:"mimeType"`
	Content     string `json:"content"`
	Binary      bool   `json:"binary,omitempty"`
}

func NewTextResource(name string, description string, uri string, mimeType string, content string) *Resource {
	if mimeType == "" {
		mimeType = util.MIMETypeJSON
	}
	return &Resource{Name: name, Description: description, URI: uri, MIMEType: mimeType, Content: content}
}

func NewBlobResource(name string, description string, uri string, mimeType string, content []byte) *Resource {
	if mimeType == "" {
		mimeType = util.MIMETypeJSON
	}
	out := base64.StdEncoding.EncodeToString(content)
	return &Resource{Name: name, Description: description, URI: uri, MIMEType: mimeType, Content: out, Binary: true}
}

func (r *Resource) ToMCP() (mcp.Resource, error) {
	return mcp.NewResource(r.URI, r.Name, mcp.WithResourceDescription(r.Description)), nil
}

func (r *Resource) IconSafe() string {
	return util.Choose(r.Icon == "", "file", r.Icon)
}

func (r *Resource) Extension() string {
	return util.ExtensionFromMIME(r.MIMEType)
}

func (r *Resource) Handler(as *app.State, logger util.Logger) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		var ret []mcp.ResourceContents
		if r.Binary {
			ret = append(ret, mcp.BlobResourceContents{URI: r.URI, MIMEType: r.MIMEType, Blob: r.Content})
		} else {
			ret = append(ret, mcp.TextResourceContents{URI: r.URI, MIMEType: r.MIMEType, Text: r.Content})
		}
		return ret, nil
	}
}

type Resources []*Resource

func (r Resources) Get(n string) *Resource {
	return lo.FindOrElse(r, nil, func(x *Resource) bool {
		return x.Name == n
	})
}

func (s *Server) AddResources(as *app.State, logger util.Logger, resources ...*Resource) error {
	for _, r := range resources {
		s.Resources = append(s.Resources, r)
		m, err := r.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddResource(m, r.Handler(as, logger))
	}
	return nil
}
