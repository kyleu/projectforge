# Templates

This project uses [quicktemplate][1] for HTML and SQL templating. 
For HTML, a top-level page is expected to define a struct that embeds `layout.Basic` (or some other layout) and implements `layout.Page`.

A simple page would look like this:

```html
{% import (
  "myproject/app"
  "myproject/app/controller/cutil"
  "myproject/app/util"
  "myproject/views/components"
  "myproject/views/layout"
) %}

{% code type Basic struct {
  layout.Basic
  MOTD string
} %}

{% func (p *Basic) Body(as *app.State, ps *cutil.PageState) %}
  <div class="card">
    <h3>{%= components.SVGIcon(`app`, ps) %} {%s util.AppName %}</h3>
    <div class="mt">
      <p>This is my app!</p>
    </div>
  </div>
{% endfunc %}
```

[1]: https://github.com/valyala/quicktemplate
