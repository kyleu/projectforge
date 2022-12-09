# Code Highlighting

Using a custom pipeline powered by [chroma](https://github.com/alecthomas/chroma), your HTML templates can include colorized formatted code, with or without line numbers.
Hundreds of languages are supported, and the resulting HTML fragment is formatted for clean markup.

For Markdown files, it's as easy as including a triple-backtick code block with a language specified. It's used for this document's code examples. 

For `go` and HTML Templates, you'll need to either call one of the utility methods and display the unescaped markup, or use the predefined JSON components.

### Golang
```go
import "{your project}/app/controller/cutil"

content := `<div>
  Hello!
</div>`

// format a string with a specific language
resultHTML, err := cutil.FormatLang(content, "html")

// or, ignore errors
resultHTML := cutil.FormatLangIgnoreErrors(content, "html")

// or, detect the language from a filename
resultHTML, err := cutil.FormatFilename(content, "index.html")

// or, format the JSON representation of any object
resultHTML, err := cutil.FormatJSON(myObject)
```

### HTML Templates

You'll need to import `views/components` at the top of your template.

```html
<!-- produces an HTML table containing formatted code with line numbers -->
{%= components.JSON(myObject) %}
<!-- produces a modal window containing the code -->
{%= components.JSONModal("unique-key", "Title", myObject, ps) %}
```
