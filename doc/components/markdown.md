# Markdown

Utilities are provided to render Markdown files as HTML. 
This is used to great effect in the `docbrowse` module, which exposes all Markdown in `/doc`.

For rendering documentation, a helper is provided. It accepts a transformation function:

```go
title, html, err := doc.HTML("key", "file.md", func(s string) (string, string, error) {
	// FormatCleanMarkup formats the HTML for displaying in the app, and inserts an optional icon 
	return cutil.FormatCleanMarkup(s, "my-icon")
})
```

To render your own files, simply call the following utility method, then clean it for your template:

```go
import (
	"github.com/gomarkdown/markdown"

	"{your_app}/app/controller/cutil"
)

html := string(markdown.ToHTML(data, nil, nil))

title, html, err := cutil.FormatCleanMarkup(s, "my-icon")
```
