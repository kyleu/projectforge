# Templates

Project Forge uses [quicktemplate][1] for HTML and SQL templating, providing a fast, type-safe, and efficient templating system. The template system is designed to generate clean, semantic HTML while maintaining excellent performance and developer experience.

## Overview

The template system provides:
- **Type Safety**: Compile-time checking of template variables and functions
- **High Performance**: Templates are compiled to Go code for maximum speed
- **Layout System**: Hierarchical layouts with component composition
- **Component Integration**: Seamless integration with Project Forge components
- **SQL Templates**: Type-safe SQL query generation

## Basic Template Structure

### Page Template Anatomy

Every HTML page template follows a consistent structure with embedded layouts and type-safe data:

```html
{% import (
  "your-project/app"
  "your-project/app/controller/cutil"
  "your-project/app/util"
  "your-project/views/components"
  "your-project/views/layout"
) %}

{% code type PageMOTD struct {
  layout.Basic
  MOTD string
} %}

{% func (p *PageMOTD) Body(as *app.State, ps *cutil.PageState) %}
  <div class="card">
    <h3>{%= components.SVGIcon(`app`, ps) %} {%s util.AppName %}</h3>
    <div class="mt">
      <p>This is my app!</p>
      <p class="motd">{%s p.MOTD %}</p>
    </div>
  </div>
{% endfunc %}
```

### Template Usage in Controllers

Integrate templates with your controller logic:

```go
func SomeAction(w http.ResponseWriter, r *http.Request) {
	Act("some.action", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := "SomeResult"
		ps.SetTitleAndData("Page Title", ret)
		return Render(r, as, &views.Home{Data: ret}, ps)
	})
}
```

## Template Compilation

Templates are compiled to Go code for optimal performance:

```bash
# Generate template code
./bin/templates.sh

# Templates are compiled to .go files
# Example: views/Home.html -> views/Home.html.go
```


### Security Considerations

1. **HTML Escaping**: Use `{%s %}` for automatic HTML escaping
2. **Raw HTML**: Use `{%s= %}` only for trusted content
3. **URL Escaping**: Use `{%u %}` for URL parameters
4. **JavaScript Escaping**: Use `{%j %}` for JavaScript strings

```html
<!-- Safe HTML escaping -->
<p>User input: {%s userInput %}</p>

<!-- Raw HTML (trusted content only) -->
<div class="content">{%s= trustedHTMLContent %}</div>

<!-- URL escaping -->
<a href="/search?q={%u searchQuery %}">Search</a>

<!-- JavaScript escaping -->
<script>
var message = {%j jsMessage %};
</script>
```

## Testing Templates

### Template Testing

Test your templates with unit tests:

```go
func TestUserCard(t *testing.T) {
    user := &models.User{
        ID:    1,
        Name:  "John Doe",
        Email: "john@example.com",
        IsAdmin: true,
    }

    ps := &cutil.PageState{
        Profile: &models.Profile{ID: 1},
    }

    var buf bytes.Buffer
    WriteUserCard(&buf, user, ps)

    html := buf.String()

    // Test content
    assert.Contains(t, html, "John Doe")
    assert.Contains(t, html, "john@example.com")
    assert.Contains(t, html, "badge-admin")

    // Test structure
    assert.Contains(t, html, `data-user-id="1"`)
    assert.Contains(t, html, `href="/users/1"`)
}
```

This comprehensive template system provides the foundation for building maintainable, performant, and secure web applications with Project Forge.

[1]: https://github.com/valyala/quicktemplate
