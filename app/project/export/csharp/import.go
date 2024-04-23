package csharp

import "projectforge.dev/projectforge/app/util"

type Imports []string

func (i Imports) Render(linebreak string) string {
	ret := &util.StringSlice{}
	for _, imp := range i {
		ret.Pushf("using %s;", imp)
	}
	return ret.Join(linebreak)
}

func (i Imports) RenderHTML(linebreak string) string {
	ret := &util.StringSlice{}
	for _, imp := range i {
		ret.Pushf("@using %s;", imp)
	}
	return ret.Join(linebreak)
}
