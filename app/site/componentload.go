package site

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/vsite"
)

var componentMetadata = map[string][]string{
	"accordion":    {"map", "collapsible UI for multiple sections"},
	"arguments":    {"mail-open", "collect a set of arguments from a web client"},
	"autocomplete": {"glasses", "enhances an input to support server-driven search"},
	"code":         {"file-text", "script-free code syntax highlighting with theme support"},
	"dom":          {"gamepad", "provides TypeScript methods for manipulating the DOM"},
	"editor":       {"edit", "form editor with lots of neat features"},
	"flash":        {"fire", "temporary notifications to the user"},
	"form":         {"list", "components for building forms backed by Golang data"},
	"icons":        {"icons", "in-document SVG references with theming support"},
	"jsx":          {"dna", "templates for TypeScript objects, using HTML syntax"},
	"link":         {"link", "enhances links with confirmation prompts and other utilities"},
	"loadscreen":   {"refresh", "an interstitial page appearing before a long request"},
	"markdown":     {"desktop", "renders Markdown as HTML"},
	"menu":         {"list", "hierarchical menu with icon support and clean markup"},
	"modal":        {"folder-open", "a modal window that appears over the current page"},
	"search":       {"search", "search framework for custom routines and generated code"},
	"table":        {"table", "utilities for resizable and sortable tables"},
	"tabs":         {"handle", "tabbed navigation component for multiple panels"},
	"tags":         {"hashtag", "drag/drop tag editor with accessibility support"},
	"templates":    {"sitemap", "dynamic HTML pages from a templating engine with full Go support "},
	"theme":        {"gift", "light/dark mode support, theme editor and gallery"},
	"websocket":    {"hammer", "supports bidirectional communication between a client and server"},
}

func loadComponents() (vsite.Components, error) {
	files, err := doc.FS.ReadDir("components")
	if err != nil {
		return nil, err
	}
	ret := make(vsite.Components, 0, len(files))
	for _, file := range files {
		key := strings.TrimSuffix(file.Name(), util.ExtMarkdown)
		md := lo.ValueOr(componentMetadata, key, []string{"star", "a web component without documentation"})
		title, html, err := componentTemplate(key)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &vsite.Component{Key: strings.TrimSuffix(file.Name(), util.ExtMarkdown), Title: title, Description: md[1], Icon: md[0], HTML: html})
	}
	return ret, nil
}
