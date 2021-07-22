package cutil

import (
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/pkg/errors"

	"$PF_PACKAGE$/app/util"
)

var (
	lineNums   *html.Formatter
	noLineNums *html.Formatter
)

func FormatJSON(v interface{}) (string, error) {
	return FormatLang(util.ToJSON(v), "json")
}

func FormatLang(content string, lang string) (string, error) {
	l := lexers.Get(lang)
	return FormatString(content, l)
}

func FormatFilename(content string, filename string) (string, error) {
	l := lexers.Match(filename)
	if l == nil {
		l = lexers.Fallback
	}
	return FormatString(content, l)
}

func FormatString(content string, l chroma.Lexer) (string, error) {
	if l == nil {
		return "", errors.New("no lexer available for this content")
	}
	s := styles.MonokaiLight
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
	ret = strings.ReplaceAll(ret, ">\n", ">")
	return ret, nil
}
