package cutil

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"$PF_PACKAGE$/app/util"
)

var (
	lineNums   *html.Formatter
	noLineNums *html.Formatter
)

func Format(v interface{}, lang string) (string, error) {
	return FormatString(util.ToJSON(v), lang)
}

func FormatString(content string, lang string) (string, error) {
	s := styles.MonokaiLight
	l := lexers.Get(lang)
	var f *html.Formatter
	if strings.Contains(content, "\n") {
		if lineNums == nil {
			lineNums = html.New(html.WithClasses(true), html.WithLineNumbers(true), html.LineNumbersInTable(true))
		}
		f = lineNums
	} else {
		if noLineNums == nil {
			noLineNums = html.New(html.WithClasses(true))
		}
		f = noLineNums
	}
	i, err := l.Tokenise(nil, content)
	if err != nil {
		return "", errors.Wrap(err, "can't tokenize")
	}
	x := &strings.Builder{}
	err = f.Format(x, s, i)
	if err != nil {
		return "", errors.Wrap(err, "can't format")
	}

	ret := x.String()
	ret = strings.ReplaceAll(ret, "\n</span>", "<br /></span>")
	ret = strings.ReplaceAll(ret, "</span>\n", "</span><br />")
	ret = strings.ReplaceAll(ret, "\n<span", "<br /><span")
	ret = strings.ReplaceAll(ret, "\n", "")
	return ret, nil
}
