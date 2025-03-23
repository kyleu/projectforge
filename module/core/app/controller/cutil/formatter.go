package cutil

import (
	"fmt"
	h "html"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

var (
	lineNums   *html.Formatter
	noLineNums *html.Formatter
)

func FormatJSON(v any, lineNumLinkAndTitle ...string) (string, error) {
	return FormatLang(util.ToJSON(v), util.KeyJSON, lineNumLinkAndTitle...)
}

func FormatLang(content string, lang string, lineNumLinkAndTitle ...string) (string, error) {
	l := lexers.Get(lang)
	return FormatString(content, l, lineNumLinkAndTitle...)
}

func FormatLangIgnoreErrors(content string, lang string, lineNumLinkAndTitle ...string) string {
	ret, err := FormatLang(content, lang, lineNumLinkAndTitle...)
	if err != nil {
		return fmt.Sprintf("encoding error: %s\n%s", err.Error(), content)
	}
	return ret
}

func FormatFilename(content string, filename string, lineNumLinkAndTitle ...string) (string, error) {
	l := lexers.Match(filename)
	if l == nil {
		l = lexers.Fallback
	}
	return FormatString(content, l, lineNumLinkAndTitle...)
}

func FormatString(content string, l chroma.Lexer, lineNumLinkAndTitle ...string) (string, error) {
	if l == nil {
		return "", errors.New("no lexer available for this content")
	}
	s := styles.MonokaiLight
	var f *html.Formatter
	if strings.Contains(strings.TrimSpace(content), "\n") {
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
	if l.Config().Name == "SQL" {
		ret = strings.ReplaceAll(ret, `<span class="err">$</span>`, `<span class="mi">$</span>`)
	}
	if len(lineNumLinkAndTitle) > 0 {
		var title string
		if len(lineNumLinkAndTitle) > 1 {
			title = lineNumLinkAndTitle[1]
		}
		ret = injectLinks(ret, lineNumLinkAndTitle[0], title)
	}
	ret = strings.Replace(ret, `<td class="lntd"><pre tabindex="0" class="chroma"><span class="lnt">1<br /></span></pre></td>`, "", 1)
	return ret, nil
}

var injectLinksRegex = regexp.MustCompile(`<span class="lnt">(.*?)</span>`)

func injectLinks(ret string, url string, title string) string {
	return injectLinksRegex.ReplaceAllStringFunc(ret, func(match string) string {
		content := injectLinksRegex.FindStringSubmatch(match)[1]
		var num int
		for _, x := range strings.TrimSpace(content) {
			if x >= '0' && x <= '9' {
				num = num*10 + int(x-'0')
			} else {
				break
			}
		}
		u := strings.ReplaceAll(url, "{}", fmt.Sprint(num))
		t := strings.ReplaceAll(title, "{}", fmt.Sprint(num))
		anchor := fmt.Sprintf(`<a target="_blank" rel="noopener noreferrer" title=%q href=%q>`, t, u)
		return strings.Replace(match, content, anchor+content+"</a>", 1)
	})
}

func FormatMarkdown(s string) (string, error) {
	match, end := "<pre><code class=\"language-", "</code></pre>"
	idx := strings.Index(s, match)
	for idx > -1 {
		startQuote := idx + len(match)
		endQuote := strings.Index(s[startQuote:], "\"")
		lang := s[startQuote : startQuote+endQuote]
		if lang == "shell" {
			lang = "bash"
		}
		contentStart := startQuote + endQuote + 2
		contentEnd := strings.Index(s[startQuote:], end) + startQuote
		content := s[contentStart:contentEnd]
		content = h.UnescapeString(content)
		code, err := FormatLang(content, lang)
		if err != nil {
			return "", err
		}
		s = s[:idx] + code + s[contentEnd+len(end):]
		idx = strings.Index(s, match)
	}
	return s, nil
}

func FormatMarkdownClean(s string, icon string) (string, string, error) {
	ret, err := FormatMarkdown(s)
	if err != nil {
		return "", "", err
	}
	var title string
	if h1Idx := strings.Index(ret, "<h1>"); h1Idx > -1 {
		if h1EndIdx := strings.Index(ret, "</h1>"); h1EndIdx > -1 {
			title = s[h1Idx+4 : h1EndIdx]
		}
		ic := fmt.Sprintf(`<svg class="icon" style="width: 20px; height: 20px;"><use xlink:href="#svg-%s" /></svg> `, icon)
		ret = ret[:h1Idx+4] + ic + ret[h1Idx+4:]
		ret = strings.ReplaceAll(ret, "<h3>", "<h4>")
		ret = strings.ReplaceAll(ret, "</h3>", "</h4>")
		ret = strings.ReplaceAll(ret, "<h2>", "<h4>")
		ret = strings.ReplaceAll(ret, "</h2>", "</h4>")
		ret = strings.ReplaceAll(ret, "<h1>", "<h3 style=\"margin-top: 0;\">")
		ret = strings.ReplaceAll(ret, "</h1>", "</h3>")
		ret = strings.ReplaceAll(ret, "<a href=\"http", "<a target=\"_blank\" rel=\"noopener noreferrer\" href=\"http")
	}
	return title, ret, nil
}
