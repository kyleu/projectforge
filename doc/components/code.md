# Code Highlighting

A powerful syntax highlighting system that transforms code snippets into beautifully formatted, colorized HTML on the server. Built on the robust [Chroma](https://github.com/alecthomas/chroma) library, this component supports hundreds of programming languages and produces clean, accessible markup perfect for documentation, tutorials, and code examples.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Extensive Language Support**: Hundreds of programming languages and file formats
- **Multiple Integration Methods**: Markdown, Go templates, and direct API calls
- **Clean HTML Output**: Semantic markup optimized for accessibility and styling
- **Line Numbers**: Optional line numbering for code references
- **Theme Support**: Integrates with Project Forge's theming system
- **Performance Optimized**: Fast server-side rendering with minimal overhead

## Go Template Integration

Import the components package and use the provided template functions:

```html
{%- import "views/components" -%}

<!-- Display any language with syntax highlighting -->
{%s= cutil.FormatLangIgnoreErrors(codeSnippet, "go") %}

<!-- Display JSON data with syntax highlighting -->
{%= components.JSON(myObject) %}

<!-- Display JSON in a modal dialog -->
{%= components.JSONModal("config", "Configuration", configObject, ps) %}
```

### 3. Go Code Integration

Use the utility functions directly in your Go code:

```go
import "{your project}/app/controller/cutil"

// Format code with explicit language specification
content := `function greet(name) {
    return "Hello, " + name + "!";
}`

resultHTML, err := cutil.FormatLang(content, "javascript")
if err != nil {
    log.Printf("Formatting error: %v", err)
    return
}

// Format code with automatic language detection from filename
cssContent := `body { margin: 0; padding: 20px; }`
resultHTML, err := cutil.FormatFilename(cssContent, "styles.css")

// Format with error handling ignored (for non-critical formatting)
resultHTML := cutil.FormatLangIgnoreErrors(content, "javascript")

// Format any Go object as JSON
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email,omitzero"`
}

user := User{ID: 1, Name: "John Doe", Email: "john@example.com"}
jsonHTML, err := cutil.FormatJSON(user)
```

## API Reference

### Go Functions

#### FormatLang
```go
func FormatLang(content, language string) (string, error)
```
Formats code content with the specified language syntax highlighting.

**Parameters:**
- `content` (string): The source code to format
- `language` (string): Language identifier (e.g., "go", "javascript", "python")

**Returns:**
- `string`: HTML-formatted code with syntax highlighting
- `error`: Error if formatting fails

**Example:**
```go
html, err := cutil.FormatLang(`print("Hello World")`, "python")
```

#### FormatLangIgnoreErrors
```go
func FormatLangIgnoreErrors(content, language string) string
```
Same as FormatLang but returns empty string on errors instead of failing.

**Example:**
```go
html := cutil.FormatLangIgnoreErrors(userCode, "go")
```

#### FormatFilename
```go
func FormatFilename(content, filename string) (string, error)
```
Automatically detects language from file extension and formats accordingly.

**Parameters:**
- `content` (string): The source code to format
- `filename` (string): Filename with extension for language detection

**Example:**
```go
html, err := cutil.FormatFilename(cssCode, "main.css")
html, err := cutil.FormatFilename(goCode, "handler.go")
html, err := cutil.FormatFilename(jsCode, "app.js")
```

#### FormatJSON
```go
func FormatJSON(obj interface{}) (string, error)
```
Formats any Go object as syntax-highlighted JSON.

**Example:**
```go
type Config struct {
    Port     int      `json:"port"`
    Database string   `json:"database"`
    Features []string `json:"features"`
}

config := Config{
    Port:     8080,
    Database: "postgresql://localhost/myapp",
    Features: []string{"auth", "logging", "metrics"},
}

html, err := cutil.FormatJSON(config)
```

### Template Functions

#### JSON Component
```go
{%= components.JSON(object) %}
```
Displays an object as formatted JSON in an HTML table with line numbers.

**Example:**
```html
{%= components.JSON(user) %}
{%= components.JSON(configuration) %}
{%= components.JSON(apiResponse) %}
```

#### JSONModal Component
```go
{%= components.JSONModal(id, title, object, pageState) %}
```
Creates a modal dialog containing formatted JSON data.

**Parameters:**
- `id` (string): Unique identifier for the modal
- `title` (string): Modal window title
- `object` (interface{}): Object to display as JSON
- `pageState` (*cutil.PageState): Current page state

**Example:**
```html
{%= components.JSONModal("user-data", "User Details", user, ps) %}
{%= components.JSONModal("config-view", "Application Configuration", config, ps) %}
```

## Supported Languages

The component supports hundreds of languages through Chroma. Here are some commonly used ones:
