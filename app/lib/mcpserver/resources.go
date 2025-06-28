package mcpserver

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

const uriTemplate = "project-file://%s/%s"

var ProjectContentResource = &ResourceTemplate{
	Name:        "project_content",
	Description: "Get details of a specific file within a project managed by " + util.AppName,
	URI:         fmt.Sprintf(uriTemplate, "{id}", "{path*}"),
	Args: util.FieldDescs{
		{Key: "id", Description: "Project ID", Type: "string"},
		{Key: "path", Description: "Path to the file within the project"},
	},
	Fn: projectContentHandler,
}

func projectContentHandler(_ context.Context, as *app.State, req mcp.ReadResourceRequest, args util.ValueMap, logger util.Logger) (string, string, any, error) {
	id, err := args.GetString("id", false)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "must provide [id] argument")
	}
	pth, err := args.GetString("path", false)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "must provide [path] argument")
	}
	prj := as.Services.Projects.Projects().Get(id)
	if prj == nil {
		return "", "", nil, errors.Errorf("project [%s] not found", id)
	}
	fs, err := as.Services.Projects.GetFilesystem(prj)
	if err != nil {
		return "", "", nil, errors.Wrapf(err, "can't create [%s] filesystem", id)
	}
	b, err := fs.ReadFile(pth)
	if err != nil {
		return "", "", nil, errors.Wrapf(err, "unable to read [%s] file at [%s]", id, pth)
	}
	ret := string(b)
	u := fmt.Sprintf(uriTemplate, id, pth)
	ext := filepath.Ext(u)
	mt := util.MIMEFromExtension(ext)
	if mt == "" {
		mt = "text/plain"
	}
	return u, mt, ret, nil
}
