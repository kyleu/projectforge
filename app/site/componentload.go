package site

import (
	"strings"

	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/vsite"
)

var componentMetadata = map[string][]string{
	"accordion":    []string{"map", "collapsible UI for multiple sections"},
	"autocomplete": []string{"glasses", "enhances an input to support server-driven search"},
	"code":         []string{"file-text", "script-free code syntax highlighting with theme support"},
	"dom":          []string{"file-contract", "provides TypeScript methods for manipulating the DOM"},
	"editor":       []string{"edit", "form editor with lots of neat features"},
	"flash":        []string{"fire", "temporary notifications to the user"},
	"icons":        []string{"cocktail", "in-document SVG references with theming support"},
	"link":         []string{"link", "enhances links with confirmation prompts and other utilities"},
	"markdown":     []string{"desktop", "renders Markdown as HTML"},
	"menu":         []string{"list", "hierarchical menu with icon support and clean markup"},
	"modal":        []string{"folder-open", "a modal window that appears over the current page"},
	"table":        []string{"table", "utilities for resizable and sortable tables"},
	"tabs":         []string{"handle", "tabbed navigation component for multiple panels"},
	"tags":         []string{"hashtag", "drag/drop tag editor with accessibility support"},
	"templates":    []string{"sitemap", "dynamic HTML pages from a templating engine with full Go support "},
	"theme":        []string{"gift", "light/dark mode support, theme editor and gallery"},
	"websocket":    []string{"hammer", "supports bidirectional communication between a client and server"},
}

func loadComponents() (vsite.Components, error) {
	files, err := doc.FS.ReadDir("components")
	if err != nil {
		return nil, err
	}
	ret := make(vsite.Components, 0, len(files))
	for _, file := range files {
		key := strings.TrimSuffix(file.Name(), ".md")
		md, ok := componentMetadata[key]
		if !ok {
			md = []string{"star", "a web component without documentation"}
		}
		title, html, err := componentTemplate(key, md[0])
		if err != nil {
			return nil, err
		}
		ret = append(ret, &vsite.Component{Key: strings.TrimSuffix(file.Name(), ".md"), Title: title, Description: md[1], Icon: md[0], HTML: html})
	}
	return ret, nil
}
