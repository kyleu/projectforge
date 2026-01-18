# Icons

A comprehensive SVG icon system that provides efficient, scalable, and customizable icons throughout your Project Forge application. The icon system uses an optimized pipeline that inlines SVG icons only once while allowing individual styling for each reference. Instead of loading separate icon files or using icon fonts, SVG icons are processed through a pipeline that allows them to be referenced multiple times while only including the actual SVG markup once in the final HTML output.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Efficient Loading**: Each icon SVG is included only once in the markup
- **Individual Styling**: Each icon reference can have unique styling and classes
- **Scalable**: Vector-based icons that scale perfectly at any size
- **Customizable**: Easy to add, modify, and manage custom icons using Project Forge
- **Performance Optimized**: Minimal overhead with smart inlining
- **Template Integration**: Seamless integration with Go templates and TypeScript

## How It Works

The icon system uses a reference-based approach:

1. **Icon Definition**: SVG icons are stored in the system and processed through a pipeline
2. **Reference Creation**: Templates and code create references to icons using helper functions
3. **Automatic Inclusion**: Referenced icons are automatically included at the end of the template
4. **Individual Styling**: Each reference can have unique size, classes, and styling
5. **Deduplication**: The same icon used multiple times only includes the SVG definition once

## Prerequisites

You'll need to import `views/components` at the top of your template to use icon helpers:

```go
{%- import "views/components" -%}
```

## Basic Usage

### Go Template Usage

Use the `SVGRef` helper to reference icons in your Go templates:

```html
<button>{%= components.SVGButton("star", ps) %}</button>
{%= components.SVGRef("star", 20, 20, "some_css_class", ps) %}
{%= components.SVGSimple("star", 20, ps) %}
{%= components.SVGIcon("star", ps) %}
{%= components.SVGBreadcrumb("star", ps) %}
{%= components.SVGInline("star", 18, ps) %}
{%= components.SVGLink("star", ps) %}
```

### TypeScript Usage

In TypeScript/JavaScript code, use the `svgRef` utility function:

```typescript
import {svgRef} from "./util"

// svgRef(key: string, size: number): string
someElement.innerHTML = svgRef("star", 20);
```

**Parameters:**
- `key`: string - The icon identifier/name
- `size`: number - Icon size (applied to both width and height)

## Icon Management

### Adding Custom Icons

To add custom icons to your Project Forge application:

1. **Access the SVG Management UI**: Navigate to the "SVG" page in your Project Forge admin interface
2. **Line Awesome**: Use a keyword to use an icon from [Line Awesome](https://icons8.com/line-awesome)
3. **Simple Icons**: Use a keyword starting with "brand-" to use an icon from [Simple Icons](https://simpleicons.org)
4. **Set Icon Keys**: Assign unique identifiers to your icons
5. **Optimize**: The system automatically optimizes SVGs for web use
6. **Reference**: Use the new icons with the `SVGRef` helpers
