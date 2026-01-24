package action

import (
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	projecttemplate "projectforge.dev/projectforge/app/project/template"
)

func runTemplate(path string, content string, ctx *projecttemplate.Context) (string, error) {
	t, err := template.New(path).Delims(delimStart, delimEnd).Parse(content)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create template for [%s]", path)
	}

	res := &strings.Builder{}
	err = t.Execute(res, ctx)
	if err != nil {
		return "", errors.Wrapf(err, "unable to execute template for [%s]", path)
	}
	return res.String(), nil
}

func runTemplateFile(f *file.File, ctx *projecttemplate.Context) (string, error) {
	return runTemplate(f.FullPath(), f.Content, ctx)
}
