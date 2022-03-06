package action

import (
	"strings"
	"text/template"

	"projectforge.dev/app/project"
	"github.com/pkg/errors"
	"projectforge.dev/app/file"
)

func runTemplate(path string, content string, ctx *project.TemplateContext) (string, error) {
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

func runTemplateFile(f *file.File, ctx *project.TemplateContext) (string, error) {
	return runTemplate(f.FullPath(), f.Content, ctx)
}
