package action

import (
	"strings"
	"text/template"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/project"
	"github.com/pkg/errors"
)

func runTemplate(f *file.File, ctx *project.TemplateContext) (string, error) {
	t, err := template.New(f.FullPath()).Delims(delimStart, delimEnd).Parse(f.Content)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create template for [%s]", f.FullPath())
	}

	res := &strings.Builder{}
	err = t.Execute(res, ctx)
	if err != nil {
		return "", errors.Wrapf(err, "unable to execute template for [%s]", f.FullPath())
	}
	return res.String(), nil
}
